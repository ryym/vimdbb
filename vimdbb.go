package vimdbb

type QueryPayload struct {
	Query string
}

type Message struct {
	Id      float64
	Command string
	Payload string
}
type Result struct {
	Rows string
}
