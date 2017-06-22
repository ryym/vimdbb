package vimdbb

type QueryPayload struct {
	Query string
}

type Action struct {
	Id   float64
	Type string
}
type Result struct {
	Rows string
}
