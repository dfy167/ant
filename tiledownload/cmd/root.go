package cmd

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"enen.com/superenen/map/tiledownload/td"
	"github.com/spf13/cobra"
)

var outputpath string
var boundary string
var maptype int
var maxzoom int
var minzoom int
var mapapikey string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tiledownload",
	Short: "little tools",
	Long: `
	
A little tools to download map tile

Example: 
	tiledownload.exe -o d:/tilemap/tdtnormal -z 10 -Z 17 -t 1 -b 119.365,44.22,123.700,47.688 -k 2ce94f67e58faa24beb7cb8a09780552

Attention: 
	1. Please create output dir by youself
	2. Do not run multiple processes on the same computer or on diffrent computer with same internet IP  . It will not be faster, but slower.
	3. Recommend to run 16，17，18 level  alone. When breaks, it will not resume but starts from the very first.	
	`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if !exists(outputpath) {
			fmt.Println("out put dir not exists")
			return nil
		}

		mapT, err := ckeckMaptype(maptype)
		if err != nil {
			fmt.Println(err)
			return nil
		}

		if maxzoom > 18 || maxzoom < 1 {
			fmt.Println("max zoom should between 1 and 18")
			return nil
		}

		if minzoom > 18 || minzoom < 1 {
			fmt.Println("min zoom should between 1 and 18")
			return nil
		}

		if maxzoom < minzoom {
			fmt.Println("max zoom should larger than min zoom")
			return nil
		}

		bd, err := ckeckBounds(boundary)
		if err != nil {
			fmt.Println(err)
			return nil
		}

		err = td.CheckKey(mapapikey, mapT)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		start := time.Now()
		// fmt.Println(outputpath, mapT, mapapikey, maxzoom, minzoom, bd)
		err = td.Tdfast(outputpath, mapT, mapapikey, maxzoom, minzoom, bd)
		// err = td.Td(outputpath, mapT, mapapikey, maxzoom, minzoom, bd)

		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("used time: %s\n", time.Since(start))
		return nil
	},
}

// Execute Execute
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&outputpath, "outputpath", "o", "", "output Dir , create by youself. split with / (required)")
	rootCmd.Flags().StringVarP(&boundary, "boundary", "b", "", `boundary to download  (required)
	Format: leftbotton to righttop. example: 108,24,115, 31 
	Earth LngLatBounds:{LngLat{-180, -85.0511}, LngLat{180, 85.0511}} `)

	rootCmd.Flags().IntVarP(&maptype, "maptype", "t", 1, `maptype to download . (required)
	1:Tianditu Normal Map
	2:Tianditu Normal Annotion
	3:Tianditu Terrain Map
	4:TiandituTerrain Annotion
	`)

	rootCmd.Flags().IntVarP(&minzoom, "minzoom", "z", 1, "min zoom level  to download . 1-18 (required)")
	rootCmd.Flags().IntVarP(&maxzoom, "maxzoom", "Z", 1, "max zoom level to download . 1-18 (required)")
	rootCmd.Flags().StringVarP(&mapapikey, "mapapikey", "k", "", "api key of map provider . (required)")

	rootCmd.MarkFlagRequired("outputpath")
	rootCmd.MarkFlagRequired("boundary")
	rootCmd.MarkFlagRequired("maptype")
	rootCmd.MarkFlagRequired("maxzoom")
	rootCmd.MarkFlagRequired("minzoom")
	rootCmd.MarkFlagRequired("mapapikey")

}

func exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// 判断所给路径是否为文件夹
func isDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// 判断所给路径是否为文件
func isFile(path string) bool {
	return !isDir(path)
}

func ckeckBounds(b string) (bounds td.LngLatBounds, err error) {
	bs := strings.Split(b, ",")
	if len(bs) != 4 {
		err = errors.New("boundary is wrong")
		return
	}

	bs1, err := strconv.ParseFloat(bs[0], 64)
	if err != nil {
		err = errors.New("boundary is wrong")
	}
	bs2, err := strconv.ParseFloat(bs[1], 64)
	if err != nil {
		err = errors.New("boundary is wrong")
	}
	bs3, err := strconv.ParseFloat(bs[2], 64)
	if err != nil {
		err = errors.New("boundary is wrong")
	}
	bs4, err := strconv.ParseFloat(bs[3], 64)
	if err != nil {
		err = errors.New("boundary is wrong")
	}

	if bs1 > bs3 {
		err = errors.New("boundary is wrong")
	}
	if bs2 > bs4 {
		err = errors.New("boundary is wrong")
	}

	if bs1 < -180 || bs1 > 180 || bs3 < -180 || bs3 > 180 || bs2 < -85.0511 || bs2 > 85.0511 || bs4 < -85.0511 || bs4 > 85.0511 {
		err = errors.New("boundary is wrong")
	}
	bounds = td.NewBounds(bs1, bs2, bs3, bs4)
	return
}

func ckeckMaptype(b int) (mapt string, err error) {
	switch b {
	case 0:
		return td.Mapbox, err
	case 1:
		return td.TdtNormalMap, err
	case 2:
		return td.TdtNormalAnnotion, err
	case 3:
		return td.TdtTerrainMap, err
	case 4:
		return td.TdtTerrainAnnotion, err
	}
	return "", errors.New("maptype not exists")
}
