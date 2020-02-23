package main

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty"
)

type poiResult struct {
	Count      string `json:"count"`
	Info       string `json:"info"`
	Infocode   string `json:"infocode"`
	Pois       []poi  `json:"pois"`
	Status     string `json:"status"`
	Suggestion struct {
		Cities   []interface{} `json:"cities"`
		Keywords []interface{} `json:"keywords"`
	} `json:"suggestion"`
}

type poi struct {
	Adcode  string `json:"adcode"`
	Address string `json:"address"`
	Adname  string `json:"adname"`
	Alias   string `json:"alias"`
	BizExt  struct {
		Cost   string `json:"cost"`
		Rating string `json:"rating"`
	} `json:"biz_ext"`
	BizType      string        `json:"biz_type"`
	BusinessArea string        `json:"business_area"`
	Children     []interface{} `json:"children"`
	Citycode     string        `json:"citycode"`
	Cityname     string        `json:"cityname"`
	DiscountNum  string        `json:"discount_num"`
	Distance     string        `json:"distance"`
	Email        string        `json:"email"`
	EntrLocation string        `json:"entr_location"`
	Event        []interface{} `json:"event"`
	ExitLocation []interface{} `json:"exit_location"`
	Gridcode     string        `json:"gridcode"`
	GroupbuyNum  string        `json:"groupbuy_num"`
	ID           string        `json:"id"`
	Importance   []interface{} `json:"importance"`
	IndoorData   struct {
		Cmsid     []interface{} `json:"cmsid"`
		Cpid      []interface{} `json:"cpid"`
		Floor     []interface{} `json:"floor"`
		Truefloor []interface{} `json:"truefloor"`
	} `json:"indoor_data"`
	IndoorMap string `json:"indoor_map"`
	Location  string `json:"location"`
	Match     string `json:"match"`
	Name      string `json:"name"`
	NaviPoiid string `json:"navi_poiid"`
	Pcode     string `json:"pcode"`
	Photos    []struct {
		Title []interface{} `json:"title"`
		URL   string        `json:"url"`
	} `json:"photos"`
	Pname     string        `json:"pname"`
	Poiweight []interface{} `json:"poiweight"`
	Postcode  []interface{} `json:"postcode"`
	Recommend string        `json:"recommend"`
	Shopid    []interface{} `json:"shopid"`
	Shopinfo  string        `json:"shopinfo"`
	Tag       []interface{} `json:"tag"`
	Tel       string        `json:"tel"`
	Timestamp []interface{} `json:"timestamp"`
	Type      string        `json:"type"`
	Typecode  string        `json:"typecode"`
	Website   []interface{} `json:"website"`
}

func (p poi) String() string {
	return fmt.Sprintln(spaceD(p.ID), spaceD(p.Name), spaceD(p.Type), spaceD(p.Typecode), spaceD(p.Address), spaceD(p.Cityname), spaceD(p.Adname), spaceD(p.Location), spaceD(p.Alias))
}

func spaceD(s string) string {
	return strings.Join(strings.Fields(s), "")
}

type point struct {
	Lng float64
	Lat float64
}

type rectangle struct {
	PointLT point
	PointRB point
}

func (r rectangle) check() bool {
	return r.PointLT.Lng < r.PointRB.Lng && r.PointLT.Lat > r.PointRB.Lat
}

func (r rectangle) polygon() string {
	return fmt.Sprintf("%f,%f|%f,%f", r.PointLT.Lng, r.PointLT.Lat, r.PointRB.Lng, r.PointRB.Lat)
}

func (r rectangle) quadtree() []rectangle {
	halflng, halflat := math.Abs(r.PointRB.Lng-r.PointLT.Lng)/2, math.Abs(r.PointLT.Lat-r.PointRB.Lat)/2

	return []rectangle{
		rectangle{r.PointLT, point{round(r.PointLT.Lng + halflng), round(r.PointLT.Lat - halflat)}},
		rectangle{point{round(r.PointLT.Lng + halflng), r.PointLT.Lat}, point{r.PointRB.Lng, round(r.PointLT.Lat - halflat)}},
		rectangle{point{r.PointLT.Lng, round(r.PointLT.Lat - halflat)}, point{round(r.PointLT.Lng + halflng), r.PointRB.Lat}},
		rectangle{point{round(r.PointLT.Lng + halflng), round(r.PointLT.Lat - halflat)}, r.PointRB}}
}

type minRec struct {
	Rec   rectangle
	Types string
	Count int
	Err   error
}

