package network

import (
	"bufio"
	"fmt"
	"helloServer/measure"
	"helloServer/utils"
	"io"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

const (
	RX_INDEX = 1
	TX_INDEX = 9
)

type metric struct {
	iface           string
	previousRxBytes uint64
	previousTxBytes uint64
}

type deviceStat struct {
	iface   string
	rxBytes uint64
	txBytes uint64
}

func New() *metric {
	return &metric{}
}

func (mt *metric) Process(ms *measure.Measure) error {
	stats, err := readNetDev("/proc/net/dev")
	if err != nil {
		return err
	}

	stat, ok := mt.selectDevice(stats)
	if !ok {
		return errors.New("network interface not found")
	}

	ms.Network.Iface = stat.iface
	ms.Network.IPaddress = getIPv4ForInterface(stat.iface)

	if mt.iface != stat.iface {
		mt.iface = stat.iface
		mt.previousRxBytes = stat.rxBytes
		mt.previousTxBytes = stat.txBytes
		ms.Network.RxUsage = 0
		ms.Network.TxUsage = 0
		return nil
	}

	if stat.rxBytes < mt.previousRxBytes || stat.txBytes < mt.previousTxBytes {
		mt.previousRxBytes = stat.rxBytes
		mt.previousTxBytes = stat.txBytes
		ms.Network.RxUsage = 0
		ms.Network.TxUsage = 0
		return nil
	}

	ms.Network.RxUsage = float64(stat.rxBytes-mt.previousRxBytes) / measure.KB
	ms.Network.TxUsage = float64(stat.txBytes-mt.previousTxBytes) / measure.KB

	mt.previousRxBytes = stat.rxBytes
	mt.previousTxBytes = stat.txBytes

	// daily record
	daily := utils.NowDaily()
	if daily.After(ms.Network.DailyTime) { // 2024-04-12 > 2024-04-11
		ms.Network.DailyTime = daily
		// utils.dailyWriteCSV
		// csv format
		// 2024-04-11,{first_byte},{last_byte - first_byte}
	}

	return nil
}

func (mt *metric) selectDevice(stats []deviceStat) (deviceStat, bool) {
	return mt.selectDeviceWithActive(stats, activeInterfaceNames())
}

func (mt *metric) selectDeviceWithActive(stats []deviceStat, active map[string]struct{}) (deviceStat, bool) {
	if mt.iface != "" {
		for _, stat := range stats {
			if stat.iface == mt.iface {
				return stat, true
			}
		}
	}

	for _, stat := range stats {
		if stat.iface == "lo" {
			continue
		}
		if _, ok := active[stat.iface]; ok {
			return stat, true
		}
	}

	for _, stat := range stats {
		if stat.iface != "lo" {
			return stat, true
		}
	}
	return deviceStat{}, false
}

func readNetDev(path string) ([]deviceStat, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("open %q", path))
	}
	defer f.Close()

	stats, err := parseNetDev(f)
	if err != nil {
		return nil, err
	}
	return stats, nil
}

func parseNetDev(r io.Reader) ([]deviceStat, error) {
	scanner := bufio.NewScanner(r)
	var stats []deviceStat
	for scanner.Scan() {
		stat, ok, err := parseNetDevLine(scanner.Text())
		if err != nil {
			return nil, err
		}
		if ok {
			stats = append(stats, stat)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, errors.Wrap(err, "scan /proc/net/dev")
	}
	if len(stats) == 0 {
		return nil, errors.New("network statistics not found")
	}
	return stats, nil
}

func parseNetDevLine(line string) (deviceStat, bool, error) {
	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 {
		return deviceStat{}, false, nil
	}

	iface := strings.TrimSpace(parts[0])
	fields := strings.Fields(parts[1])
	if len(fields) < 16 {
		return deviceStat{}, false, fmt.Errorf("invalid network stat line: %q", line)
	}

	rxBytes, err := strconv.ParseUint(fields[RX_INDEX-1], 10, 64)
	if err != nil {
		return deviceStat{}, false, errors.Wrapf(err, "parse rx bytes for %s", iface)
	}
	txBytes, err := strconv.ParseUint(fields[TX_INDEX-1], 10, 64)
	if err != nil {
		return deviceStat{}, false, errors.Wrapf(err, "parse tx bytes for %s", iface)
	}

	return deviceStat{iface: iface, rxBytes: rxBytes, txBytes: txBytes}, true, nil
}

func activeInterfaceNames() map[string]struct{} {
	active := make(map[string]struct{})
	interfaces, err := net.Interfaces()
	if err != nil {
		return active
	}

	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		active[iface.Name] = struct{}{}
	}
	return active
}

func getIPv4ForInterface(name string) string {
	iface, err := net.InterfaceByName(name)
	if err != nil {
		return ""
	}

	addresses, err := iface.Addrs()
	if err != nil {
		return ""
	}
	for _, addr := range addresses {
		ipnet, ok := addr.(*net.IPNet)
		if !ok || ipnet.IP.IsLoopback() {
			continue
		}
		if ip := ipnet.IP.To4(); ip != nil {
			return ip.String()
		}
	}
	return ""
}

func GetLocalIP() (string, error) {
	ips, err := GetLocalIPs()
	if err != nil {
		return "", err
	}
	if len(ips) == 0 {
		return "", nil
	}
	return ips[0].String(), nil
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
	ip, err := GetLocalIP()
	if err != nil {
		return err
	}
	ms.Network.IPaddress = ip
	ms.Network.DailyTime = utils.NowDaily()
	// 최초 실행 시간 저장

	if err := utils.InitializeCSV(utils.Metric_NETWORK_CSV); err != nil {
		return err
	}

	return nil
}
