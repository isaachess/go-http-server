package main

import (
	"./ihttp"
	"fmt"
)

func main() {
	ihttp.ListenAndServe(":3001", handleMessage)
}

func handleMessage(request *ihttp.Request) {
	fmt.Println(request.Headers)
	fmt.Println(request.Body)
	fmt.Println(request.Method)
	fmt.Println(request.Url)
	fmt.Println(request.JsonBody())
}
