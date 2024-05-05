package main

import (
	"fmt"
	"helloServer/agent"
	"helloServer/server"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mbndr/figlet4go"
)

// qr코드를 이용한 앱등록 서버(qr코드) <-> 핸드폰앱에서 qr코드로 인식 기능 사용
func main() {
	banner()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	agent := agent.New()
	server := server.New()

	go agent.Start()
	go server.Start()

	<-sigs

	agent.Close()
	server.Close()
}

func banner() {
	ascii := figlet4go.NewAsciiRender()
	renderStr, err := ascii.Render("Hello, Server!")
	if err != nil {
		log.Fatal("banner: ", err.Error())
	}
	fmt.Println(renderStr)
}
