package headers

type Headers map[string]string

func New() Headers {
	return make(Headers)
}

func (h Headers) Get(key string) string {
	return h[key]
}

func (h Headers) Set(key string, value string) {
	h[key] = value
}
