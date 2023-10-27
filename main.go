package main

import (
	"fmt"
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

	var seeds []*collect.Task
	for i := 25; i <= 100; i += 25 {
		str := fmt.Sprintf("https://www.douban.com/group/szsh/discussion?start=%d", i)
		seeds = append(seeds, &collect.Task{
			Url:      str,
			WaitTime: 1 * time.Second,
			MaxDepth: 5,
			RootReq: &collect.Request{
				ParseFunc: doubangroup.ParseURL,
			},
		})
	}
	var f collect.Fetcher = collect.BrowserFetch{
		Timeout: time.Second * 3000,
		Logger:  logger,
	}
	s := engine.NewSchedule(
		engine.WithLogger(logger),
		engine.WithWorkCount(5),
		engine.WithFetcher(f),
		engine.WithSeeds(seeds),
	)

	s.Run()
}
