package proxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func NewProxy(targetURL string) (*httputil.ReverseProxy, error) {
	url, err := url.Parse(targetURL)
	if err != nil {
		return nil, err
	}
	proxy := httputil.NewSingleHostReverseProxy(url)
	oridir := proxy.Director
	proxy.Director = func(request *http.Request) {
		oridir(request)
		request.Header.Add("ex-proxy", "veeya")
	}
	proxy.ModifyResponse = func(response *http.Response) error {
		response.Header.Add("ex-backend-server", url.Host)
		return nil
	}
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("proxy error: %v", err)
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte("502"))
	}
	return proxy, nil
}
