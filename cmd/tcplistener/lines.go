package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
)

func getLinesChannel(f net.Conn) <-chan string {

	ch := make(chan string)

	go func() {
		defer f.Close()
		defer close(ch)
		curLine := ""
		for {
			readMsg := make([]byte, 8)
			n, err := f.Read(readMsg)
			if err != nil {
				if curLine != "" {
					ch <- curLine
				}
				if errors.Is(err, io.EOF) {
					break
				}
				fmt.Printf("error: %v\n", err.Error())
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
