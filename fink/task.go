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
			imageUrls, err := finker.Parser(ctx, urlStr)
			if err != nil {
				break
			}
			for i, imageUrl := range imageUrls {
				filename := filepath.Join(dir, "finkapp", fmt.Sprintf("fink-%d-%d.jpg", time.Now().UnixNano(), i))
				err = finker.SaveFile(ctx, imageUrl, filename)
				if err != nil {
					log.Printf("保存图片失败: %s - %+v", filename, err)
					continue
				}
				log.Printf("保存图片成功: %s", filename)
			}
		case <-ctx.Done():
			return nil
		}
	}
}
