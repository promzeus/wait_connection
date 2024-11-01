package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/shirou/gopsutil/net"
)

func main() {
	portPtr := flag.Int("port", 8080, "port to check connections on")
	expectedNumberConn := flag.Int("conn", 1, "expected number of connections")
	workingTimeLimit := flag.Duration("t", 0, "timeout duration for the program (e.g., 30s or 1m)")
	statesPtr := flag.String("states", "ESTABLISHED,FIN_WAIT_1,FIN_WAIT_2,CLOSE_WAIT,TIME_WAIT", "TCP states to monitor, separated by commas")

	flag.Parse()

	monitoredStates := strings.Split(*statesPtr, ",")
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
			if conn.Laddr.Port == uint32(*portPtr) && contains(monitoredStates, conn.Status) {
				connCount++
			}
		}

		if connCount != prevConnCount {
			fmt.Printf("Active connections on port %d with states %v: %d\n", *portPtr, monitoredStates, connCount)
			prevConnCount = connCount
		}

		if connCount <= int(*expectedNumberConn) {
			duration := time.Since(startTime)
			fmt.Printf("Expected number of active connections reached. Exiting. It took %.2f seconds to close.\n", duration.Seconds())
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

func contains(states []string, state string) bool {
	for _, s := range states {
		if s == state {
			return true
		}
	}
	return false
}
