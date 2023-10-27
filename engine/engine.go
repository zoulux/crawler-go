package engine

import (
	"go.uber.org/zap"

	"github.com/zoulux/crawler-go/collect"
)

type ScheduleEngine struct {
	requestCh chan *collect.Request
	workCh    chan *collect.Request
	out       chan collect.ParseResult
	WorkCount int
	Logger    *zap.Logger
	Fetcher   collect.Fetcher
	Seeds     []*collect.Request
}

func (s *ScheduleEngine) Run() {
	requestCh := make(chan *collect.Request)
	workCh := make(chan *collect.Request)
	out := make(chan collect.ParseResult)
	s.requestCh = requestCh
	s.workCh = workCh
	s.out = out
	go s.Schedule()
	for i := 0; i < s.WorkCount; i++ {
		go s.CreateWork()
	}
	s.HandleResult()
}

func (s *ScheduleEngine) Schedule() {
	var reqQueue = s.Seeds
	go func() {
		for {
			var req *collect.Request
			var ch chan *collect.Request
			if len(reqQueue) > 0 {
				req = reqQueue[0]
				reqQueue = reqQueue[1:]
				ch = s.workCh
			}

			select {
			case r := <-s.requestCh:
				reqQueue = append(reqQueue, r)
			case ch <- req:

			}
		}
	}()
}

func (s *ScheduleEngine) CreateWork() {

	for {
		r := <-s.workCh
		body, err := s.Fetcher.Get(r)
		if err != nil {
			s.Logger.Error("can't fetch", zap.Error(err))
			continue
		}
		result := r.ParseFunc(body, r)
		s.out <- result
	}
}

func (s *ScheduleEngine) HandleResult() {
	for {
		select {
		case result := <-s.out:
			for _, req := range result.Requests {
				s.requestCh <- req
			}

			for _, item := range result.Items {
				// TODO
				s.Logger.Sugar().Info("get result ", item)
			}
		}
	}
}
