package main

//import "bufio"
//import "encoding/csv"
import "net/http"
import "net/url"
import "compress/bzip2"
import "fmt"
import "gopkg.in/alecthomas/kingpin.v2"

//import "io"
import "io/ioutil"
import "os"
import "time"
import "math"

import "github.com/amsokol/go-grib2"
import "github.com/PuerkitoBio/fetchbot"

var (
	files = kingpin.Arg("files", "Filename to load").Strings()
	urls  = kingpin.Flag("url", "URLs to crawl").Strings()
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type latlon struct {
	lat float64
	lon float64
}

func latlon_dist(a latlon, b latlon) float64 {
	// Calculate with opposite lon and same lon, in case lon is near 180
	x := math.Pow(a.lon-b.lon, 2)
	y := math.Pow(a.lat-b.lat, 2)
	dist := math.Sqrt(x + y)

	flip_x := math.Pow(a.lon+b.lon, 2)
	flip_dist := math.Sqrt(flip_x + y)
	return math.Min(dist, flip_dist)
}

type entry struct {
	coord latlon
	time  time.Time
	value float32
}

func parse_file(filename string, targ latlon, entries map[string][]entry) {
	f, err := os.Open(filename)
	check(err)
	defer f.Close()

	uncompressed_data := bzip2.NewReader(f)
	check(err)

	data, err := ioutil.ReadAll(uncompressed_data)
	check(err)

	gribs, err := grib2.Read(data)
	check(err)

	for _, g := range gribs {
		var closest_entry *entry
		var dist float64
		for _, v := range g.Values {
			lon := v.Longitude
			if lon > 180.0 {
				lon -= 360.0
			}

			new_entry := entry{
				coord: latlon{lat: v.Latitude, lon: lon},
				time:  g.VerfTime,
				value: v.Value,
			}
			new_dist := latlon_dist(targ, new_entry.coord)
			if closest_entry == nil {
				closest_entry = &new_entry
				dist = new_dist
			} else {
				if new_dist < dist {
					closest_entry = &new_entry
					dist = latlon_dist(targ, new_entry.coord)
				}
			}
		}
		if closest_entry != nil {
			entries[g.Name] = append(entries[g.Name], *closest_entry)
		}
	}
}

func query_urls(urls []string, entries map[string][]entry) {
	mux := fetchbot.NewMux()

	mux.HandleErrors(fetchbot.HandlerFunc(func(ctx *fetchbot.Context, res *http.Response, err error) {
		fmt.Printf("[ERR] %s %s - %s\n", ctx.Cmd.Method(), ctx.Cmd.URL(), err)
	}))

	hosts := make([]string, 0)
	for _, query_url := range urls {
		u, err := url.Parse(query_url)
		check(err)
		hosts = append(hosts, u.Host)
	}

	mux.Response().Method("GET").ContentType("text/html").Handler(fetchbot.HandlerFunc(
		func(ctx *fetchbot.Context, res *http.Response, err error) {
			fmt.Printf("Hello, %q", ctx.Cmd.URL())
		}))

	for _, host := range hosts {
		mux.Response().Method("HEAD").Host(host).ContentType("text/html").Handler(fetchbot.HandlerFunc(
			func(ctx *fetchbot.Context, res *http.Response, err error) {
				if _, err := ctx.Q.SendStringGet(ctx.Cmd.URL().String()); err != nil {
					fmt.Printf("[ERR] %s %s - %s\n", ctx.Cmd.Method(), ctx.Cmd.URL(), err)
				}
			}))
		}

	f := fetchbot.New(mux)
	queue := f.Start()
	for _, query_url := range urls {
		queue.SendStringHead(query_url)
	}
	queue.
}

func main() {
	kingpin.Parse()
	arns_latlon := latlon{lat: 52.534080, lon: 13.438340}
	entries := make(map[string][]entry)
	if *files != nil {
		for _, filename := range *files {
			parse_file(filename, arns_latlon, entries)
		}
	}
	if *urls != nil {
		query_urls(*urls, entries)
	}
	for _, entry := range entries["TMP"] {
		fmt.Println(
			"hey %v, %v, %v",
			entry.coord,
			entry.time.Format("2006-01-02T15:04:05-0700"),
			entry.value,
		)
	}
}
