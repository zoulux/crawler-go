package collect

type Request struct {
	Url       string
	Cookie    string
	ParseFunc func([]byte) ParseResult
}

type ParseResult struct {
	Requests []*Request
	Items    []interface{}
}
