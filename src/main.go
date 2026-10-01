package main

import (
	"bypaws/src/web"
	"log"
	"net/http"
	_ "net/http/pprof"
)



func main() {
	go func() {
		log.Println("[pprof] 메모리 분석 서버 실행 중...")
		err := http.ListenAndServe(":6060", nil)
		if err != nil {
			log.Printf("[pprof] 분석 서버 오류: %v", err)
		}
	}()

	web.Run()
}