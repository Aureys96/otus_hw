package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	timeout := flag.Duration("timeout", 10*time.Second, "connection timeout")
	flag.Parse()

	if len(flag.Args()) != 2 {
		log.Fatalln("host and port must be provided")
	}
	host := flag.Arg(0)
	port := flag.Arg(1)

	client := NewTelnetClient(net.JoinHostPort(host, port), *timeout, os.Stdin, os.Stdout)

	if err := client.Connect(); err != nil {
		log.Fatalln("failed to connect:", err)
	}
	defer client.Close()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT)
	defer cancel()

	go func() {
		client.Send()
		cancel()
	}()
	go func() {
		client.Receive()
		cancel()
	}()

	<-ctx.Done()
}
