package main

import (
	"log"

	_ "github.com/lib/pq"

	"github.com/dfy167/ant/amap"
	"xorm.io/xorm"
)

func main() {
	var err error
	// 数据库
	// engine, err := xorm.NewEngine("sqlite3", "./mimo.db")
	engine, err := xorm.NewEngine("postgres", "postgres://postgres:3@localhost:5432/sp?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	h := amap.NewHandler(engine)

	err = h.SyncDb()
	if err != nil {
		log.Fatal(err)
	}
	amap.SetResty()
	recChangSha := amap.Rectangle{PointLT: amap.Point{Lng: 111.89, Lat: 28.66}, PointRB: amap.Point{Lng: 114.24, Lat: 27.85}}
	// recChangSha = rectangle{point{111.89, 28.66}, point{112.00, 28.00}}
	// recChangSha = rectangle{point{108, 31}, point{115, 24}}
	// 长沙株洲湘潭
	// recChangSha = amap.Rectangle{amap.Point{111.89, 28.66}, amap.Point{114.24, 26.14}}
	// test
	// recChangSha = amap.Rectangle{PointLT: amap.Point{Lng: 112.89, Lat: 28.66}, PointRB: amap.Point{Lng: 112.90, Lat: 28.25}}

	types := []string{}
	types = []string{"010000", "020000", "030000", "040000", "050000", "060000", "070000", "080000", "090000", "100000", "110000", "120000", "130000", "140000", "150000", "160000", "170000", "180000", "190000", "200000", "220000", "970000", "990000"}
	// types = []string{"110000", "120000", "130000", "140000", "150000"}
	// types = []string{"010000", "020000", "030000", "040000", "050000", "060000", "070000", "080000", "090000", "100000"}
	// types = []string{"160000", "170000", "180000", "190000", "200000", "220000", "970000", "990000"}
	// types = []string{"010000"}
	// types = []string{"070000"}
	// var pois []amap.Poi

	var mrun = amap.NewMRun(5)

	for _, typess := range types {
		go func(typess string) {
			defer mrun.Done()
			mrun.Add(1)
			pois := amap.RecTypePois(recChangSha, typess)
			var dbPoi amap.DbPoi
			for _, poi := range pois {
				dbPoi.Base = poi
				dbPoi.Source = "gaode"
				dbPoi.PoiKey = poi.ID
				err := h.InsertPoi(dbPoi)
				if err != nil {
					log.Fatal(err.Error())
				}
			}
		}(typess)
		// pois = append(pois, amap.RecTypePois(recChangSha, typess)...)
	}
	mrun.Wait()
	// fmt.Println(pois)

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
