package agent

import (
	"flag"
	"helloServer/agent/httpd"
	"helloServer/agent/metrics"
	"helloServer/agent/metrics/cpu"
	"helloServer/agent/metrics/disk"
	"helloServer/agent/metrics/memory"
	"helloServer/agent/metrics/network"
	"helloServer/agent/metrics/system"
	"log"
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
	measure   *metrics.Measure
	period    time.Duration
	config    Config
	server    Server
	// log 추가
}

type Processor interface {
	Process(*metrics.Measure) error
	Once(*metrics.Measure) error
}

func (a *Agent) addmetric(process ...Processor) {
	if len(process) < 1 {
		log.Println("Requires at least one required metric")
		os.Exit(1)
	}

	a.processor = make([]Processor, 0, len(process))
	a.processor = append(a.processor, process...)
}

type Server interface {
	Serve()
	Close()
}

func (a *Agent) addHttpd(server Server) {
	a.server = server
}

func Start() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	agent := &Agent{measure: &metrics.Measure{}, period: 1}
	agent.argumentParse()
	agent.addHttpd(httpd.New())
	agent.addmetric(system.New(), cpu.New(), memory.New(), disk.New(), network.New())

	go agent.Run()
	go agent.Serve()

	<-sigs
	agent.Close()
}

func (a *Agent) Close() {
	log.Println("Fiber Shutdown")
	a.server.Close()
}

func (a *Agent) Serve() {
	a.server.Serve()
}

func (a *Agent) Run() {
	measure := &metrics.Measure{}

	if err := a.OnceProcess(measure); err != nil {
		panic(err) // OnceProcess 는 반드시 실행, 에러 발생 시 패닉 후 에러 파악
	}

	for {
		now := time.Now()

		for i := 0; i < len(a.processor); i++ {
			if err := a.processor[i].Process(measure); err != nil {
				log.Printf("[%d] processor error: %s", i, err.Error())
				// TODO... error 처리 작성
			}
		}

		measure.Elapse = time.Since(now).String()
		measure.Show()

		time.Sleep(a.period * time.Second)
	}
}

func (a *Agent) OnceProcess(measure *metrics.Measure) error {
	for i := 0; i < len(a.processor); i++ {
		if err := a.processor[i].Once(measure); err != nil {
			log.Printf("[%d] once error: %s\n", i, err.Error())
			return err
		}
	}

	return nil
}

func (a *Agent) argumentParse() {
	port := flag.String("port", "9227", "htpp server port ex) -port=8080")
	flag.Parse()

	a.config.port = *port
}