type minRecPage struct {
	Rec   rectangle
	Types string
	Page  string
}

func round(f float64) float64 {
	n10 := math.Pow10(6)
	return math.Trunc(f*n10) / n10
}

var gaoDePolygonURL = "https://restapi.amap.com/v3/place/polygon"
var gaoDeDetailURL = "https://www.amap.com/detail/get/detail"

var key = "aaa8abdaf05433e3702eae99964cc8c6"

// var key = "935c7385f239000f98ade53bbbc002e7"

func cutRec(rec rectangle, types string) (recCutresult []minRec) {
	count, err := recCount(rec, types)
	if err != nil {
		fmt.Println(rec, types, count, err)
		recCutresult = append(recCutresult, minRec{rec, types, count, err})
	} else if count <= 800 && count > 0 {
		fmt.Println(rec, types, count, err)
		recCutresult = append(recCutresult, minRec{rec, types, count, err})
	} else if count > 800 {
		// fmt.Println("cuting:", rec, types, count, err)
		rec4s := rec.quadtree()
		for _, rec4 := range rec4s {
			recCutresult = append(recCutresult, cutRec(rec4, types)...)
		}
	}
	return
}

func recCount(rec rectangle, types string) (count int, err error) {
	para := map[string]string{
		"types":   types,
		"offset":  "1",
		"polygon": rec.polygon(),
	}
	poiResult1, err := recRequest(para)
	if err != nil {
		return
	}
	count, err = strconv.Atoi(poiResult1.Count)
	if err != nil {
		return
	}
	return
}

func minRecPagePois(minRecPage minRecPage) (pois []poi, err error) {
	para := map[string]string{
		"types":   minRecPage.Types,
		"offset":  "20",
		"polygon": minRecPage.Rec.polygon(),
		"page":    minRecPage.Page,
	}
	result, err := recRequest(para)
	if err != nil {
		return
	}
	pois = result.Pois
	return
}

func minRecPagesPois(minRecPages []minRecPage) (pois []poi) {
	for _, minRecPage := range minRecPages {
		pagePois, err := minRecPagePois(minRecPage)
		if err == nil {
			pois = append(pois, pagePois...)
		} else {
			fmt.Println(minRecPages, err)
		}
	}
	return
}

func minRecPages(mRec minRec) (minRecPages []minRecPage) {
	for page := int(math.Ceil(float64(mRec.Count) / 20)); page > 0; page-- {
		minRecPages = append(minRecPages, minRecPage{mRec.Rec, mRec.Types, strconv.Itoa(page)})
	}
	return
}

func minRecsPages(mRecs []minRec) (mrp []minRecPage) {
	for _, mRec := range mRecs {
		mrp = append(mrp, minRecPages(mRec)...)
	}
	return
}

func recTypePages(rec rectangle, types string) (mrp []minRecPage) {
	cutrec := cutRec(rec, types)
	mrp = minRecsPages(cutrec)
	return
}

func recTypePois(rec rectangle, types string) (pois []poi) {
	pages := recTypePages(rec, types)
	pois = minRecPagesPois(pages)
	return
}

func recRequest(para map[string]string) (result poiResult, err error) {
	para["key"] = key
	resp, err := resty.
		SetTimeout(10 * time.Second).
		SetRetryCount(5).
		SetRetryWaitTime(10 * time.Second).
		SetRetryMaxWaitTime(65 * time.Second).
		R().
		SetQueryParams(para).
		Get(gaoDePolygonURL)
	if err != nil {
		return
	}
	json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return
	}
	if result.Status != "1" || result.Infocode != "10000" {
		err = fmt.Errorf(result.Status, result.Infocode, result.Info)
		return
	}
	return
}

