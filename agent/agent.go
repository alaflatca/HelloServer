package agent

import (
	"helloServer/agent/metrics/cpu"
	"helloServer/agent/metrics/disk"
	"helloServer/agent/metrics/memory"
	"helloServer/agent/metrics/network"
	"helloServer/agent/metrics/system"
	"helloServer/measure"
	"log"
	"os"
	"sync"
	"time"
)

type Agent struct {
	processor []Processor
	measure   measure.Measure
	period    time.Duration
	breakFlag bool
	debugFlag bool
}

type Processor interface {
	Process(*measure.Measure) error
	Once(*measure.Measure) error
}

func (a *Agent) addmetric(process ...Processor) {
	if len(process) < 1 {
		log.Println("Requires at least one required metric")
		os.Exit(1)
	}

	a.processor = make([]Processor, 0, len(process))
	a.processor = append(a.processor, process...)
}

func New() *Agent {
	return &Agent{measure: measure.Measure{}, period: 1 * time.Second, debugFlag: true}
}

func (a *Agent) Start() {
	a.addmetric(system.New(), cpu.New(), memory.New(), disk.New(), network.New())
	measure.Subscribe("period", func(data interface{}) {
		period, ok := data.(time.Duration)
		if !ok {
			return
		}
		mtx := sync.Mutex{}
		mtx.Lock()
		a.period = period
		mtx.Unlock()
	})
	a.Run()
}

func (a *Agent) Close() {
	log.Println("Agent Close")
	a.breakFlag = true
}

func (a *Agent) Run() {
	ms := &measure.Measure{}

	if err := a.OnceProcess(ms); err != nil {
		panic(err) // OnceProcess 는 반드시 실행, 에러 발생 시 패닉 후 에러 파악
	}

	for !a.breakFlag {
		for i := 0; i < len(a.processor); i++ {
			if err := a.processor[i].Process(ms); err != nil {
				log.Printf("[%d] processor error: %s", i, err.Error())
				// TODO... error 처리 작성
			}
		}

		if a.debugFlag {
			ms.Show()
		}

		measure.Set("l", ms)
		time.Sleep(a.period)
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
