package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

const E621_ADDR = "https://e621.net/"
const CDN_HOST = "https://static1.e621.net/"

var originURL, _ = url.Parse(E621_ADDR)
var staticURL, _ = url.Parse(CDN_HOST)
var transport = &http.Transport{
	MaxIdleConns: 5,
	MaxIdleConnsPerHost: 5,
	IdleConnTimeout: 5 * time.Second,
}

func main() {
	proxyer := httputil.NewSingleHostReverseProxy(originURL)
	proxyer.Transport = transport
	proxyer.ModifyResponse = rewriter

	staticProxyer := httputil.NewSingleHostReverseProxy(staticURL)
	staticProxyer.Transport = transport

	http.HandleFunc("/cdn/", func(writer http.ResponseWriter, req *http.Request) {
		after := strings.TrimPrefix(req.URL.Path, "/cdn")

		req.Host = staticURL.Host
		req.URL.Path = after
		staticProxyer.ServeHTTP(writer, req)
	})
	http.HandleFunc("/", func(writer http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/sw.js" {
			writer.Header().Set("Content-Type", "text/javascript; charset=utf-8")
			writer.Header().Set("Cache-Control", "no-cache")
			writer.Write([]byte(serviceWorkerJS))

			return
		}
		
		req.Header.Del("Accept-Encoding")
		req.Header.Set("Origin", originURL.Scheme+"://"+originURL.Host)
		req.Host = originURL.Host
		proxyer.ServeHTTP(writer, req)
	})

	log.Printf("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}