package fink

import (
	"context"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestDownloadFink_Parser(t *testing.T) {
	convey.Convey("DownloadFink_Parser", t, func() {
		downloader := NewFinkDownload()
		convey.Convey("DownloadFink_Parser_OK", func() {
			urlStr := "https://www.finkapp.cn/post/finka-tDyCfMcMYftPrCYDWVbkEQ?shareBy=PHvjIvbfDC0"

			s, err := downloader.Parser(context.Background(), urlStr)

			convey.So(err, convey.ShouldBeNil)
			t.Logf("%s", s)
		})
	})
}
