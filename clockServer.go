// Clock Server is a concurrent TCP server that periodically writes the time.
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
)

func TimeIn(t time.Time, name string) (time.Time, error) {
	loc, err := time.LoadLocation(name)
	if err == nil {
		t = t.In(loc)
	}
	return t, err
}

func handleConn(c net.Conn) {
	defer c.Close()

	for {
		var timeObj, timeErr = TimeIn(time.Now(), os.Getenv("TZ"))
		if timeErr == nil {
			_, err := io.WriteString(c, os.Getenv("TZ")+"\t= "+timeObj.Local().Format("15:04:05\n"))
			if err != nil {
				return // e.g., client disconnected
			}
			time.Sleep(1 * time.Second)
		} else {
			fmt.Println("Invalid timezone")
			fmt.Println()
			os.Exit(-1)
		}
	}

}

func main() {

	if len(os.Args) < 3 {
		fmt.Println("Usage: go run clockServer.go -port <portNum>")
		fmt.Println()
		os.Exit(-1)
	}

	var portNum = flag.String("port", "9090", "Port specified for clock connection")
	flag.Parse()

	listener, err := net.Listen("tcp", "localhost:"+*portNum)
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err) // e.g., connection aborted
			continue
		}
		go handleConn(conn) // handle connections concurrently
	}
}
