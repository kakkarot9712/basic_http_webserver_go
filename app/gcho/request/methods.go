package request

import (
	"slices"
)

const (
	GET     = "GET"
	POST    = "POST"
	PATCH   = "PATCH"
	DELETE  = "DELETE"
	OPTIONS = "OPTIONS"
	PUT     = "PUT"
)

func supportedMethods() []string {
	return []string{
		GET,
		POST,
		PATCH,
		DELETE,
		OPTIONS,
		PUT,
	}
}

func checkMethod(m string) (ok bool) {
	return slices.Contains(supportedMethods(), m)
}
