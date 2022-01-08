package utils

import (
	"io"
	"log"
)

// SafeClose 安全关闭
func SafeClose(closer io.Closer) {
	if closer != nil {
		if err := closer.Close(); err != nil {
			log.Printf("关闭失败: errmsg[%+v]", err)
		}
	}
}
