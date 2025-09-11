package main

import (
	"fmt"
	"log"
	"net"

	"github.com/GLobyNew/httpfromtcp/internal/request"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatalf("Error creating listener: %v", err)
	}
	defer listener.Close()

	fmt.Println("Server has started!")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Error accepting connection: %v", err)
		}
		fmt.Println("Connection has been established!")

		r, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatalf("error getting request from reader: %v", err)
		}
		fmt.Println("Request line:")
		fmt.Printf("- Method: %v\n", r.RequestLine.Method)
		fmt.Printf("- Target: %v\n", r.RequestLine.RequestTarget)
		fmt.Printf("- Version: %v\n", r.RequestLine.HttpVersion)

	}
}
