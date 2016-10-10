package main

import (
	"bytes"
	"fmt"
	"net"
	"strconv"
	"strings"
)

type Headers map[string]string

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
	bufferSize := 2048
	buf := make([]byte, bufferSize)
	lineEnd := "\r\n"
	headerEnd := []byte(lineEnd + lineEnd)
	headersDone := false
	var message []byte
	var headers Headers
	var body []byte
	var contentLength int
	var bodyStartIndex int
	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading:", err.Error())
			break
		}
		message = append(message, buf[:n]...)

		// handling headers
		if headersDone == false {
			headerEndIndex := bytes.Index(message, headerEnd)
			if headerEndIndex > 0 {
				headers = parseHeaders(string(message[:headerEndIndex]), lineEnd)
				contentLength, _ = getContentLength(headers)
				bodyStartIndex = headerEndIndex + len(headerEnd)
				headersDone = true
			}
		}

		// handling body
		if headersDone == true {
			bodyLength := len(message[bodyStartIndex:])
			if bodyLength >= contentLength {
				body = message[bodyStartIndex : bodyStartIndex+contentLength]
				break
			}
		}
	}
	handleMessage(headers, body)
}

func parseHeaders(headers string, lineEnd string) Headers {
	splitHeaders := strings.Split(headers, lineEnd)
	finalHeaders := map[string]string{"request": splitHeaders[0]}
	for _, value := range splitHeaders[1:] {
		split := strings.SplitN(value, ":", 2)
		finalHeaders[strings.Trim(split[0], " ")] = strings.Trim(split[1], " ")
	}
	return finalHeaders
}

func getContentLength(headers Headers) (int, error) {
	return strconv.Atoi(headers["Content-Length"])
}

func handleMessage(headers Headers, body []byte) {
	fmt.Println("handling message")
	fmt.Println(headers)
	fmt.Println(string(body))
}
