package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"Heimdall/io"
)

const (
	Port    = 9999
	Batch   = 128
	MsgSize = 64
	Total   = 2000000
)

func main() {
	if len(os.Args) < 2 {
		return
	}
	if os.Args[1] == "server" {
		runServer()
	} else {
		runClient()
	}
}

func runServer() {
	// 1. Initialize Engine with New API
	cfg := io.Config{BatchSize: Batch, MsgSize: MsgSize}
	engine, err := io.NewEngine(Port, cfg)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Heimdall Raw Network Engine [v1.1] Listening on UDP :%d\n", Port)

	count := 0
	start := time.Now()
	for count < Total {
		n := engine.Poll()
		if n > 0 {
			engine.Send(n)
			count += n
		}
	}

	elapsed := time.Since(start)
	stats := engine.GetStats()

	fmt.Printf("\n--- Heimdall Professional Benchmark Result ---\n")
	fmt.Printf("Recv Packets: %d\n", stats["recv_packets"])
	fmt.Printf("Syscalls:     %d\n", stats["syscalls"])
	fmt.Printf("Errors:       %d\n", stats["errors"])
	fmt.Printf("Throughput:   %.2f pps (1 Core)\n", float64(count)/elapsed.Seconds())
}

func runClient() {
	addr, _ := net.ResolveUDPAddr("udp", "127.0.0.1:9999")
	conn, _ := net.DialUDP("udp", nil, addr)
	payload := make([]byte, MsgSize)
	copy(payload, "heimdall-key-01")
	fmt.Println("Client sending...")
	for i := 0; i < Total; i++ {
		conn.Write(payload)
		if i%50000 == 0 {
			time.Sleep(1 * time.Microsecond)
		}
	}
}
