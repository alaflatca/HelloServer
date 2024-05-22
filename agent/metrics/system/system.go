package system

import (
	"bufio"
	"fmt"
	"helloServer/measure"
	"log"
	"os"
	"os/user"
	"path/filepath"
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

type ProcessInfo struct {
	Uid  string `json:"uid"`
	Pid  string `json:"pid"`
	PPid string `json:"ppid"`
	Time string `json:"time"`
	Cmd  string `json:"cmd"`
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

func LookupUID(uid string) string {
	f, err := os.Open("/etc/passwd")
	if err != nil {
		return ""
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		split := strings.Split(line, ":")
		userName := split[0]
		uid := split[2]
	}
}

func (mt *metric) Once(measure *measure.Measure) error {
	osRelease, err := mt.GetOsRelease()
	if err != nil {
		return err
	}

	measure.System.OsRelease = osRelease
	return nil
}

// openat(AT_FDCWD, "/proc/659/stat", O_RDONLY) = 6
// openat(AT_FDCWD, "/proc/659/status", O_RDONLY) = 6
// openat(AT_FDCWD, "/proc/659/cmdline", O_RDONLY) = 6
func GetProcessList() ([]*ProcessInfo, error) {
	dirEntry, err := os.ReadDir("/proc")
	if err != nil {
		return nil, err
	}

	processInfoList := make([]*ProcessInfo, 0)
	for _, dir := range dirEntry {
		if dir.IsDir() {
			name := dir.Name()
			if '0' < name[0] && name[0] <= '9' {
				info := &ProcessInfo{}
				procStatus(info, filepath.Join("/proc", name, "status"))
				procCmdline(info, filepath.Join("/proc", name, "cmdline"))
				processInfoList = append(processInfoList, info)
			}
		}
	}

	return processInfoList, nil
}

func procStatus(info *ProcessInfo, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(f)

	for i := 0; i < 9; i++ {
		if !scanner.Scan() {
			return scanner.Err()
		}

		line := scanner.Text()
		split := strings.Fields(line)
		switch split[0] {
		case "Pid:":
			info.Pid = strings.TrimSpace(split[1])
		case "PPid:":
			info.PPid = strings.TrimSpace(split[1])
		case "Uid:":
			uid := strings.TrimSpace(split[1])
			user, err := user.LookupId(uid)
			if err != nil {
				return err
			}
			fmt.Printf("uid: %s, name: %s\n", uid, user.Name)
			info.Uid = user.Name
		default:
			continue
		}
	}

	return nil
}

func procCmdline(info *ProcessInfo, path string) error {
	cmd, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	info.Cmd = string(cmd)
	return nil
}
