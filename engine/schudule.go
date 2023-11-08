package engine

import (
	"fmt"
	"sync"

	"go.uber.org/zap"

	"github.com/zoulux/crawler-go/collect"
)

type Crawler struct {
	options
	out         chan collect.ParseResult
	Visited     map[string]bool
	VisitedLock sync.Mutex
}

type Scheduler interface {
	Schedule()
	Push(...*collect.Request)
	Pull() *collect.Request
}

type Schedule struct {
	requestCh chan *collect.Request
	workCh    chan *collect.Request
	reqQueue  []*collect.Request
	Logger    *zap.Logger
}

func (s *Schedule) Schedule() {
	for {
		var req *collect.Request
		var ch chan *collect.Request
		if len(s.reqQueue) > 0 {
			req = s.reqQueue[0]
			s.reqQueue = s.reqQueue[1:]
			ch = s.workCh
		}
		select {
		case req := <-s.requestCh:
			s.reqQueue = append(s.reqQueue, req)
		case ch <- req:
			fmt.Println(123)
		}
	}
}

func (s *Schedule) Push(requests ...*collect.Request) {
	for _, req := range requests {
		s.requestCh <- req
	}
}

func (s *Schedule) Pull() *collect.Request {
	r := <-s.workCh
	return r
}

func (e *Crawler) Run() {
	go e.Schedule()

	for i := 0; i < e.WorkCount; i++ {
		go e.CreateWork()
	}
	e.HandleResult()
}

func (e *Crawler) Schedule() {
	var reqQueue []*collect.Request

	for _, seed := range e.Seeds {

		seed.RootReq.Task = seed
		seed.RootReq.Url = seed.Url
		reqQueue = append(reqQueue, seed.RootReq)
	}
	go e.scheduler.Schedule()
	go e.scheduler.Push(reqQueue...)
}

func (e *Crawler) CreateWork() {

	for {
		r := e.scheduler.Pull()
		if err := r.Check(); err != nil {
			e.Logger.Error("crawler check error", zap.Error(err))
			continue
		}

		if e.HasVisited(r) {
			e.Logger.Debug(" request has visited", zap.String("url", r.Url))
			continue
		}

		e.StoreVisited(r)

		body, err := r.Task.Fetcher.Get(r)
		if err != nil {
			e.Logger.Error("can't fetch", zap.Error(err))
			continue
		}
		result := r.ParseFunc(body, r)
		if len(result.Requests) > 0 {
			go e.scheduler.Push(result.Requests...)
		}
		e.out <- result
	}
}

func (e *Crawler) HandleResult() {
	for {
		select {
		case result := <-e.out:
			for _, item := range result.Items {
				// TODO
				e.Logger.Sugar().Info("get result ", item)
			}
		}
	}
}

func (e *Crawler) HasVisited(r *collect.Request) bool {
	e.VisitedLock.Lock()
	defer e.VisitedLock.Unlock()

	unique := r.Unique()
	return e.Visited[unique]
}

func (e *Crawler) StoreVisited(reqs ...*collect.Request) {
	e.VisitedLock.Lock()
	defer e.VisitedLock.Unlock()
	for _, req := range reqs {
		unique := req.Unique()
		e.Visited[unique] = true
	}
}

func NewEngine(opts ...Option) *Crawler {
	options := defaultOptions

	for _, opt := range opts {
		opt(&options)
	}
	e := &Crawler{}
	e.Visited = make(map[string]bool, 100)
	out := make(chan collect.ParseResult)
	e.out = out
	e.options = options
	return e
}

func NewSimpleScheduler() *Schedule {
	s := &Schedule{}
	requestCh := make(chan *collect.Request)
	workCh := make(chan *collect.Request)
	s.requestCh = requestCh
	s.workCh = workCh
	return s
}
