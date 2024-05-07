package system

import (
	"bufio"
	"fmt"
	"helloServer/measure"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

func init() {
	if _, err := os.Stat("/etc/os-release"); os.IsNotExist(err) {
		log.Println("'/etc/os-release' is not exist")
		os.Exit(1)
	}
	if _, err := os.Stat("/proc/uptime"); os.IsNotExist(err) {
		log.Println("'/proc/uptime' is not exist")
		os.Exit(1)
	}
}

type metric struct {
	OS     string
	Uptime string
}

func New() *metric {
	return &metric{}
}

func (mt *metric) Process(measure *measure.Measure) error {
	uptime, err := GetUptime()
	if err != nil {
		return errors.Wrap(err, "Failed to Uptime")
	}
	measure.System.Uptime = uptime

	return nil
}

func GetUptime() (string, error) {
	uptimeBytes, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return "", errors.Wrap(err, "Failed to read file '/proc/uptime'")
	}
	fields := strings.Fields(string(uptimeBytes))
	if len(fields) < 2 {
		return "", errors.New("/proc/uptime read value abnormal")
	}

	field := fields[0]

	for i := 0; i < len(field); i++ {
		if field[i] == '.' {
			field = field[:i]
		}
	}
	uptime, err := strconv.ParseInt(field, 10, 64)
	if err != nil {
		return "", errors.Wrap(err, "uptime parse ( string -> int )")
	}

	var days, hour int64

	days = uptime / 86400
	uptime = uptime % 86400
	hour = uptime / 3600

	return fmt.Sprintf("%d days %d hour", days, hour), nil
}

func (mt *metric) GetOsRelease() (string, error) {
	f, err := os.Open("/etc/os-release")
	if err != nil {
		return "", errors.Wrap(err, "Failed to open '/etc/os-release'")
	}
	defer f.Close()

	var osRelease string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		split := strings.Split(line, "=")
		if len(split) < 2 {
			continue
		}

		if split[0] != "PRETTY_NAME" {
			continue
		}

		osRelease = strings.Trim(split[1], "\"")
	}

	if scanner.Err() != nil {
		return "", errors.Wrap(scanner.Err(), "Failed to os-release scanner")
	}

	return osRelease, nil
}

func (mt *metric) Once(measure *measure.Measure) error {
	osRelease, err := mt.GetOsRelease()
	if err != nil {
		return err
	}

	measure.System.OsRelease = osRelease
	return nil
}
