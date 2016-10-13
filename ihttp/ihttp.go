package ihttp

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
type JSONBody map[string]interface{}
type HTTPHandler func(*Request)

type Request struct {
	Headers Headers
	Body    Body
	Url     string
	Method  string
}

const lineEnd = "\r\n"

func (r *Request) JsonBody() JSONBody {
	var dat map[string]interface{}
	json.Unmarshal(r.Body, &dat)
	return dat
}

func (r *Request) addData(headers Headers, body Body) {
	r.Headers = headers
	r.Body = body
	method, url, _ := methodUrlVersion(headers)
	r.Method = method
	r.Url = url
}

func ListenAndServe(host string, handler HTTPHandler) {
	ln, err := net.Listen("tcp", host)
	if err != nil {
		// handle error
	}
	fmt.Printf("Successfully listening on : %s\n", host)
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
		}
		go handleConnection(conn, handler)
	}
}

func methodUrlVersion(headers Headers) (string, string, string) {
	requestDetails := headers["request"]
	splitted := strings.Split(requestDetails, " ")
	return splitted[0], splitted[1], splitted[2]
}

func handleConnection(conn net.Conn, handler HTTPHandler) {
	bufferSize := 2048
	buf := make([]byte, bufferSize)
	headerEnd := []byte(lineEnd + lineEnd)
	headersDone := false
	var (
		message        []byte
		headers        Headers
		body           Body
		contentLength  int
		bodyStartIndex int
	)

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
	handler(request)
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
