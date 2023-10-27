package doubangroup

import (
	"regexp"

	"github.com/zoulux/crawler-go/collect"
)

const cityListRe = `(https://www.douban.com/group/topic/[0-9a-z]+/)"[^>]*>([^<]+)</a>`
const contentRe = `<div class="topic-content">[\s\S]*?阳台[\s\S]*?<div`

func ParseURL(contents []byte, req *collect.Request) collect.ParseResult {
	var re = regexp.MustCompile(cityListRe)

	matches := re.FindAllSubmatch(contents, -1)
	result := collect.ParseResult{}
	for _, m := range matches {
		u := string(m[1])
		result.Requests = append(result.Requests, &collect.Request{
			Url:    u,
			Cookie: req.Cookie,
			ParseFunc: func(c []byte, req *collect.Request) collect.ParseResult {
				return GetContent(c, u)
			},
		})
	}
	return result
}

func GetContent(contents []byte, u string) collect.ParseResult {
	var re = regexp.MustCompile(contentRe)
	if !re.Match(contents) {
		return collect.ParseResult{
			Items: []interface{}{},
		}
	}
	return collect.ParseResult{
		Items: []interface{}{u},
	}
}
