package cpu

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseCPUStatLine(t *testing.T) {
	total, idle, err := parseCPUStatLine("cpu  100 20 30 800 50 5 6 7 8 9")
	if err != nil {
		t.Fatal(err)
	}

	if total != 1035 {
		t.Fatalf("total = %d, want 1035", total)
	}
	if idle != 850 {
		t.Fatalf("idle = %d, want 850", idle)
	}
}

func TestParseCPUStatLineInvalidNumber(t *testing.T) {
	if _, _, err := parseCPUStatLine("cpu  100 bad 30 800 50"); err == nil {
		t.Fatal("expected parse error")
	}
}

func TestParseCPUStatLineMalformed(t *testing.T) {
	tests := []string{
		"",
		"cpu0 100 20 30 800 50",
		"cpu 100 20 30",
	}

	for _, tt := range tests {
		t.Run(tt, func(t *testing.T) {
			if _, _, err := parseCPUStatLine(tt); err == nil {
				t.Fatal("expected malformed line error")
			}
		})
	}
}

func TestCalculateUsage(t *testing.T) {
	usage, ok := calculateUsage(1000, 800, 1200, 900)
	if !ok {
		t.Fatal("expected usage calculation")
	}
	if usage != 50 {
		t.Fatalf("usage = %d, want 50", usage)
	}
}

func TestCalculateUsageUsesTotalIdleDeltas(t *testing.T) {
	usage, ok := calculateUsage(100, 60, 160, 90)
	if !ok {
		t.Fatal("expected usage calculation")
	}

	if usage != 50 {
		t.Fatalf("usage = %d, want 50", usage)
	}
}

func TestCalculateUsageRejectsCounterReset(t *testing.T) {
	if _, ok := calculateUsage(1200, 900, 1000, 800); ok {
		t.Fatal("expected counter reset to be rejected")
	}
}

func TestCalculateUsageRejectsInvalidIdleDelta(t *testing.T) {
	if _, ok := calculateUsage(100, 10, 110, 30); ok {
		t.Fatal("expected idle delta greater than total delta to be rejected")
	}
}

func TestReadCPUTimes(t *testing.T) {
	path := filepath.Join(t.TempDir(), "stat")
	if err := os.WriteFile(path, []byte("cpu  100 20 30 800 50 5 6 7 8 9\n"), 0644); err != nil {
		t.Fatal(err)
	}

	total, idle, err := readCPUTimes(path)
	if err != nil {
		t.Fatal(err)
	}
	if total != 1035 {
		t.Fatalf("total = %d, want 1035", total)
	}
	if idle != 850 {
		t.Fatalf("idle = %d, want 850", idle)
	}
}

func TestReadCPUTimesMissingFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing-stat")
	if _, _, err := readCPUTimes(path); err == nil {
		t.Fatal("expected missing file error")
	}
}