type detail struct {
	Status string `json:"status"`
	Data   struct {
		Base struct {
			PoiTag            string `json:"poi_tag"`
			Code              string `json:"code"`
			ImportanceVipFlag int    `json:"importance_vip_flag"`
			CityAdcode        string `json:"city_adcode"`
			Telephone         string `json:"telephone"`
			NewType           string `json:"new_type"`
			CityName          string `json:"city_name"`
			NewKeytype        string `json:"new_keytype"`
			Checked           string `json:"checked"`
			Title             string `json:"title"`
			CreFlag           int    `json:"cre_flag"`
			StdTTag0V         string `json:"std_t_tag_0_v"`
			NaviGeometry      string `json:"navi_geometry"`
			Classify          string `json:"classify"`
			Business          string `json:"business"`
			ShopInfo          struct {
				Claim int `json:"claim"`
			} `json:"shop_info"`
			PoiTagHasTTag int    `json:"poi_tag_has_t_tag"`
			Pixelx        string `json:"pixelx"`
			Pixely        string `json:"pixely"`
			Geodata       struct {
				Aoi []struct {
					Name    string  `json:"name"`
					Mainpoi string  `json:"mainpoi"`
					Area    float64 `json:"area"`
				} `json:"aoi"`
			} `json:"geodata"`
			Poiid           string `json:"poiid"`
			Distance        int    `json:"distance"`
			Name            string `json:"name"`
			StdVTag0V       string `json:"std_v_tag_0_v"`
			EndPoiExtension string `json:"end_poi_extension"`
			Y               string `json:"y"`
			X               string `json:"x"`
			Address         string `json:"address"`
			Bcs             string `json:"bcs"`
			Tag             string `json:"tag"`
		} `json:"base"`
		Spec struct {
			MiningShape struct {
				Aoiid  string `json:"aoiid"`
				Center string `json:"center"`
				Level  int    `json:"level"`
				SpType string `json:"sp_type"`
				Area   string `json:"area"`
				Shape  string `json:"shape"`
				Type   int    `json:"type"`
			} `json:"mining_shape"`
			SpPic []interface{} `json:"sp_pic"`
		} `json:"spec"`
		Residential struct {
			BuildingTypes   string        `json:"building_types"`
			SrcTypeMix      string        `json:"src_type_mix"`
			SrcID           string        `json:"src_id"`
			IsCommunity     int           `json:"is_community"`
			Business        string        `json:"business"`
			Price           string        `json:"price"`
			HaveSchDistrict int           `json:"have_sch_district"`
			PropertyFee     string        `json:"property_fee"`
			AreaTotal       string        `json:"area_total"`
			PropertyCompany string        `json:"property_company"`
			VolumeRate      float64       `json:"volume_rate"`
			GreenRate       string        `json:"green_rate"`
			SrcType         string        `json:"src_type"`
			Intro           string        `json:"intro"`
			HxpicInfo       []interface{} `json:"hxpic_info"`
			Developer       string        `json:"developer"`
		} `json:"residential"`
		Deep struct {
			BuildingTypes   string        `json:"building_types"`
			SrcTypeMix      string        `json:"src_type_mix"`
			SrcID           string        `json:"src_id"`
			IsCommunity     int           `json:"is_community"`
			Business        string        `json:"business"`
			Price           string        `json:"price"`
			HaveSchDistrict int           `json:"have_sch_district"`
			PropertyFee     string        `json:"property_fee"`
			AreaTotal       string        `json:"area_total"`
			PropertyCompany string        `json:"property_company"`
			VolumeRate      float64       `json:"volume_rate"`
			GreenRate       string        `json:"green_rate"`
			SrcType         string        `json:"src_type"`
			Intro           string        `json:"intro"`
			HxpicInfo       []interface{} `json:"hxpic_info"`
			Developer       string        `json:"developer"`
		} `json:"deep"`
		Rti struct {
			ReviewEntrance  int           `json:"review_entrance"`
			ReviewSummary   string        `json:"review_summary"`
			ReviewCount     int           `json:"review_count"`
			HasDiscountFlag int           `json:"has_discount_flag"`
			ReviewLabels    []interface{} `json:"review_labels"`
		} `json:"rti"`
		Review struct {
			Comment []struct {
				AosTagScore      float64       `json:"aos_tag_score"`
				Recommend        string        `json:"recommend"`
				HighQuality      int           `json:"high_quality"`
				Labels           []interface{} `json:"labels"`
				ReviewID         string        `json:"review_id"`
				AuthorProfileurl string        `json:"author_profileurl"`
				ReviewWeburl     string        `json:"review_weburl"`
				ReviewWapurl     string        `json:"review_wapurl"`
				Review           string        `json:"review"`
				Author           string        `json:"author"`
				GoldNum          int           `json:"gold_num"`
				QualityFlag      int           `json:"quality_flag"`
				GoldType         string        `json:"gold_type"`
				Score            int           `json:"score"`
				LikeNum          string        `json:"like_num"`
				ReviewAppurl     struct {
					IosAppurl     string `json:"ios_appurl"`
					AndroidAppurl string `json:"android_appurl"`
				} `json:"review_appurl"`
				Time     string `json:"time"`
				SrcName  string `json:"src_name"`
				SrcType  string `json:"src_type"`
				AuthorID int    `json:"author_id"`
			} `json:"comment"`
		} `json:"review"`
		SrcInfo  []interface{} `json:"src_info"`
		ShareURL string        `json:"share_url"`
	} `json:"data"`
}

