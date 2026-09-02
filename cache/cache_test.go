package cache

import (
	"helloServer/measure"
	"testing"
)

func TestSetGetReturnsSnapshot(t *testing.T) {
	key := "snapshot-test"
	Set(key, &measure.Measure{})

	first := Get(key)
	if first == nil {
		t.Fatal("expected cached measure")
	}
	first.Cpu.Usage = 99

	second := Get(key)
	if second == nil {
		t.Fatal("expected cached measure")
	}
	if second.Cpu.Usage == 99 {
		t.Fatal("Get returned shared measure pointer")
	}
}
