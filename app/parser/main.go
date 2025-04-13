package parser

type RequestStage uint8

const (
	Location RequestStage = iota
	Headers
	Body
)

func RequestStructure() *[]RequestStage {
	return &[]RequestStage{
		Location,
		Headers,
		Body,
	}
}
