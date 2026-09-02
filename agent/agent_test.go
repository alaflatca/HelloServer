package agent

import (
	"context"
	"errors"
	"helloServer/measure"
	"sync"
	"testing"
	"time"
)

type testProcessor struct {
	processCh chan struct{}
	onceErr   error
}

func (p *testProcessor) Once(_ *measure.Measure) error {
	return p.onceErr
}

func (p *testProcessor) Process(_ *measure.Measure) error {
	select {
	case p.processCh <- struct{}{}:
	default:
	}
	return nil
}

func TestRunStopsOnContextCancellation(t *testing.T) {
	processor := &testProcessor{processCh: make(chan struct{}, 2)}
	agent := &Agent{
		processor: []Processor{processor},
		period:    time.Hour,
		periodCh:  make(chan time.Duration, 1),
		debugFlag: false,
	}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)

	go func() {
		done <- agent.Run(ctx)
	}()

	waitForProcess(t, processor.processCh, "initial collection did not run")
	cancel()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("Run returned error: %v", err)
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("Run did not stop after context cancellation")
	}
}

func TestRunAppliesUpdatedPeriod(t *testing.T) {
	processor := &testProcessor{processCh: make(chan struct{}, 5)}
	agent := &Agent{
		processor: []Processor{processor},
		period:    time.Hour,
		periodCh:  make(chan time.Duration, 1),
		debugFlag: false,
	}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)

	go func() {
		done <- agent.Run(ctx)
	}()

	waitForProcess(t, processor.processCh, "initial collection did not run")
	agent.UpdatePeriod(10 * time.Millisecond)
	waitForProcess(t, processor.processCh, "period update did not wake collection loop")
	waitForProcess(t, processor.processCh, "updated ticker period did not trigger collection")

	cancel()
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("Run returned error: %v", err)
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("Run did not stop after period update test")
	}

	if agent.period != 10*time.Millisecond {
		t.Fatalf("period = %s, want 10ms", agent.period)
	}
}

func TestUpdatePeriodConcurrentDoesNotBlock(t *testing.T) {
	agent := &Agent{
		period:   time.Hour,
		periodCh: make(chan time.Duration, 1),
	}

	done := make(chan struct{})
	go func() {
		var wg sync.WaitGroup
		for i := 0; i < 100; i++ {
			i := i
			wg.Add(1)
			go func() {
				defer wg.Done()
				agent.UpdatePeriod(time.Duration(i+1) * time.Millisecond)
			}()
		}
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(200 * time.Millisecond):
		t.Fatal("concurrent UpdatePeriod calls blocked")
	}
}

func TestRunReturnsOnceProcessError(t *testing.T) {
	expected := errors.New("once failed")
	agent := &Agent{
		processor: []Processor{&testProcessor{processCh: make(chan struct{}, 1), onceErr: expected}},
		period:    time.Hour,
		periodCh:  make(chan time.Duration, 1),
		debugFlag: false,
	}

	err := agent.Run(context.Background())
	if !errors.Is(err, expected) {
		t.Fatalf("err = %v, want %v", err, expected)
	}
}

func waitForProcess(t *testing.T, ch <-chan struct{}, message string) {
	t.Helper()
	select {
	case <-ch:
	case <-time.After(300 * time.Millisecond):
		t.Fatal(message)
	}
}
