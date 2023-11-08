package collect

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
)

type Request struct {
	unique    string
	Task      *Task
	Url       string
	Method    string
	Depth     int
	ParseFunc func([]byte, *Request) ParseResult
}

func (r *Request) Check() error {
	if r.Depth > r.Task.MaxDepth {
		return errors.New("max depth limit reached")
	}
	return nil
}

func (r *Request) Unique() string {
	block := md5.Sum([]byte(r.Url + r.Method))
	return hex.EncodeToString(block[:])
}

type ParseResult struct {
	Requests []*Request
	Items    []interface{}
}
