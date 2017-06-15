package vimdbb

type Payload struct {
	Query string
}

type Action struct {
	Id      float64
	Payload Payload
}
type Result struct {
	Rows string
}
