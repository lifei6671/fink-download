package fink

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"time"
)

var (
	ch     = make(chan string, 100)
	finker = NewFinkDownload()
)

func Push(urlStr string) {
	ch <- urlStr
}

func Run(ctx context.Context, dir string) error {
	for {
		select {
		case urlStr, ok := <-ch:
			if !ok {
				log.Println("管道已关闭!")
				return nil
			}
			imageStr, err := finker.Parser(ctx, urlStr)
			if err != nil {
				break
			}

			filename := filepath.Join(dir, "finkapp", fmt.Sprintf("fink-%d.jpg", time.Now().UnixNano()))

			err = finker.SaveFile(ctx, imageStr, filename)
			if err != nil {
				break
			}
			log.Printf("保存图片成功: %s", filename)
		case <-ctx.Done():
			return nil
		}
	}
}
