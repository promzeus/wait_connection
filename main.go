package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/shirou/gopsutil/net"
)

func main() {
	portPtr := flag.Int("port", 8080, "port to check connections on")
	expectedNumberConn := flag.Int("conn", 1, "expected number of connections")
	workingTimeLimit := flag.Duration("t", 0, "timeout duration for the program (e.g., 30s or 1m)")
	flag.Parse()

	prevConnCount := -1

	startTime := time.Now()

	for {
		conns, err := net.Connections("tcp")
		if err != nil {
			fmt.Println("Error getting network connections:", err)
			os.Exit(1)
		}

		connCount := 0

		for _, conn := range conns {
			if conn.Laddr.Port == uint32(*portPtr) && conn.Status == "ESTABLISHED" {
				connCount++
			}
		}

		if connCount != prevConnCount {
			fmt.Printf("Active connections on port %d: %d\n", *portPtr, connCount)
			prevConnCount = connCount
		}

		if connCount <= int(*expectedNumberConn) {
			fmt.Println("Expected number of active connections reached. Exiting.")
			os.Exit(0)
		}

		if *workingTimeLimit > 0 {
			elapsed := time.Since(startTime)
			if elapsed >= *workingTimeLimit {
				fmt.Printf("Timeout of %v reached. Exiting.\n", *workingTimeLimit)
				os.Exit(0)
			}
		}

		time.Sleep(time.Second)
	}
}
