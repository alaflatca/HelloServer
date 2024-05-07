package main

import (
	"bytes"
	"encoding/json"
	"helloServer/measure"
	"net"
	"testing"
)

const max int = 4

func TestClient(t *testing.T) {
	conn, err := net.Dial("tcp", "218.239.147.15:9227")
	if err != nil {
		t.Fatal(err)
	}

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
			break
		}

		t.Logf("%+v\n", measure)
	}

}
