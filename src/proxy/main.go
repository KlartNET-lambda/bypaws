package proxy

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const proxyListURL = "https://raw.githubusercontent.com/proxifly/free-proxy-list/refs/heads/main/proxies/countries/JP/data.txt"

func testProxyConnection(proxyURL *url.URL) bool {
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest("HEAD", "https://e621.net/", nil)
	if err != nil {
		return false
	}

	res, err := client.Do(req)
	if err != nil {
		return false
	}
	defer res.Body.Close()

	return res.StatusCode > 0
}

func GetSingleProxy() (*url.URL, error) {
	response, err := http.Get(proxyListURL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	scanner := bufio.NewScanner(response.Body)
	log.Println("========== 프록시 유효성 검증 시작 ==========")
	for scanner.Scan() {
		rawURL := strings.TrimSpace(scanner.Text())

		parsedURL, err := url.Parse(rawURL)
		if err != nil {
			continue
		}

		log.Printf("검증 중... (%s)", parsedURL.String())
		if testProxyConnection(parsedURL) {
			log.Println("검증 성공")
			log.Println("========== 프록시 유효성 검증 종료 ==========")
			return parsedURL, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return nil, fmt.Errorf("사용 가능한 SOCKS5 프록시가 목록에 없음")
}