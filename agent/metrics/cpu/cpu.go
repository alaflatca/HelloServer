package cpu

import (
	"bufio"
	"fmt"
	"helloServer/measure"
	"helloServer/utils"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type metric struct {
	previousIdle uint64
	previousCPU  uint64
	criteriaTime time.Time
}

func New() *metric {
	return &metric{
		criteriaTime: utils.Now(),
	}
}

func (mt *metric) Process(measure *measure.Measure) error {
	total, idle, err := readCPUTimes("/proc/stat")
	if err != nil {
		return err
	}

	if mt.previousCPU == 0 {
		mt.previousCPU = total
		mt.previousIdle = idle
		return nil
	}

	cpuUsage, ok := calculateUsage(mt.previousCPU, mt.previousIdle, total, idle)
	if !ok {
		mt.previousCPU = total
		mt.previousIdle = idle
		return nil
	}
	measure.Cpu.Usage = cpuUsage
	mt.previousCPU = total
	mt.previousIdle = idle

	// cpu usage csv
	if time.Since(mt.criteriaTime) > (10 * time.Second) {
		mt.criteriaTime = time.Now()
		// record
		record := []string{utils.NowString(), strconv.Itoa(int(cpuUsage))}
		if err := utils.WriteCSV(utils.Metric_CPU_CSV, record); err != nil {
			return err
		}
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

func readCPUTimes(path string) (uint64, uint64, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, 0, errors.Wrap(err, fmt.Sprintf("%q open error", path))
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return 0, 0, errors.Wrap(err, "scan cpu stat")
		}
		return 0, 0, errors.New("cpu stat is empty")
	}
	total, idle, err := parseCPUStatLine(scanner.Text())
	if err != nil {
		return 0, 0, err
	}
	if err := scanner.Err(); err != nil {
		return 0, 0, errors.Wrap(err, "scan cpu stat")
	}
	return total, idle, nil
}

func parseCPUStatLine(line string) (uint64, uint64, error) {
	fields := strings.Fields(line)
	if len(fields) < 5 {
		return 0, 0, fmt.Errorf("invalid cpu stat line: %q", line)
	}
	if fields[0] != "cpu" {
		return 0, 0, fmt.Errorf("invalid cpu stat prefix: %q", fields[0])
	}

	values := make([]uint64, 0, len(fields)-1)
	for _, field := range fields[1:] {
		value, err := strconv.ParseUint(field, 10, 64)
		if err != nil {
			return 0, 0, errors.Wrapf(err, "parse cpu stat field %q", field)
		}
		values = append(values, value)
	}

	var total uint64
	for _, value := range values {
		total += value
	}

	idle := values[3]
	if len(values) > 4 {
		idle += values[4]
	}
	return total, idle, nil
}

func calculateUsage(previousTotal, previousIdle, currentTotal, currentIdle uint64) (int64, bool) {
	if currentTotal <= previousTotal || currentIdle < previousIdle {
		return 0, false
	}

	totalDelta := currentTotal - previousTotal
	idleDelta := currentIdle - previousIdle
	if totalDelta == 0 || idleDelta > totalDelta {
		return 0, false
	}

	usage := float64(totalDelta-idleDelta) * 100 / float64(totalDelta)
	return int64(usage + 0.5), true
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

		split := strings.SplitN(line, ":", 2)
		if len(split) < 2 {
			continue
		}
		if strings.TrimSpace(split[0]) != "model name" {
			continue
		}
		modelName = strings.TrimSpace(split[1])
		break
	}

	if scanner.Err() != nil {
		return "", errors.Wrap(scanner.Err(), "scan /proc/cpuinfo")
	}
	if modelName == "" {
		return "", errors.New("cpu model name not found")
	}

	return modelName, nil
}
