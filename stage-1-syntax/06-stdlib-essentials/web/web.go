// Package web 演示 net/http handler、client 与 httptest 友好的设计。
//
// import 路径：just-go/stage-1-syntax/06-stdlib-essentials/web
package web

import (
	"fmt"
	"io"
	"net/http"
)

// HelloHandler 返回一个基础 HTTP handler。
func HelloHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		if name == "" {
			name = "gopher"
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprintf(w, "hello, %s", name)
	})
}

// FetchText 使用传入的 client 请求 URL 并返回响应文本。
func FetchText(client *http.Client, url string) (string, error) {
	if client == nil {
		client = http.DefaultClient
	}
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
