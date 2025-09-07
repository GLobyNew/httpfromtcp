package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	f, err := os.Open("messages.txt")
	if err != nil {
		log.Fatalf("Error openning file: %v", err)
	}
	defer f.Close()

	ch := getLinesChannel(f)

	fmt.Printf("read: %s\n", <-ch)

}
