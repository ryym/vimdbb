package vimdbb

type QueryPayload struct {
	ConnectionURL string
	QueryID string
	Query   string
}

type Message struct {
	ID      float64
	Command string
	Payload string
}

type Result struct {
	QueryID string
	Rows    string
}
