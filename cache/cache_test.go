package cache

import (
	"helloServer/measure"
	"sync"
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

func TestConcurrentSetGetReturnsIndependentSnapshots(t *testing.T) {
	key := "snapshot-concurrent-test"
	Set(key, &measure.Measure{})

	const iterations = 1000
	var wg sync.WaitGroup

	for writer := 0; writer < 4; writer++ {
		writer := writer
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				ms := &measure.Measure{}
				ms.Cpu.Usage = int64(writer*iterations + i)
				ms.Memory.Total = float64(i + 1)
				ms.Network.Iface = "writer"
				Set(key, ms)
			}
		}()
	}

	for reader := 0; reader < 8; reader++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				ms := Get(key)
				if ms == nil {
					continue
				}
				ms.Cpu.Usage = -1
				ms.Network.Iface = "mutated-by-reader"
			}
		}()
	}

	wg.Wait()

	final := Get(key)
	if final == nil {
		t.Fatal("expected cached measure")
	}
	if final.Cpu.Usage == -1 || final.Network.Iface == "mutated-by-reader" {
		t.Fatal("reader mutation leaked into cached snapshot")
	}
}
