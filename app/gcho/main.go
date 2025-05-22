package gcho

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"net"

	"github.com/codecrafters-io/http-server-starter-go/app/gcho/headers"
	"github.com/codecrafters-io/http-server-starter-go/app/gcho/request"
	"github.com/codecrafters-io/http-server-starter-go/app/gcho/response"
)

type Context struct {
	context.Context
	conn    net.Conn
	version string
	Request request.Request
	res     response.Response
}

func NewContext(conn net.Conn) (*Context, error) {
	buf := bufio.NewReader(conn)
	req, err := request.Parse(buf)
	if err != nil {
		return nil, err
	}
	ctx := Context{
		Request: req,
		version: req.Version,
		conn:    conn,
		res:     response.New(),
	}
	return &ctx, nil
}

func (ctx *Context) Write(status int, p []byte) error {
	ctx.res.Status = status
	if p != nil {
		_, err := ctx.res.Write(p)
		if err != nil {
			return err
		}
	}
	resInfo := fmt.Sprintf("%v %v %v\r\n", ctx.Request.Version, ctx.res.Status, ctx.res.Message())
	headers := ctx.res.Headers
	buf := bytes.NewBuffer([]byte(resInfo))
	for k, v := range headers {
		buf.WriteString(fmt.Sprintf("%v: %v\r\n", k, v))
	}
	_, err := buf.WriteString("\r\n")
	if err != nil {
		return err
	}
	_, err = buf.Write(ctx.res.Bytes())
	if err != nil {
		return err
	}
	ctx.conn.Write(buf.Bytes())
	return nil
}

func (ctx *Context) Headers() headers.Headers {
	return ctx.res.Headers
}
