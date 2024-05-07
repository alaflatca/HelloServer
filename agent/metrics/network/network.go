package network

import (
	"bufio"
	"helloServer/measure"
	"helloServer/utils"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

func init() {
	if _, err := os.Stat("/proc/net/dev"); os.IsNotExist(err) {
		log.Println("'/proc/net/dev' is not exist")
		os.Exit(1)
	}

	if _, err := os.Stat("daily_network.csv"); os.IsNotExist(err) {
		f, err := os.Create("daily_network.csv")
		if err != nil {
			os.Exit(1)
		}
		f.Close()
	}
}

const (
	IFACE    = 0
	RX_INDEX = 1
	TX_INDEX = 9
)

type metric struct {
	previousRxBytes float64
	previousTxBytes float64
}

func New() *metric {
	return &metric{}
}

func (mt *metric) Process(ms *measure.Measure) error {
	f, err := os.Open("/proc/net/dev")
	if err != nil {
		return errors.Wrap(err, "Failed to open '/proc/net/dev'")
	}
	defer f.Close()

	// Inter-|   Receive                                                    |  Transmit
	//  face |bytes    packets errs drop fifo frame compressed multicast|   bytes    packets errs drop fifo colls carrier compressed
	//  lo:    441220207  267382    0    0    0     0          0         0  441220207  267382    0    0    0     0       0          0
	//  eth0: 2652673230 1791793    0    0    0     0          0      1543  30048187  411110    0    0    0     0       0          0
	scanner := bufio.NewScanner(f)
	scanner.Scan()
	scanner.Scan()
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 17 {
			continue
		}

		if fields[0] == "lo:" {
			continue
		}

		ms.Network.Iface = fields[IFACE]
		rxBytes, _ := strconv.ParseFloat(fields[RX_INDEX], 64)
		txBytes, _ := strconv.ParseFloat(fields[TX_INDEX], 64)

		if mt.previousRxBytes == 0 && mt.previousTxBytes == 0 {
			mt.previousRxBytes = rxBytes
			mt.previousTxBytes = txBytes
			return nil
		}

		ms.Network.RxUsage = float64(rxBytes-mt.previousRxBytes) / measure.KB
		ms.Network.TxUsage = float64(txBytes-mt.previousTxBytes) / measure.KB

		mt.previousRxBytes = rxBytes
		mt.previousTxBytes = txBytes

		// daily record
		daily := utils.NowDaily()
		if daily.After(ms.Network.DailyTime) { // 2024-04-12 > 2024-04-11
			ms.Network.DailyTime = daily
			// utils.dailyWriteCSV
			// csv format
			// 2024-04-11,{first_byte},{last_byte - first_byte}
		}

	}
	return nil
}

func GetLocalIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

func GetLocalIPs() ([]net.IP, error) {
	var ips []net.IP
	addresses, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	for _, addr := range addresses {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ips = append(ips, ipnet.IP)
			}
		}
	}
	return ips, nil
}

func (mt *metric) Once(ms *measure.Measure) error {
	ms.Network.IPaddress = GetLocalIP()
	ms.Network.DailyTime = utils.NowDaily()
	// 최초 실행 시간 저장

	if err := utils.InitializeCSV(utils.Metric_NETWORK_CSV); err != nil {
		return err
	}

	return nil
}
