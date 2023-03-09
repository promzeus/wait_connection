package main

import (
	"flag"
	"fmt"
	"github.com/shirou/gopsutil/net"
	"os"
	"strconv"
	"time"
)

func main() {

	portPtr := flag.Int("port", 8080, "port to check connections on")
	flag.Parse()

	//logFile, err := os.OpenFile("/dev/fd/1", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer logFile.Close()
	//log.SetOutput(logFile)

	connCount := 0
	for {
		conns, err := net.Connections("tcp")
		if err != nil {
			fmt.Println("Error getting network connections:", err)
			os.Exit(1)
		}
		for _, conn := range conns {
			if conn.Laddr.Port == uint32(*portPtr) {
				connCount++
			}
		}
		//fmt.Fprintf(logFile, "Connections on port %d: %d\n", *portPtr, connCount)
		fmt.Printf("Connections on port %d: %d\n", *portPtr, connCount)
		if connCount <= 1 {
			os.Exit(0)
		}
		connCount = 0
		time.Sleep(time.Second)
	}
}

// atoi преобразует строку в целое число.
func atoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}