func setResty() {
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

func requestDetail(id string) (result detail, err error) {
	resp, err := resty.
		R().
		SetQueryParams(map[string]string{"id": id}).
		Get(gaoDeDetailURL)
	if err != nil {
		return
	}
	json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return
	}
	if result.Status != "1" {
		err = fmt.Errorf(id, result.Status)
		return
	}
	return
}

func requestDetails(ids []string) (result []detail) {
	for _, id := range ids {
		r, err1 := requestDetail(id)
		if err1 == nil {
			result = append(result, r)
		}
	}
	return
}

func main() {
	setResty()
	recChangSha := rectangle{point{111.89, 28.66}, point{114.24, 27.85}}
	// recChangSha = rectangle{point{111.89, 28.66}, point{112.00, 28.00}}
	// recChangSha = rectangle{point{108, 31}, point{115, 24}}
	// 长沙株洲湘潭
	recChangSha = rectangle{point{111.89, 28.66}, point{114.24, 26.14}}

	types := []string{}
	// types = []string{"010000", "020000", "030000", "040000", "050000", "060000", "070000", "080000", "090000", "100000", "110000", "120000", "130000", "140000", "150000", "160000", "170000", "180000", "190000", "200000", "220000", "970000", "990000"}
	// types = []string{"110000", "120000", "130000", "140000", "150000"}
	// types = []string{"010000", "020000", "030000", "040000", "050000", "060000", "070000", "080000", "090000", "100000"}
	// types = []string{"160000", "170000", "180000", "190000", "200000", "220000", "970000", "990000"}
	// types = []string{"010000"}
	types = []string{"070000"}
	var pois []poi
	for _, typess := range types {
		pois = append(pois, recTypePois(recChangSha, typess)...)
	}
	fmt.Println(pois)

	// maxRoutineNum := 1
	// ch := make(chan string, maxRoutineNum)
	// fi, err := os.Open("id.txt")
	// if err != nil {
	// 	fmt.Printf("Error: %s\n", err)
	// 	return
	// }
	// defer fi.Close()
	// var ids []string
	// br := bufio.NewReader(fi)
	// for {
	// 	a, _, c := br.ReadLine()
	// 	if c == io.EOF {
	// 		break
	// 	}
	// 	ids = append(ids, string(a))
	// }
	// for _, id := range ids {
	// 	ch <- id
	// 	go printResult(id, ch)
	// }
	// time.Sleep(10 * time.Second)

	// fi, err := os.Open("id.txt")
	// if err != nil {
	// 	fmt.Printf("Error: %s\n", err)
	// 	return
	// }
	// defer fi.Close()
	// var ids []string
	// br := bufio.NewReader(fi)
	// for {
	// 	a, _, c := br.ReadLine()
	// 	if c == io.EOF {
	// 		break
	// 	}
	// 	ids = append(ids, string(a))
	// }
	// for _, id := range ids {
	// 	r, err := requestDetail(id)
	// 	if err == nil {
	// 		fmt.Println(id, r.Data.Spec.MiningShape.Shape, "type:"+strconv.Itoa(r.Data.Spec.MiningShape.Type), "sptype:"+r.Data.Spec.MiningShape.SpType)
	// 	} else if r.Status == "6" {
	// 		fmt.Println(id, "err:toofast")
	// 		break

	// 	} else if r.Status == "8" {
	// 		fmt.Println(id, "err:notfounddetail")
	// 	} else {
	// 		fmt.Println(id, "err"+r.Status, err)
	// 		break
	// 	}
	// }

}

func printResult(id string, ch chan string) {
	r, err := requestDetail(id)
	if err == nil {
		fmt.Println(id, r.Data.Spec.MiningShape.Shape, "type:"+strconv.Itoa(r.Data.Spec.MiningShape.Type), "sptype:"+r.Data.Spec.MiningShape.SpType)
	} else if r.Status == "6" {
		fmt.Println(id, "err:toofast")
		time.Sleep(10 * time.Second)

	} else if r.Status == "8" {
		fmt.Println(id, "err:notfounddetail")
	} else {
		fmt.Println(id, "err"+r.Status)
		time.Sleep(10 * time.Second)
	}
	<-ch
}
