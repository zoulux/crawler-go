package main

import (
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/zoulux/crawler-go/collect"
	"github.com/zoulux/crawler-go/engine"
	"github.com/zoulux/crawler-go/log"
	"github.com/zoulux/crawler-go/parse/doubangroup"
)

var headerRe = regexp.MustCompile(`<div class="news_li"[\s\S]*?<h2>[\s\S]*?<a.*?target="_blank">([\s\S]*?)</a>`)

func main() {
	plugin, c := log.NewFilePlugin("/dev/stdout", log.InfoLevel)
	defer c.Close()
	logger := log.NewLogger(plugin)
	logger.Info("log init end")

	var f collect.Fetcher = collect.BrowserFetch{
		Timeout: time.Second * 3000,
		Logger:  logger,
	}

	var seeds []*collect.Task
	for i := 25; i <= 100; i += 25 {
		str := fmt.Sprintf("https://www.douban.com/group/szsh/discussion?start=%d", i)
		seeds = append(seeds, &collect.Task{
			Url:      str,
			WaitTime: 1 * time.Second,
			MaxDepth: 5,
			Fetcher:  f,
			Cookie:   "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3282.140 Safari/537.36",
			RootReq: &collect.Request{
				Method:    http.MethodGet,
				ParseFunc: doubangroup.ParseURL,
			},
		})
	}

	s := engine.NewEngine(
		engine.WithFetcher(f),
		engine.WithLogger(logger),
		engine.WithWorkCount(5),
		engine.WithSeeds(seeds),
		engine.WithScheduler(engine.NewSimpleScheduler()),
	)

	s.Run()
}
