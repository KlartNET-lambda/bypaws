package main

import (
	"io"
	"net/http"
	"strings"
)

func rewriter(res *http.Response) error {
	body, err := io.ReadAll(io.LimitReader(res.Body, 10<<20))
	res.Body.Close()
	if err != nil {
		return err
	}

	out := string(body)
	out = strings.Replace(out, "<head>", `<head><script>navigator.serviceWorker.register("/sw.js")</script>`, 1)

	res.Header.Del("Content-Security-Policy")
	res.Header.Del("Content-Length")
	res.ContentLength = -1
	res.Body = io.NopCloser(strings.NewReader(out))

	return nil
}