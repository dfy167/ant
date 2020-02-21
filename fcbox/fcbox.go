package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/go-resty/resty"
)

var client *resty.Client

type poiResult struct {
	Code       int     `json:"code"`
	EngDesc    string  `json:"engDesc"`
	ChnDesc    string  `json:"chnDesc"`
	Detail     string  `json:"detail"`
	Content    []poi   `json:"content"`
	TotalCount float64 `json:"totalCount"`
	PageSize   int     `json:"pageSize"`
	PageNo     int     `json:"pageNo"`
}

type poi struct {
	ServiceType       string  `json:"serviceType"`
	Address           string  `json:"address"`
	ServiceTime       string  `json:"serviceTime"`
	ServicePLointName string  `json:"servicePointName"`
	Longitude         float64 `json:"longitude"`
	Latitude          float64 `json:"latitude"`
	ImagePath         string  `jsIon:"imagePath"`
	ID                string  `json:"id"`
	ServiceScope      string  `json:"serviceScope"`
	Distance          string  `json:"distance"`
}

func (p poi) print() {
	fmt.Println(fmt.Sprintf("%s,%s,%s,%s,%f,%f", p.ServiceType, p.Address, p.ServiceTime, p.ServicePLointName, p.Longitude, p.Latitude))
}

func round(f float64) float64 {
	n10 := math.Pow10(6)
	return math.Trunc(f*n10) / n10
}

func ffloat64(f float64) string {
	return fmt.Sprintf("%f", f)
}

var fcboxDetailURL = "https://www.fcbox.com/serviceNodeQuery/nearServiceNode"

type point struct {
	longitude float64
	latitude  float64
}

func requestPagePois(p point, pageNo int) (result poiResult, err error) {
	para := map[string]string{
		"longitude": ffloat64(p.longitude),
		"latitude":  ffloat64(p.latitude),
		"pageNo":    strconv.Itoa(pageNo),
	}
	result, err = request(para)
	return
}

func request(para map[string]string) (result poiResult, err error) {
	getClient()
	resp, err := client.
		SetTimeout(10 * time.Second).
		SetRetryCount(5).
		SetRetryWaitTime(10 * time.Second).
		SetRetryMaxWaitTime(65 * time.Second).
		R().
		SetQueryParams(para).
		Get(fcboxDetailURL)
	if err != nil {
		return
	}
	json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return
	}
	if result.Code != 0 || result.EngDesc != "success" {
		err = fmt.Errorf("%d,%s", result.Code, result.ChnDesc)
		return
	}
	return
}

func getClient() {
	client = resty.New()
	setResty()
	username := "I3E31659135705650222"
	// 密码请到用户中心-我的订单页面查询
	password := "RcO5dXEXpj9Akm2F"
	proxyURL := fmt.Sprintf("http://%s:%s@dyn.horocn.com:50000", username, password)
	client.SetProxy(proxyURL)
}

func setResty() {
	client.
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
		SetRetryCount(2).
		SetRetryWaitTime(10 * time.Second).
		SetRetryMaxWaitTime(65 * time.Second)
}

func main() {
	//雨花区 https://www.fcbox.com/serviceNodeQuery/nearServiceNode?longitude=113.038017&latitude=28.13771&type=&pageNo=1&_=1582260373182
	// 天心区 https://www.fcbox.com/serviceNodeQuery/nearServiceNode?longitude=112.9962&latitude=28.14447&type=&pageNo=1&_=1582260334125
	// 芙蓉区 https://www.fcbox.com/serviceNodeQuery/nearServiceNode?longitude=113.032539&latitude=28.185386&type=&pageNo=2&_=1582254606574

	// 株洲市
	// 荷塘区 113.173487,27.855929
	// 芦淞区 113.152724,27.785070
	// 石峰区 113.117732,27.875445
	// 天元区  113.082216,27.826867
	// 株洲县 113.144006,27.699346
	// 攸县 113.396404,27.014607
	// 茶陵县 113.539280,26.777492
	// 炎陵县 113.772655,26.489902
	// 醴陵市 113.496894,27.646130
	// 禄口区 113.134002,27.827550
	// 云龙示范区 113.160642,27.903426

	// 湘潭市
	// 雨湖区 112.903317,27.854705
	// 岳塘区 112.925371,27.808646
	// 湘潭县 112.950781,27.778947
	// 湘乡市 112.550581,27.718313
	// 韶山市 112.526670,27.914958
	// points := []point{}

	var p point
	p = point{112.526670, 27.914958}

	maxPage := 40
	for page := 1; page <= maxPage; page++ {

		fcbox2, err2 := requestPagePois(p, page)
		if err2 != nil {
			log.Fatal(err2)
		}
		fmt.Println(page, int(math.Ceil(fcbox2.TotalCount/10)))
		maxPage = int(math.Ceil(fcbox2.TotalCount / 10))
		for _, p := range fcbox2.Content {
			p.print()
		}
	}

}
