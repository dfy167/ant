package amap

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty"
)

// PoiResult PoiResult
type PoiResult struct {
	Count      string `json:"count"`
	Info       string `json:"info"`
	Infocode   string `json:"infocode"`
	Pois       []Poi  `json:"pois"`
	Status     string `json:"status"`
	Suggestion struct {
		Cities   []interface{} `json:"cities"`
		Keywords []interface{} `json:"keywords"`
	} `json:"suggestion"`
}

// Poi Poi
type Poi struct {
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

func (p Poi) String() string {
	return fmt.Sprintln(spaceD(p.ID), spaceD(p.Name), spaceD(p.Type), spaceD(p.Typecode), spaceD(p.Address), spaceD(p.Cityname), spaceD(p.Adname), spaceD(p.Location), spaceD(p.Alias))
}

func spaceD(s string) string {
	return strings.Join(strings.Fields(s), "")
}

// Point Point
type Point struct {
	Lng float64
	Lat float64
}

// Rectangle Rectangle
type Rectangle struct {
	PointLT Point
	PointRB Point
}

func (r Rectangle) check() bool {
	return r.PointLT.Lng < r.PointRB.Lng && r.PointLT.Lat > r.PointRB.Lat
}

func (r Rectangle) polygon() string {
	return fmt.Sprintf("%f,%f|%f,%f", r.PointLT.Lng, r.PointLT.Lat, r.PointRB.Lng, r.PointRB.Lat)
}

func (r Rectangle) quadtree() []Rectangle {
	halflng, halflat := math.Abs(r.PointRB.Lng-r.PointLT.Lng)/2, math.Abs(r.PointLT.Lat-r.PointRB.Lat)/2

	return []Rectangle{
		{r.PointLT, Point{round(r.PointLT.Lng + halflng), round(r.PointLT.Lat - halflat)}},
		{Point{round(r.PointLT.Lng + halflng), r.PointLT.Lat}, Point{r.PointRB.Lng, round(r.PointLT.Lat - halflat)}},
		{Point{r.PointLT.Lng, round(r.PointLT.Lat - halflat)}, Point{round(r.PointLT.Lng + halflng), r.PointRB.Lat}},
		{Point{round(r.PointLT.Lng + halflng), round(r.PointLT.Lat - halflat)}, r.PointRB}}
}

type minRec struct {
	Rec   Rectangle
	Types string
	Count int
	Err   error
}

type minRecPage struct {
	Rec   Rectangle
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

func cutRec(rec Rectangle, types string) (recCutresult []minRec) {
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

func recCount(rec Rectangle, types string) (count int, err error) {
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

func minRecPagePois(minRecPage minRecPage) (pois []Poi, err error) {
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

func minRecPagesPois(minRecPages []minRecPage) (pois []Poi) {
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

func recTypePages(rec Rectangle, types string) (mrp []minRecPage) {
	cutrec := cutRec(rec, types)
	mrp = minRecsPages(cutrec)
	return
}

// RecTypePois RecTypePois
func RecTypePois(rec Rectangle, types string) (pois []Poi) {
	pages := recTypePages(rec, types)
	pois = minRecPagesPois(pages)
	return
}

func recRequest(para map[string]string) (result PoiResult, err error) {
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

// Detail Detail
type Detail struct {
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

func requestDetail(id string) (result Detail, err error) {
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

func requestDetails(ids []string) (result []Detail) {
	for _, id := range ids {
		r, err1 := requestDetail(id)
		if err1 == nil {
			result = append(result, r)
		}
	}
	return
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
