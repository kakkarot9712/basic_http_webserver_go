package response

func (r *Response) Message() string {
	messages := map[int]string{
		200: "OK",
		404: "Not Found",
	}
	return messages[r.Status]
}
