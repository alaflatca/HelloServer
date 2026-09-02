package system

import (
	"strings"
	"testing"
)

func TestParseUptime(t *testing.T) {
	uptime, err := parseUptime("90061.47 123456.00")
	if err != nil {
		t.Fatal(err)
	}

	if uptime != "1 days 1 hour" {
		t.Fatalf("uptime = %q, want %q", uptime, "1 days 1 hour")
	}
}

func TestParseUptimeInvalid(t *testing.T) {
	if _, err := parseUptime("not-a-number 123456.00"); err == nil {
		t.Fatal("expected parse error")
	}
}

func TestParseOSRelease(t *testing.T) {
	release, err := parseOSRelease(strings.NewReader(`NAME="Ubuntu"
VERSION="20.04.6 LTS (Focal Fossa)"
PRETTY_NAME="Ubuntu 20.04.6 LTS"
`))
	if err != nil {
		t.Fatal(err)
	}

	if release != "Ubuntu 20.04.6 LTS" {
		t.Fatalf("release = %q, want %q", release, "Ubuntu 20.04.6 LTS")
	}
}

func TestParseOSReleaseMissingPrettyName(t *testing.T) {
	if _, err := parseOSRelease(strings.NewReader(`NAME="Ubuntu"`)); err == nil {
		t.Fatal("expected missing PRETTY_NAME error")
	}
}
