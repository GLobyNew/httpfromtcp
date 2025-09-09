package main

import (
	"fmt"
	"log"
	"net"
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

		ch := getLinesChannel(conn)
		for s := range ch {
			fmt.Println(s)
		}

	}
}
