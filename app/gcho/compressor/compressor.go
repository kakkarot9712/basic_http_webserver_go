package compressor

import (
	"bytes"
	"compress/gzip"
	"strings"

	"github.com/codecrafters-io/http-server-starter-go/app/gcho/request"
)

func SupportedSchema() map[string]bool {
	return map[string]bool{
		"gzip": true,
	}
}

func ParseAcceptEncoders(req request.Request) []string {
	if req.Headers.Get("Accept-Encoding") == "" {
		return []string{}
	}
	providedEncoders := strings.Split(req.Headers.Get("Accept-Encoding"), ", ")
	encoders := []string{}
	for _, e := range providedEncoders {
		if SupportedSchema()[e] {
			encoders = append(encoders, e)
		}
	}
	return encoders
}

func Compress(encodingShema string, raw []byte) []byte {
	var compressed bytes.Buffer
	switch encodingShema {
	case "gzip":
		gzipEncoder := gzip.NewWriter(&compressed)
		gzipEncoder.Write(raw)
		gzipEncoder.Close()
	}
	return compressed.Bytes()
}
