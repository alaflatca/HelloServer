package agent

import (
	"helloServer/agent/measure"
	"helloServer/agent/metrics/cpu"
	"helloServer/agent/metrics/disk"
	"helloServer/agent/metrics/memory"
	"helloServer/agent/metrics/network"
	"helloServer/agent/metrics/system"
	"log"
	"os"
	"time"
)

type Agent struct {
	processor []Processor
	measure   measure.Measure
	period    time.Duration
	breakFlag bool
	// log 추가
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
	return &Agent{measure: measure.Measure{}, period: 1}
}

func (a *Agent) Start() {
	a.addmetric(system.New(), cpu.New(), memory.New(), disk.New(), network.New())
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

	for {
		if a.breakFlag {
			break
		}

		ms = &measure.Measure{}
		for i := 0; i < len(a.processor); i++ {
			if err := a.processor[i].Process(ms); err != nil {
				log.Printf("[%d] processor error: %s", i, err.Error())
				// TODO... error 처리 작성
			}
		}
		ms.Show()
		measure.Set("l", ms)

		time.Sleep(a.period * time.Second)
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
