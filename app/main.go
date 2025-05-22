package main

import (
	"fmt"
	"net"
	"os"

	"github.com/codecrafters-io/http-server-starter-go/app/hub"
)

func main() {
	// Uncomment this block to pass the first stage
	//
	h := hub.NewHub(10)
	h.Start()
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		h.WorkReceiver <- &hub.Task{
			C: conn,
		}
	}
}
