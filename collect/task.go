package collect

import (
	"time"
)

type Task struct {
	Url      string
	Cookie   string
	WaitTime time.Duration
	MaxDepth int
	RootReq  *Request
	Fetcher  Fetcher
}
