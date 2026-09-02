package system

import (
	"bufio"
	"fmt"
	"helloServer/measure"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

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
	return parseUptime(string(uptimeBytes))
}

func parseUptime(raw string) (string, error) {
	fields := strings.Fields(raw)
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

	return parseOSRelease(f)
}

func parseOSRelease(r io.Reader) (string, error) {
	var osRelease string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		split := strings.SplitN(line, "=", 2)
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
	if osRelease == "" {
		return "", errors.New("PRETTY_NAME not found in /etc/os-release")
	}

	return osRelease, nil
}

func LookupUID(uid string) string {
	usr, err := user.LookupId(uid)
	if err != nil {
		return uid
	}
	return usr.Username
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
		if !dir.IsDir() {
			continue
		}

		name := dir.Name()
		if _, err := strconv.Atoi(name); err != nil {
			continue
		}

		info := &ProcessInfo{}
		if err := procStatus(info, filepath.Join("/proc", name, "status")); err != nil {
			if os.IsNotExist(err) || os.IsPermission(err) {
				continue
			}
			return nil, errors.Wrapf(err, "read process status for pid %s", name)
		}
		if err := procCmdline(info, filepath.Join("/proc", name, "cmdline")); err != nil {
			if os.IsNotExist(err) {
				continue
			}
			if os.IsPermission(err) {
				info.Cmd = ""
			} else {
				return nil, errors.Wrapf(err, "read process cmdline for pid %s", name)
			}
		}
		processInfoList = append(processInfoList, info)
	}

	return processInfoList, nil
}

func procStatus(info *ProcessInfo, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()
		split := strings.Fields(line)
		if len(split) < 2 {
			continue
		}
		switch split[0] {
		case "Pid:":
			info.Pid = strings.TrimSpace(split[1])
		case "PPid:":
			info.PPid = strings.TrimSpace(split[1])
		case "Uid:":
			uid := strings.TrimSpace(split[1])
			user, err := user.LookupId(uid)
			if err != nil {
				info.Uid = uid
				continue
			}
			info.Uid = user.Username
		default:
			continue
		}
	}

	return scanner.Err()
}

func procCmdline(info *ProcessInfo, path string) error {
	cmd, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	info.Cmd = strings.TrimSpace(strings.ReplaceAll(string(cmd), "\x00", " "))
	return nil
}
