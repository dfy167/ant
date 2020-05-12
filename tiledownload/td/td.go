package td

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"gopkg.in/resty.v1"
)

var startTime time.Time
var processes = 12

const (
	// Mapbox Mapbox
	Mapbox string = "mapbox"

	// TdtNormalMap TdtNormalMap
	TdtNormalMap string = "tdt-normal-map"
	// TdtNormalAnnotion TdtNormalAnnotion
	TdtNormalAnnotion string = "tdt-normal-annotion"

	// TdtTerrainMap TdtTerrainMap
	TdtTerrainMap string = "tdt-terrain-map"
	// TdtTerrainAnnotion TdtTerrainAnnotion
	TdtTerrainAnnotion string = "tdt-terrain-annotion"
)

// Td td
func Td(baseDir, mapType, key string, maxZoom, minZoom int, bounds LngLatBounds) (err error) {

	// baseDir := "d:/tilemap/tdt"

	// key := "2ce94f67e58faa24beb7cb8a09780552"
	// mapType := "tdt-normal-map"
	// hunanLngLatBounds := LngLatBounds{LngLat{108, 24}, LngLat{115, 31}}

	for i := minZoom; i <= maxZoom; i++ {
		// err = hunanLngLatBounds.Download2FileFast(baseDir, key, i, mapType)
		startTime = time.Now()
		err = bounds.download2File(baseDir, key, i, mapType)
		// fmt.Println("zoomlevel:", i, " finished")
		fmt.Println()
		if err != nil {
			return
		}
	}
	return
}

//Tdfast Tdfast
func Tdfast(baseDir, mapType, key string, maxZoom, minZoom int, bounds LngLatBounds) (err error) {

	for i := minZoom; i <= maxZoom; i++ {
		startTime = time.Now()
		err = bounds.download2FileFast(baseDir, key, i, mapType)
		if err != nil {
			return
		}
	}
	return
}

func newClient() (r *resty.Client) {
	client := resty.New()
	client.SetRetryCount(3).SetRetryWaitTime(1 * time.Second)
	client.SetHeaders(map[string]string{
		"Accept":          "text/html, application/xhtml+xml, application/xml; q=0.9, */*; q=0.8",
		"Accept-Encoding": "gzip, deflate",
		"Accept-Language": "zh-Hans-CN, zh-Hans; q=0.8, en-US; q=0.5, en; q=0.3",
		"Cache-Control":   "max-age=0",
		"Connection":      "Keep-Alive",
		"Referer":         "http://map.tianditu.gov.cn/",
		"User-Agent":      userAgent(),
	})
	return client
}

type tileURL interface {
	Url(mapType string) (url string)
}

// LngLat LngLat
type LngLat struct {
	lng float64
	lat float64
}

// NewBounds NewBounds
func NewBounds(b1, b2, b3, b4 float64) (b LngLatBounds) {
	b = LngLatBounds{LngLat{b1, b2}, LngLat{b3, b4}}
	return
}

// LngLatBounds LngLatBounds
type LngLatBounds struct {
	Sw, Ne LngLat
}

// Tile Tile
type Tile struct {
	x, y, z int
}

func lngLat2Tile(lng, lat float64, z int) (x, y int) {
	x = int(math.Floor((lng + 180.0) / 360.0 * (math.Exp2(float64(z)))))
	y = int(math.Floor((1.0 - math.Log(math.Tan(lat*math.Pi/180.0)+1.0/math.Cos(lat*math.Pi/180.0))/math.Pi) / 2.0 * (math.Exp2(float64(z)))))
	if x < 0 {
		x = 0
	} else if x > int(math.Exp2(float64(z)))-1 {
		x = int(math.Exp2(float64(z))) - 1
	}
	if y < 0 {
		y = 0
	} else if y > int(math.Exp2(float64(z)))-1 {
		x = int(math.Exp2(float64(z))) - 1
	}
	return
}

func (b LngLatBounds) tiles(z int) (tiles []Tile) {
	minX, maxY := lngLat2Tile(b.Sw.lng, b.Sw.lat, z)
	maxX, minY := lngLat2Tile(b.Ne.lng, b.Ne.lat, z)

	for x := minX; x <= maxX; x++ {
		for y := minY; y <= maxY; y++ {
			tiles = append(tiles, Tile{x, y, z})
		}
	}
	return
}

func (t Tile) tile2URL(key, maptype string) string {
	switch maptype {
	case Mapbox:
		return fmt.Sprintf("https://api.mapbox.com/v4/mapbox.mapbox-streets-v7/%d/%d/%d.vector.pbf?access_token=%s", t.z, t.x, t.y, key)
	case TdtNormalMap:
		return fmt.Sprintf("http://t%d.tianditu.gov.cn/DataServer?T=vec_w&x=%d&y=%d&l=%d&tk=%s", rand.Intn(7), t.x, t.y, t.z, key)
	case TdtNormalAnnotion:
		return fmt.Sprintf("http://t%d.tianditu.gov.cn/DataServer?T=cva_w&x=%d&y=%d&l=%d&tk=%s", rand.Intn(7), t.x, t.y, t.z, key)
	case TdtTerrainMap:
		return fmt.Sprintf("http://t%d.tianditu.gov.cn/DataServer?T=img_w&x=%d&y=%d&l=%d&tk=%s", rand.Intn(7), t.x, t.y, t.z, key)
	case TdtTerrainAnnotion:
		return fmt.Sprintf("http://t%d.tianditu.gov.cn/DataServer?T=cia_w&x=%d&y=%d&l=%d&tk=%s", rand.Intn(7), t.x, t.y, t.z, key)
	}
	return ""
}

