package cpu

import (
	"bufio"
	"helloServer/metrics"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

func init() {
	if _, err := os.Stat("/proc/stat"); os.IsNotExist(err) {
		log.Println("'/proc/stat' is not exist")
		os.Exit(1)
	}
}

type metric struct {
	previousCPU int64
}

func New() *metric {
	return &metric{}
}

func (mt *metric) Process(measure *metrics.Measure) error {
	f, err := os.Open("/proc/stat")
	if err != nil {
		return errors.Wrap(err, "'/proc/stat' open error")
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Scan()

	if scanner.Err() != nil {
		return errors.Wrap(scanner.Err(), "scanner scan error")
	}

	split := strings.Split(scanner.Text(), " ")
	text := split[2]

	currentCPU, err := strconv.ParseInt(text, 10, 64)
	if err != nil {
		return errors.Wrap(err, "strconv parseInt error")
	}

	if mt.previousCPU == 0 {
		mt.previousCPU = currentCPU
	}

	cpuUsage := currentCPU - mt.previousCPU
	measure.Cpu.Usage = cpuUsage

	mt.previousCPU = currentCPU

	return nil
}

func (mt *metric) Once(measure *metrics.Measure) error {
	return nil
}
