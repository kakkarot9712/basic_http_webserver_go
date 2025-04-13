package request

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/codecrafters-io/http-server-starter-go/app/http"
)

type Request struct {
	Path    string
	Method  http.Method
	Version string
}

func Parse(r *bufio.Reader) (req Request, err error) {
	line, err := r.ReadString(byte('\n'))
	if err != nil {
		return Request{}, err
	}
	line = strings.Trim(line, "\r\n")
	reqInfo := strings.SplitN(line, " ", 3)
	if len(reqInfo) != 3 {
		return Request{}, fmt.Errorf("inavalid request structure found")
	}
	if !http.CheckMethod(reqInfo[0]) {
		return Request{}, fmt.Errorf("invalid method detected")
	}
	if !http.CheckVersion(reqInfo[2]) {
		return Request{}, fmt.Errorf("unsupported http version found")
	}
	req.Method = http.Method(reqInfo[0])
	req.Path = reqInfo[1]
	req.Version = reqInfo[2]
	// for {
	// 	header, err := r.ReadString(byte('\n'))
	// 	if err != nil {
	// 		return Request{}, err
	// 	}
	// 	if header == "\r\n" {
	// 		break
	// 	}
	// }
	return req, nil
}
