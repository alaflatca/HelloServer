package network

import (
	"bufio"
	"helloServer/metrics"
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
}

const (
	IFACE    = 0
	RX_BYTES = 1
	TX_BYTES = 9
)

type metric struct {
	previousRxBytes float64
	previousTxBytes float64
}

func New() *metric {
	return &metric{}
}

func (mt *metric) Process(measure *metrics.Measure) error {
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

		measure.Network.Iface = fields[IFACE]
		rxBytes, _ := strconv.ParseFloat(fields[RX_BYTES], 64)
		txBytes, _ := strconv.ParseFloat(fields[TX_BYTES], 64)

		if mt.previousRxBytes == 0 && mt.previousTxBytes == 0 {
			mt.previousRxBytes = rxBytes
			mt.previousTxBytes = txBytes
			return nil
		}

		measure.Network.RxUsage = (rxBytes - mt.previousRxBytes) / metrics.KB
		measure.Network.TxUsage = (txBytes - mt.previousTxBytes) / metrics.KB

		mt.previousRxBytes = rxBytes
		mt.previousTxBytes = txBytes

		// // // fmt.Println("IP:\t\t", GetLocalIP())
		// // fmt.Printf("Rx:\t\t%0.2fKB\n", RxUsage/1024)
		// fmt.Printf("Tx:\t\t%0.2fKB\n", TxUsage/1024)
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

func (mt *metric) Once(measure *metrics.Measure) error {
	measure.Network.IPaddress = GetLocalIP()
	return nil
}
