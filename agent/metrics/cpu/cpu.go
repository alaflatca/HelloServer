package cpu

import (
	"bufio"
	"fmt"
	"helloServer/measure"
	"helloServer/utils"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

func init() {
	if _, err := os.Stat("/proc/stat"); os.IsNotExist(err) {
		log.Println("'/proc/stat' is not exist")
		os.Exit(1)
	}
}

type metric struct {
	previousCPU  int64
	criteriaTime time.Time
}

func New() *metric {
	return &metric{
		criteriaTime: utils.Now(),
	}
}

func (mt *metric) Process(measure *measure.Measure) error {
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

	// cpu 사용률은 현재 사용률 - 이전 사용률 = 사용률
	// 현재 2초 간격으로 데이터를 가져오고 있음
	// 이전 CPU 값을 저장하고 2초 후에는 현재 cpu 값이 많이 늘어남 그래서 위 공식대로 계산 시 100을 넘게됨

	if mt.previousCPU == 0 {
		mt.previousCPU = currentCPU
		return nil
	}
	cpuUsage := currentCPU - mt.previousCPU

	measure.Cpu.Usage = cpuUsage
	mt.previousCPU = currentCPU

	// cpu usage csv
	if time.Since(mt.criteriaTime) > (10 * time.Second) {
		mt.criteriaTime = time.Now()
		// record
		record := []string{utils.NowString(), strconv.Itoa(int(cpuUsage))}
		utils.WriteCSV(utils.Metric_CPU_CSV, record)
	}

	return nil
}

func (mt *metric) Once(measure *measure.Measure) error {
	if err := utils.InitializeCSV(utils.Metric_CPU_CSV); err != nil {
		return err
	}

	modelName, err := GetCpuModelName()
	if err != nil {
		return err
	}
	measure.Cpu.ModelName = modelName

	return nil

}

func GetCpuModelName() (string, error) {
	f, err := os.Open("/proc/cpuinfo")
	if err != nil {
		return "", err
	}
	defer f.Close()

	modelName := ""
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println("line: ", line)

		split := strings.Split(line, ":")
		if len(split) < 2 {
			continue
		}
		if strings.TrimSpace(split[0]) != "model name" {
			continue
		}
		modelName = strings.TrimSpace(split[1])
		fmt.Println("modelName: ", modelName)
		break
	}

	return modelName, nil
}
