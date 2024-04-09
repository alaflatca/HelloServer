package main

import (
	"fmt"
	"helloServer/agent"
	"log"

	"github.com/mbndr/figlet4go"
)

func main() {
	// qr코드를 이용한 앱등록 서버(qr코드) <-> 핸드폰앱에서 qr코드로 인식 기능 사용
	banner()
	agent.Start()
}

func banner() {
	ascii := figlet4go.NewAsciiRender()
	renderStr, err := ascii.Render("Hello, Server")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(renderStr)
}
