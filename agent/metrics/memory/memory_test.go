package memory

import (
	"strings"
	"testing"
)

func TestParseMemInfo(t *testing.T) {
	input := `MemTotal:        8089672 kB
MemFree:         4679572 kB
MemAvailable:    6354024 kB
Buffers:          586580 kB
Cached:          1174896 kB
SwapCached:            0 kB
`

	total, available, cached, err := parseMemInfo(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}

	if total != 8089672 {
		t.Fatalf("total = %f, want 8089672", total)
	}
	if available != 6354024 {
		t.Fatalf("available = %f, want 6354024", available)
	}
	if cached != 1174896 {
		t.Fatalf("cached = %f, want 1174896", cached)
	}
}

func TestParseMemInfoMissingField(t *testing.T) {
	input := `MemTotal:        8089672 kB
Cached:          1174896 kB
`

	if _, _, _, err := parseMemInfo(strings.NewReader(input)); err == nil {
		t.Fatal("expected missing field error")
	}
}
