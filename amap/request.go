package amap

import (
	"time"

	"github.com/go-resty/resty"
)

// SetResty SetResty
func SetResty() {
	resty.
		SetHeaders(map[string]string{
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/68.0.3440.84 Safari/537.36",
			// "User-Agent" = "Mozilla/5.0 (X11; U; Linux i686; en-US; rv:1.9.2.13) Gecko/20101206 Ubuntu/10.10 (maverick) Firefox/3.6.13",
			"Accept":          "*/*",
			"Accept-Language": "en-US,en;q=0.9,zh-CN;q=0.8,zh-TW;q=0.7,zh;q=0.6",
			// "Accept-Encoding":"gzip,deflate",
			"Accept-Charset":   "GB2312,utf-8;q=0.7,*;q=0.7",
			"Keep-Alive":       "115",
			"Connection":       "keep-alive",
			"X-Requested-With": "XMLHttpRequest",
		}).
		SetTimeout(10 * time.Second).
		SetRetryCount(5).
		SetRetryWaitTime(10 * time.Second).
		SetRetryMaxWaitTime(65 * time.Second)
}
