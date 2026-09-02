package main

import (
	"context"
	"fmt"
	"helloServer/agent"
	"helloServer/server"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/mbndr/figlet4go"
)

func main() {
	if err := banner(); err != nil {
		log.Println("banner:", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	defer stop()

	agent := agent.New()
	server := server.New()
	errCh := make(chan error, 2)

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		if err := agent.Start(ctx); err != nil && ctx.Err() == nil {
			errCh <- fmt.Errorf("agent start: %w", err)
		}
	}()
	go func() {
		defer wg.Done()
		if err := server.Start(); err != nil && ctx.Err() == nil {
			errCh <- fmt.Errorf("server start: %w", err)
		}
	}()

	select {
	case <-ctx.Done():
	case err := <-errCh:
		log.Println(err)
		stop()
	}

	agent.Close()
	if err := server.Close(); err != nil {
		log.Println("server close:", err)
	}
	wg.Wait()
}

func banner() error {
	ascii := figlet4go.NewAsciiRender()
	renderStr, err := ascii.Render("Hello, Server!")
	if err != nil {
		return err
	}
	fmt.Println(renderStr)
	return nil
}
