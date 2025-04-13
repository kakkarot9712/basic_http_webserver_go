package main

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/codecrafters-io/http-server-starter-go/app/gcho"
)

func main() {
	// Uncomment this block to pass the first stage
	//
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	ctx, err := gcho.NewContext(conn)
	if err != nil {
		panic(err)
	}

	req := ctx.Request
	switch {
	case req.Path == "/":
		ctx.Write(200, nil)
	case strings.HasPrefix(req.Path, "/echo"):
		echoStr := strings.TrimPrefix(req.Path, "/echo/")
		ctx.Write(200, []byte(echoStr))
	case req.Path == "/user-agent":
		ctx.Write(200, []byte(ctx.Request.Headers.Get("User-Agent")))
	default:
		ctx.Write(404, nil)
	}
}
