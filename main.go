package main

import (
	"bytes"
	"fmt"
	"regexp"

	"github.com/PuerkitoBio/goquery"

	"github.com/zoulux/crawler-go/collect"
)

var headerRe = regexp.MustCompile(`<div class="news_li"[\s\S]*?<h2>[\s\S]*?<a.*?target="_blank">([\s\S]*?)</a>`)

func main() {
	url := "https://book.douban.com/subject/1007305/"
	var f collect.Fetcher = collect.BaseFetch{}
	body, err := f.Get(url)
	if err != nil {
		fmt.Println("read content failed:%v", err)
		return
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		fmt.Println("parse html failed:%v", err)
		return
	}

	doc.Find("div.news_li h2 a[target-_blank]").Each(func(i int, s *goquery.Selection) {
		title := s.Text()
		fmt.Printf("Review %d: %s \n", i, title)
	})

}
