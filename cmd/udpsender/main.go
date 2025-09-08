package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	udpAddr, err := net.ResolveUDPAddr("udp", ":42069")
	if err != nil {
		log.Fatalf("Error resolving UDP addr: %v", err)
	}

	udpConn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		log.Fatalf("Error creating UDP connection: %v", err)
	}
	defer udpConn.Close()

	rd := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">")
		line, err := rd.ReadString('\n')
		if err == io.EOF {
			continue
		}
		if err != nil {
			log.Fatalf("Error reading string: %v", err)
		}
		udpConn.Write([]byte(line))
	}

}
