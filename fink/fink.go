package fink

import (
	"context"
	"errors"
	"github.com/lifei6671/fink-download/internal/utils"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var DefaultUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.71 Safari/537.36 Edg/97.0.1072.55"

var (
	patternStr = `http[s]?://(?:[a-zA-Z]|[0-9]|[$-_@.&+]|[!*\(\),]|(?:%[0-9a-fA-F][0-9a-fA-F]))+`
)

type Downloader interface {
	Parser(ctx context.Context, urlStr string) (string, error)
	SaveFile(ctx context.Context, urlStr string, filename string) error
}

type downloadFink struct {
	c   http.Client
	exp *regexp.Regexp
}

func NewFinkDownload() Downloader {
	c := http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			MaxIdleConnsPerHost:   100,
			MaxConnsPerHost:       200,
		},
		Timeout: time.Second * 30,
	}

	return &downloadFink{
		c:   c,
		exp: regexp.MustCompile(patternStr),
	}
}

func (d *downloadFink) Parser(ctx context.Context, urlStr string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlStr, nil)
	if err != nil {
		log.Printf("创建请求失败: url[%s] errmsg[%+v]", urlStr, err)
		return "", err
	}
	req.Header.Set("User-Agent", DefaultUserAgent)
	req.Header.Set("Referer", urlStr)

	resp, err := d.c.Do(req)
	if err != nil {
		log.Printf("请求失败: url[%s] errmsg[%+v]", urlStr, err)
		return "", err
	}
	defer utils.SafeClose(resp.Body)

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Printf("解析返回值失败: url[%s] errmsg[%+v]", urlStr, err)
		return "", err
	}
	val, exist := doc.Find(".live-video .swiper-slide").First().Attr("style")
	if !exist {
		log.Printf("解析图片失败,属性不存在: url[%s]", urlStr)
		return "", err
	}
	imageStr := strings.TrimSuffix(d.exp.FindString(val), ")")
	if imageStr == "" {
		log.Printf("正则匹配图片失败: url[%s] value[%s]", urlStr, val)
		return "", errors.New("正则匹配图片失败：" + val)
	}
	uri, err := url.ParseRequestURI(imageStr)
	if err != nil {
		log.Printf("解析图片地址失败: url[%s]", urlStr)
		return "", err
	}
	uri.RawQuery = ""
	return uri.String(), nil
}

func (d *downloadFink) SaveFile(ctx context.Context, urlStr string, filename string) error {
	dir := filepath.Dir(filename)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0666)
		if err != nil {
			log.Printf("创建目录失败: url[%s] filename[%s] errmsg[%+v]", urlStr, filename, err)
			return err
		}
	}
	f, err := os.Create(filename)
	if err != nil {
		log.Printf("创建文件失败: url[%s] filename[%s] errmsg[%+v]", urlStr, filename, err)
		return err
	}
	defer utils.SafeClose(f)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlStr, nil)
	if err != nil {
		log.Printf("创建请求失败: url[%s] errmsg[%+v]", urlStr, err)
		return err
	}
	req.Header.Set("User-Agent", DefaultUserAgent)
	req.Header.Set("Referer", urlStr)

	resp, err := d.c.Do(req)
	if err != nil {
		log.Printf("请求失败: url[%s] errmsg[%+v]", urlStr, err)
		return err
	}
	defer utils.SafeClose(resp.Body)

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		log.Printf("保存文件失败: url[%s] errmsg[%+v]", urlStr, err)
		return err
	}
	return f.Sync()
}
