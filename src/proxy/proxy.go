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

const proxyListAddr = "https://raw.githubusercontent.com/proxifly/free-proxy-list/refs/heads/main/proxies/countries/JP/data.txt"

func testProxyConn(proxyURL *url.URL, targetURL *url.URL) bool {
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
			MaxIdleConns: 5,
			MaxIdleConnsPerHost: 5,
			IdleConnTimeout: 15 * time.Second,
		},
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest("HEAD", targetURL.String(), nil)
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

func GetSingleProxyURL(targetURL *url.URL) (*url.URL, error) {
	response, err := http.Get(proxyListAddr)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	log.Println("========== 프록시 유효성 검증 시작 ==========")

	scanner := bufio.NewScanner(response.Body)
	for scanner.Scan() {
		proxyURL, err := url.Parse(strings.TrimSpace(scanner.Text()))
		if err != nil {
			continue
		}

		log.Printf("검증 중... (%s)", proxyURL.String())
		if testProxyConn(proxyURL, targetURL) {
			log.Println("검증 성공")
			log.Println("========== 프록시 유효성 검증 종료 ==========")

			return proxyURL, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return nil, fmt.Errorf("사용 가능한 SOCKS5 프록시 없음")
}