func (t Tile) tile2Filepath(base string) string {
	return filepath.Join(base, strconv.Itoa(t.z), strconv.Itoa(t.x), strconv.Itoa(t.y)+".png")
}

func (t Tile) download2File(baseDir, key, maptype string) error {
	waitTime := 5 * time.Second
	n := 1
	for {
		err := t.download2FileWait(baseDir, key, maptype)
		if err != nil {
			time.Sleep(waitTime * time.Duration(n))
			n += n
			continue
		} else {
			return nil
		}
	}
}

func (t Tile) download2FileWait(baseDir, key, maptype string) error {
	client := newClient()
	r, err := client.R().
		SetOutput(t.tile2Filepath(baseDir)).
		Get(t.tile2URL(key, maptype))
	if contentType := r.Header().Get("Content-Type"); contentType == "text/html" {
		return errAttack
	}
	return err
}

var errAttack = errors.New("get html")

func (b LngLatBounds) download2File(baseDir, key string, z int, maptype string) (err error) {
	tiles := b.tiles(z)
	l := len(tiles)

	for i, t := range tiles {
		err = t.download2File(baseDir, key, maptype)
		if err != nil {
			return
		}

		continuesPrint(z, l, i+1)

		// fmt.Println(i, len(tiles), z)
	}
	return
}
func (b LngLatBounds) download2FileFast(baseDir, key string, z int, maptype string) (err error) {
	tiles := b.tiles(z)
	l := len(tiles)
	var wg sync.WaitGroup
	ch := make(chan int, processes)
	for i, t := range tiles {
		ch <- 1
		wg.Add(1)
		go func(t Tile, i int) {
			defer wg.Done()
			err = t.download2File(baseDir, key, maptype)
			if err != nil {
				fmt.Println(t.tile2URL(key, maptype), err)
			}
			<-ch
			continuesPrint(z, l, i+1)
		}(t, i)

	}
	wg.Wait()
	return
}

// CheckKey CheckKey
func CheckKey(key, maptype string) error {
	t := Tile{3337, 1712, 12}
	client := newClient()
	r, err := client.R().
		Get(t.tile2URL(key, maptype))
	fmt.Println(t.tile2URL(key, maptype))
	if contentType := r.Header().Get("Content-Type"); contentType == "application/octet-stream" {
		return errors.New("the map api key is not valid")
	}
	return err
}

func timeDuC(d time.Duration, p float64) (td time.Duration) {
	return time.Duration(float64(d) * p)
}

func continuesPrint(z, tiles, i int) {
	passdTime := time.Since(startTime)
	if int64(passdTime) != 0 {
		fmt.Fprintf(os.Stdout, "Zoom: %d	G/T: %d/%d	Rate: %.4f%%	Used: %s	Speed(s):%.2f 	Remaining: %s	Expected: %s    \r",
			z, i, tiles, 100*float64(i)/float64(tiles), passdTime.Truncate(time.Second),
			1000000000*float64(i)/float64(passdTime),
			time.Duration(float64(tiles-i)/float64(i)*float64(passdTime)).Truncate(time.Second),
			time.Now().Add(time.Duration(float64(tiles-i)/float64(i)*float64(passdTime))).Format("2006-01-02 15:04"))
	} else {
		fmt.Fprintf(os.Stdout, "ZoomLevel: %d Total:	%d Finished: %d  Rate:	%.4f%% UsedTime: %s Speed: N/A tiles/second \r", z, tiles, i, 100*float64(i)/float64(tiles), passdTime)
	}
}

func userAgent() string {
	ua := []string{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.132 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.113 Safari/537.36",
		"Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; WOW64; Trident/5.0)",
		"Mozilla/5.0 (Windows NT 6.3; WOW64; rv:41.0) Gecko/20100101 Firefox/41.0",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 11_2_5 like Mac OS X) AppleWebKit/604.5.6 (KHTML, like Gecko) Version/11.0 Mobile/15D60 Safari/604.1",
		"Mozilla/5.0 (Linux; Android 6.0; MYA-L22 Build/HUAWEIMYA-L22) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.84 Mobile Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_5) AppleWebKit/602.2.14 (KHTML, like Gecko) Version/10.0.1 Safari/602.2.14",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 10_1 like Mac OS X) AppleWebKit/602.2.14 (KHTML, like Gecko) Version/10.0 Mobile/14B72 Safari/602.1",
		"Mozilla/5.0 (X11; Linux x86_64; rv:42.0) Gecko/20100101 Firefox/42.0",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 10_2 like Mac OS X) AppleWebKit/602.3.12 (KHTML, like Gecko) Version/10.0 Mobile/14C92 Safari/602.1",
		"Mozilla/5.0 (iPad; CPU OS 5_0_1 like Mac OS X) AppleWebKit/534.46 (KHTML, like Gecko) Version/5.1 Mobile/9A405 Safari/7534.48.3",
		"Mozilla/5.0 (Windows NT 10.0; WOW64; rv:54.0) Gecko/20100101 Firefox/54.0",
		"Mozilla/5.0 (Linux; Android 7.0; Redmi Note 4 Build/NRD90M) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.111 Mobile Safari/537.36",
	}
	return ua[rand.Intn(11)]
}

func download2File(URL, outPut string) error {
	client := newClient()
	r, err := client.R().
		SetOutput(outPut).
		Get(URL)
	if contentType := r.Header().Get("Content-Type"); contentType == "text/html" {
		return errAttack
	}
	return err
}
