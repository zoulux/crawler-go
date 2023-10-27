package collect

import (
	"errors"
)

type Request struct {
	Task      *Task
	Url       string
	Depth     int
	ParseFunc func([]byte, *Request) ParseResult
}

func (r Request) Check() error {
	if r.Depth > r.Task.MaxDepth {
		return errors.New("max depth limit reached")
	}
	return nil
}

type ParseResult struct {
	Requests []*Request
	Items    []interface{}
}
