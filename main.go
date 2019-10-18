package main

import (
	"fmt"
	"net"
	"tiny_server/tiny_http"
)

func main() {
	port := ":8080"
	ln, err := net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Listen Port %s\n", port)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		go tiny_http.HandleConnection(conn)
	}
}
