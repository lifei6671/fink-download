package fink

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"path/filepath"
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
			for _, imageUrl := range imageUrls {
				filename := filepath.Join(dir, "finkapp", fmt.Sprintf("fink-%s.jpg", createFileName(imageUrl)))
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

func createFileName(urlStr string) string {
	h := md5.New()
	h.Write([]byte(urlStr))
	return hex.EncodeToString(h.Sum(nil))
}
