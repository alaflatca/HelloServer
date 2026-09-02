package cpu

import "testing"

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

func TestCalculateUsage(t *testing.T) {
	usage, ok := calculateUsage(1000, 800, 1200, 900)
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
