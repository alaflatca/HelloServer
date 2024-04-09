package agent

import (
	"encoding/json"
	"log"
	"net"
	"os"
	"time"
)

func (a *Agent) Serve() {
	var err error
	a.listener, err = net.Listen("tcp", "0.0.0.0:"+a.config.port)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	for {
		conn, err := a.listener.Accept()
		if err != nil {
			log.Println(err)
		}

		go a.clientHandler(conn)
	}
}

func (a *Agent) clientHandler(c net.Conn) {
	defer c.Close()

	for {
		bdata, err := json.Marshal(a.measure)
		if err != nil {
			log.Println(err)
			break
		}

		_, err = c.Write(bdata)
		if err != nil {
			log.Println(err)
			break
		}

		time.Sleep(a.period * time.Second)
	}
}
