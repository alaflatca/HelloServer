package main

import (
	"helloServer/agent/metrics/system"
	"testing"
)

func TestXxx(t *testing.T) {
	if _, err := system.GetProcessList(); err != nil {
		t.Fatal(err)
	}

}
