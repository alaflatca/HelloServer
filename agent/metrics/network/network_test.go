package network

import (
	"strings"
	"testing"
)

func TestParseNetDev(t *testing.T) {
	input := `Inter-|   Receive                                                |  Transmit
 face |bytes    packets errs drop fifo frame compressed multicast|bytes    packets errs drop fifo colls carrier compressed
   lo: 441220207 267382 0 0 0 0 0 0 441220207 267382 0 0 0 0 0 0
 ens5: 2652673230 1791793 0 0 0 0 0 1543 30048187 411110 0 0 0 0 0 0
`

	stats, err := parseNetDev(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}

	if len(stats) != 2 {
		t.Fatalf("len(stats) = %d, want 2", len(stats))
	}
	if stats[1].iface != "ens5" {
		t.Fatalf("iface = %q, want ens5", stats[1].iface)
	}
	if stats[1].rxBytes != 2652673230 {
		t.Fatalf("rxBytes = %d, want 2652673230", stats[1].rxBytes)
	}
	if stats[1].txBytes != 30048187 {
		t.Fatalf("txBytes = %d, want 30048187", stats[1].txBytes)
	}
}

func TestParseNetDevLineInvalidNumber(t *testing.T) {
	_, _, err := parseNetDevLine("ens5: bad 1791793 0 0 0 0 0 1543 30048187 411110 0 0 0 0 0 0")
	if err == nil {
		t.Fatal("expected parse error")
	}
}

func TestSelectDeviceKeepsCurrentInterface(t *testing.T) {
	mt := &metric{iface: "docker0"}
	stat, ok := mt.selectDevice([]deviceStat{
		{iface: "ens5", rxBytes: 100, txBytes: 200},
		{iface: "docker0", rxBytes: 300, txBytes: 400},
	})
	if !ok {
		t.Fatal("expected device")
	}
	if stat.iface != "docker0" {
		t.Fatalf("iface = %q, want docker0", stat.iface)
	}
}
