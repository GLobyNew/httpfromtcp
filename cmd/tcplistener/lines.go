package main

import (
	"io"
	"log"
	"net"
	"strings"
)

func getLinesChannel(f net.Conn) <-chan string {

	ch := make(chan string)

	readMsg := make([]byte, 8)
	curLine := ""
	go func() {
		for {
			n, err := f.Read(readMsg)
			if err == io.EOF {
				f.Close()
				close(ch)
				break
			}
			if err != nil {
				log.Printf("Error reading 8 bytes: %v", err)
			}

			splitLine := strings.Split(string(readMsg[:n]), "\n")
			curLine += splitLine[0]
			if len(splitLine) > 1 {
				ch <- curLine
				curLine = strings.Join(splitLine[1:], "\n")
			}
		}
	}()
	return ch

}
