package agent

import (
	"context"
	"errors"
	"helloServer/agent/metrics/cpu"
	"helloServer/agent/metrics/disk"
	"helloServer/agent/metrics/memory"
	"helloServer/agent/metrics/network"
	"helloServer/agent/metrics/system"
	"helloServer/cache"
	"helloServer/event"
	"helloServer/measure"
	"log"
	"sync"
	"time"
)

type Agent struct {
	processor []Processor
	measure   measure.Measure
	period    time.Duration
	debugFlag bool
	periodCh  chan time.Duration
	cancelMu  sync.Mutex
	cancel    context.CancelFunc
}

type Processor interface {
	Process(*measure.Measure) error
	Once(*measure.Measure) error
}

func (a *Agent) addmetric(process ...Processor) error {
	if len(process) < 1 {
		return errors.New("requires at least one metric processor")
	}

	a.processor = make([]Processor, 0, len(process))
	a.processor = append(a.processor, process...)
	return nil
}

func New() *Agent {
	return &Agent{
		measure:   measure.Measure{},
		period:    1 * time.Second,
		debugFlag: true,
		periodCh:  make(chan time.Duration, 1),
	}
}

func (a *Agent) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	a.cancelMu.Lock()
	a.cancel = cancel
	a.cancelMu.Unlock()
	defer func() {
		cancel()
		a.cancelMu.Lock()
		a.cancel = nil
		a.cancelMu.Unlock()
	}()

	if err := a.addmetric(system.New(), cpu.New(), memory.New(), disk.New(), network.New()); err != nil {
		return err
	}

	event.Subscribe("period", func(data interface{}) {
		period, ok := data.(time.Duration)
		if !ok {
			return
		}
		a.UpdatePeriod(period)
	})

	return a.Run(ctx)
}

func (a *Agent) Close() {
	log.Println("Agent Close")
	a.cancelMu.Lock()
	cancel := a.cancel
	a.cancelMu.Unlock()
	if cancel != nil {
		cancel()
	}
}

func (a *Agent) UpdatePeriod(period time.Duration) {
	if period <= 0 {
		return
	}

	select {
	case a.periodCh <- period:
	default:
		select {
		case <-a.periodCh:
		default:
		}
		select {
		case a.periodCh <- period:
		default:
		}
	}
}

func (a *Agent) Run(ctx context.Context) error {
	ms := &measure.Measure{}

	if err := a.OnceProcess(ms); err != nil {
		return err
	}

	ticker := time.NewTicker(a.period)
	defer ticker.Stop()

	for {
		for i := 0; i < len(a.processor); i++ {
			if err := a.processor[i].Process(ms); err != nil {
				log.Printf("[%d] processor error: %s", i, err.Error())
			}
		}

		if a.debugFlag {
			ms.Show()
		}

		cache.Set("l", ms)

		select {
		case <-ctx.Done():
			return nil
		case period := <-a.periodCh:
			a.period = period
			ticker.Reset(period)
			log.Printf("Agent period changed: %s", period)
		case <-ticker.C:
		}
	}
}

func (a *Agent) OnceProcess(measure *measure.Measure) error {
	for i := 0; i < len(a.processor); i++ {
		if err := a.processor[i].Once(measure); err != nil {
			log.Printf("[%d] once error: %s\n", i, err.Error())
			return err
		}
	}
	return nil
}
