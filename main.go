package main

import (
	"fmt"
	"regexp"
	"time"

	"go.uber.org/zap"

	"github.com/zoulux/crawler-go/collect"
	"github.com/zoulux/crawler-go/log"
	"github.com/zoulux/crawler-go/parse"
)

var headerRe = regexp.MustCompile(`<div class="news_li"[\s\S]*?<h2>[\s\S]*?<a.*?target="_blank">([\s\S]*?)</a>`)

func main() {
	plugin, c := log.NewFilePlugin("/dev/stdout", log.InfoLevel)
	defer c.Close()
	logger := log.NewLogger(plugin)
	logger.Info("log init end")

	var worklist []*collect.Request
	for i := 25; i <= 100; i += 25 {
		str := fmt.Sprintf("https://www.douban.com/group/szsh/discussion?start=%d", i)
		worklist = append(worklist, &collect.Request{
			Url:       str,
			ParseFunc: parse.ParseURL,
		})
	}
	var f collect.Fetcher = collect.BrowserFetch{
		Timeout: time.Second * 3000,
	}

	for len(worklist) > 0 {
		items := worklist
		worklist = nil
		for _, item := range items {
			body, err := f.Get(item)
			if err != nil {
				logger.Error("get err:", zap.Error(err))
				continue
			}
			time.Sleep(time.Second * 1)
			res := item.ParseFunc(body)
			for _, item := range res.Items {
				logger.Info("result", zap.String("get url:", item.(string)))
			}
			worklist = append(worklist, res.Requests...)
		}
	}

}
