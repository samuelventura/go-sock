package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"time"
)

func main() {
	log.SetFlags(0)

	ctrlc := make(chan os.Signal, 1)
	signal.Notify(ctrlc, os.Interrupt)

	if len(os.Args) != 3 {
		name := filepath.Base(os.Args[0])
		log.Printf("Usage: %s proxy_ep target_ep", name)
		log.Printf("Sample: %s localhost:55666 localhost:80", name)
		os.Exit(1)
	}

	proxy_ep := os.Args[1]
	target_ep := os.Args[2]

	conn, err := net.DialTimeout("tcp", proxy_ep, 5*time.Second)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	line := fmt.Sprintln(target_ep)
	n, err := conn.Write([]byte(line))
	if err != nil {
		log.Println(err)
		return
	}
	if n != len(line) {
		log.Println(fmt.Errorf("write mismatch %d %d", len(line), n))
		return
	}
	done := make(chan interface{})
	go func() {
		io.Copy(os.Stdout, conn)
		done <- true
	}()
	go func() {
		io.Copy(conn, os.Stdin)
		done <- true
	}()
	select {
	case <-ctrlc:
	case <-done:
	}
}
