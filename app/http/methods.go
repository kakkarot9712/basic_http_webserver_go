package http

import (
	"slices"
)

type Method string

const (
	GET     Method = "GET"
	POST    Method = "POST"
	PATCH   Method = "PATCH"
	DELETE  Method = "DELETE"
	OPTIONS Method = "OPTIONS"
	PUT     Method = "PUT"
)

func SupportedMethods() []Method {
	return []Method{
		GET,
		POST,
		PATCH,
		DELETE,
		OPTIONS,
		PUT,
	}
}

func CheckMethod(m string) (ok bool) {
	return slices.Contains(SupportedMethods(), Method(m))
}
