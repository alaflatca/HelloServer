package main

import (
	"bytes"
	"encoding/json"
	"helloServer/measure"
	"net"
	"os"
	"testing"
)

const max int = 4

func TestClient(t *testing.T) {
	addr := os.Getenv("HELLOSERVER_TEST_ADDR")
	if addr == "" {
		t.Skip("set HELLOSERVER_TEST_ADDR to run the legacy tcp client test")
	}

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	for i := 0; i < max; i++ {
		data := make([]byte, 1024)
		n, err := conn.Read(data)
		if err != nil {
			t.Fatal(err)
		}
		if bytes.Equal(data[:n], []byte("ping")) {
			t.Log("ping")
			continue
		}
		measure := measure.Measure{}
		if err = json.Unmarshal(data[:n], &measure); err != nil {
			t.Fatal(err)
		}

		t.Logf("%+v\n", measure)
	}

}
