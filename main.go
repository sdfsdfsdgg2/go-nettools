package main

import (
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

func scanPort(host string, port int, wg *sync.WaitGroup, results chan<- int) {
	defer wg.Done()
	addr := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", addr, 2*time.Second)
	if err != nil {
		return
	}
	conn.Close()
	results <- port
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: nettools <host>")
		os.Exit(1)
	}
	host := os.Args[1]
	fmt.Printf("Scanning %s...\n", host)

	results := make(chan int, 1024)
	var wg sync.WaitGroup

	for port := 1; port <= 1024; port++ {
		wg.Add(1)
		go scanPort(host, port, &wg, results)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for port := range results {
		fmt.Printf("  %d/tcp open\n", port)
	}
}