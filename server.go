package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"strings"
)

type Headers map[string]string
type Body []byte

type Request struct {
	headers Headers
	body    Body
	url     string
	method  string
}

func (r *Request) jsonBody() (string, error) {
	body, err := json.Marshal(string(r.body))
	return string(body), err
}

func (r *Request) addData(headers Headers, body Body) {
	r.headers = headers
	r.body = body
	method, url, _ := methodUrlVersion(headers)
	r.method = method
	r.url = url
}

func methodUrlVersion(headers Headers) (string, string, string) {
	requestDetails := headers["request"]
	splitted := strings.Split(requestDetails, " ")
	return splitted[0], splitted[1], splitted[2]
}

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
	var body Body
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
			bodyLength := len(message) - bodyStartIndex
			if bodyLength >= contentLength {
				body = message[bodyStartIndex : bodyStartIndex+contentLength]
				break
			}
		}
	}
	request := new(Request)
	request.addData(headers, body)
	handleMessage(request)
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

func handleMessage(request *Request) {
	fmt.Println(request.headers)
	fmt.Println(request.body)
	fmt.Println(request.method)
	fmt.Println(request.url)
	fmt.Println(request.jsonBody())
}
