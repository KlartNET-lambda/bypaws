package web

import (
	"bypaws/src/proxy"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

var whitelist = []string{
	"/favicon.ico",
	"/vite/assets/",
	"/images/counter/",

	"/tags/autocomplete.json",

	"/posts",
	"/popular",
}

func Run(targetURL *url.URL) {
	proxyURL, err := proxy.GetSingleProxyURL(targetURL)
	if err != nil {
		log.Fatalf("프록시 로드 실패: %v", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	proxy.Transport = &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
		MaxIdleConns: 5,
		MaxIdleConnsPerHost: 5,
		IdleConnTimeout: 60 * time.Second,
	}

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		reqURL := request.URL

		if reqURL.Path != "/" {
			isAllowed := false

			for _, allowedPath := range whitelist {
				if strings.HasPrefix(reqURL.Path, allowedPath) {
					isAllowed = true
					break
				}
			}

			if !isAllowed {
				http.Error(
					writer,
					"해당 페이지는 관리자에 의해 임시 차단되어 있습니다.",
					http.StatusForbidden,
				)
				return
			}
		}

		request.Host = targetURL.Host
		proxy.ServeHTTP(writer, request)
	})

	log.Printf("Server running on http://localhost:8080 via %s", proxyURL)
	log.Fatal(http.ListenAndServe(":8080", nil))
}