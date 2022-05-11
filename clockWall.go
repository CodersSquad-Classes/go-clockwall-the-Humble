package main

import (
	"io"
	"log"
	"net"
	"os"
	"strings"
)

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}

func displayTime(port string, done chan string) {

	conn, err := net.Dial("tcp", "localhost:"+port)
	if err != nil {
		log.Fatal(err)
	}
	mustCopy(os.Stdout, conn)
	done <- port
}

func main() {
	// Let's start the fun
	done := make(chan string, 3)
	for _, element := range os.Args[1:] {
		var port = strings.Split(element, ":")[1]
		go displayTime(port, done)
	}
	_ = <-done
}
