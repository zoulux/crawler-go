package main

import (
	"fmt"
	"regexp"
	"time"

	"github.com/zoulux/crawler-go/collect"
	"github.com/zoulux/crawler-go/proxy"
)

var headerRe = regexp.MustCompile(`<div class="news_li"[\s\S]*?<h2>[\s\S]*?<a.*?target="_blank">([\s\S]*?)</a>`)

func main() {

	proxyURLs := []string{"http://127.0.0.1:8088"}

	p, err := proxy.RoundRobinProxySwitcher(proxyURLs...)
	url := "https://book.douban.com/subject/1007305/"
	var f collect.Fetcher = collect.BrowserFetch{
		Timeout: 3000 * time.Second,
		Proxy:   p,
	}
	body, err := f.Get(url)
	if err != nil {
		fmt.Println("read content failed:%v", err)
		return
	}

	fmt.Println(string(body))
}
