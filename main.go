package main

import (
	"fmt"
	"regexp"

	"github.com/zoulux/crawler-go/collect"
)

var headerRe = regexp.MustCompile(`<div class="news_li"[\s\S]*?<h2>[\s\S]*?<a.*?target="_blank">([\s\S]*?)</a>`)

func main() {
	url := "https://book.douban.com/subject/1007305/"
	var f collect.Fetcher = collect.BrowserFetch{}
	body, err := f.Get(url)
	if err != nil {
		fmt.Println("read content failed:%v", err)
		return
	}

	fmt.Println(string(body))
}
