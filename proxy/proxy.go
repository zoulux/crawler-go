package proxy

import (
	"errors"
	"net/http"
	"net/url"
	"sync/atomic"
)

type ProxyFunc func(req *http.Request) (*url.URL, error)

func RoundRobinProxySwitcher(proxyUrls ...string) (ProxyFunc, error) {

	if len(proxyUrls) < 1 {
		return nil, errors.New("proxyUrls length must be greater than 0")
	}
	urls := make([]*url.URL, len(proxyUrls))
	for i, u := range proxyUrls {
		parsedu, err := url.Parse(u)
		if err != nil {
			return nil, err
		}
		urls[i] = parsedu
	}
	return (&roundRobinProxySwitcher{
		proxyURLs: urls,
		index:     0,
	}).GetProxy, nil
}

type roundRobinProxySwitcher struct {
	proxyURLs []*url.URL
	index     uint32
}

func (r *roundRobinProxySwitcher) GetProxy(pr *http.Request) (*url.URL, error) {
	idx := atomic.AddUint32(&r.index, 1) - 1
	u := r.proxyURLs[idx%uint32(len(r.proxyURLs))]
	return u, nil
}
