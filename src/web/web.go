package web

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

const staticHost = "static1.e621.net"
const staticPrefix = "/s1"

var whitelist = []string{
	"/favicon.ico",
	"/vite/assets/",
	"/images/counter/",

	"/tags/autocomplete.json",

	"/posts",
	"/popular",

	staticPrefix,
}

func Run(targetURL *url.URL) {
	transport := &http.Transport{
		MaxIdleConns: 5,
		MaxIdleConnsPerHost: 5,
		IdleConnTimeout: 5 * time.Second,
	}

	site := httputil.NewSingleHostReverseProxy(targetURL)
	site.Transport = transport
	site.ModifyResponse = rewriteStaticURLs

	static := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "https",
		Host: staticHost,
	})
	static.Transport = transport

	http.HandleFunc("/", func(writer http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/sw.js" {
			writer.Header().Set("Content-Type", "text/javascript; charset=utf-8")
			writer.Header().Set("Cache-Control", "no-cache")
			writer.Write([]byte(serviceWorkerJS))

			return
		}

		if req.URL.Path != "/" {
			ok := false

			for _, allowedPath := range whitelist {
				if strings.HasPrefix(req.URL.Path, allowedPath) {
					ok = true
					break
				}
			}

			if !ok {
				http.Error(
					writer,
					"해당 페이지는 관리자에 의해 임시 차단되어 있습니다.",
					http.StatusForbidden,
				)
				
				return
			}
		}
		
		after, ok := strings.CutPrefix(req.URL.Path, staticPrefix)
		if ok {
			req.URL.Path = after
			req.Host = staticHost
			static.ServeHTTP(writer, req)

			return
		}

		req.Header.Del("Accept-Encoding")
		req.Host = targetURL.Host
		site.ServeHTTP(writer, req)
	})

	log.Printf("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}