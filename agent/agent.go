package agent

import (
	"flag"
	"helloServer/metrics"
	"helloServer/metrics/cpu"
	"helloServer/metrics/disk"
	"helloServer/metrics/memory"
	"helloServer/metrics/network"
	"helloServer/metrics/system"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Config struct {
	port string
}

type Agent struct {
	processor []Processor
	listener  net.Listener
	measure   *metrics.Measure
	period    time.Duration
	config    Config
	// log 추가
}

type Processor interface {
	Process(*metrics.Measure) error
	Once(*metrics.Measure) error
}

func (a *Agent) addmetric(process ...Processor) {
	if len(process) < 4 {
		log.Println("Failed to Default metrics length ~~~~~")
		os.Exit(1)
	}

	a.processor = make([]Processor, 0, len(process))
	a.processor = append(a.processor, process...)
}

// tcp server add
func Start() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	agent := &Agent{measure: &metrics.Measure{}, period: 2}
	agent.argumentParse()
	agent.addmetric(system.New(), cpu.New(), memory.New(), disk.New(), network.New())

	go agent.Run()
	go agent.Serve()

	<-sigs
	agent.Close()
}

func (a *Agent) Close() {
	if err := a.listener.Close(); err != nil {
		log.Println("listener close: ", err.Error())
		return
	}
}

func (a *Agent) Run() {
	measure := &metrics.Measure{}
	a.OnceProcess(measure)

	for {
		now := time.Now()

		for i := 0; i < len(a.processor); i++ {
			if err := a.processor[i].Process(measure); err != nil {
				log.Printf("[%d] processor error: %s", i, err.Error())
				// TODO... error 처리 작성
			}
		}
		measure.Elapse = time.Since(now).String()
		// mutex
		*a.measure = *measure
		measure.Show()

		time.Sleep(a.period * time.Second)
		// 전체 for문에 wg?
		// tcp Listen 하는 변수와 동기화
	}
}

func (a *Agent) OnceProcess(measure *metrics.Measure) {
	for i := 0; i < len(a.processor); i++ {
		if err := a.processor[i].Once(measure); err != nil {
			log.Printf("[%d] once error: %s", i, err.Error())
			// TODO... error 처리 작성
		}
	}
}

func (a *Agent) argumentParse() {
	port := flag.String("port", "9227", "tcp server port ex) -port=8080")
	flag.Parse()

	a.config.port = *port
}
