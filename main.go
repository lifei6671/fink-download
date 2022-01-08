package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/lifei6671/fink-download/fink"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var (
	addr string
	path string
)

func main() {
	flag.StringVar(&addr, "addr", ":8089", "监听地址")
	flag.StringVar(&path, "path", "./images/", "保存图片路径")
	flag.Parse()
	if addrStr := os.Getenv("addr"); addrStr != "" {
		addr = addrStr
	}
	if pathStr := os.Getenv("path"); pathStr != "" {
		path = pathStr
	}
	var err error
	path, err = filepath.Abs(path)
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			log.Printf("解析请求失败: %+v", err)
			_, _ = fmt.Fprintf(w, "failed")
			return
		}
		content := r.FormValue("content")
		if content == "" {
			log.Println("解析请求参数失败")
			_, _ = fmt.Fprintf(w, "failed")
			return
		}
		fink.Push(content)
		_, _ = fmt.Fprintf(w, "success")
	})
	go func() {
		if err := fink.Run(context.Background(), path); err != nil {
			panic(err)
		}
	}()
	log.Printf("服务已启动: http://%s", addr)
	err = http.ListenAndServe(addr, nil)
	if err != nil {
		log.Printf("启动服务失败: addr[%s] %+v", addr, err)
	}
}
