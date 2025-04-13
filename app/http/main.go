package http

import "slices"

const (
	HTTPV1 = "HTTP/1.1"
)

func versions() []string {
	return []string{HTTPV1}
}

func CheckVersion(v string) (ok bool) {
	return slices.Contains(versions(), v)
}
