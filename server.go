package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	port := 3001
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		// handle error
	}
	fmt.Printf("Successfully listening on port: %d\n", port)
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		status, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Println(status)
	}
}
