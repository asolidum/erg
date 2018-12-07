package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	. "bitbucket.org/happyhourpal/erg/src/models"
	. "bitbucket.org/happyhourpal/erg/src/revgeos"
)

var (
	wgRGs        sync.WaitGroup
	latLon       string
	llColsString string
	fields       string
	filename     string
	sleepSeconds int
	append       bool
)

func outputLocation(rgStats *RGStats, line string) {
	if filename == "" {
		fmt.Printf("GPS Coord - %f, %f\n", rgStats.Latitude, rgStats.Longitude)
		fmt.Printf("  OpenCage - %s\n", rgStats.OpenCage)
		fmt.Printf("  FourSquare - %s\n", rgStats.Fourquare)
		fmt.Printf("  Yelp - %s\n", rgStats.Yelp)
		fmt.Printf("  MapBox - %s\n", rgStats.MapBox)
		fmt.Printf("  OSM - %s\n", rgStats.OSM)
	} else {
		output := fmt.Sprintf("%s,%s,%s,%s,%s", rgStats.OpenCage, rgStats.Fourquare, rgStats.Yelp, rgStats.MapBox, rgStats.OSM)
		if append {
			fmt.Printf("%s,%s\n", line, output)
		} else {
			fmt.Printf("%f,%f,%s\n", rgStats.Latitude, rgStats.Longitude, output)
		}
	}
}

func retrieveLocation(latitude float64, longitude float64, line string) {
	defer wgRGs.Done()

	var rgStats *RGStats
	rgStats = new(RGStats)
	rgStats.Latitude = latitude
	rgStats.Longitude = longitude
	rgStats.RadiusMeters = 10

	ch := make(chan string)
	go GetOpenCage(ch, rgStats)
	go GetFourSquare(ch, rgStats)
	go GetYelp(ch, rgStats)
	go GetMapbox(ch, rgStats)
	go GetOSM(ch, rgStats)

	sitesVisited := 5

siteVisits:
	for sitesVisited > 0 {
		select {
		// case site := <-ch:
		// 	fmt.Printf("Visited %s\n", site)
		case <-ch:
			sitesVisited--
		case <-time.After(5 * time.Second):
			fmt.Printf("\033[1;33mTimed out for %d\n\033[0;37m", rgStats.Latitude)
			break siteVisits
		}
	}
	outputLocation(rgStats, line)
}

func getLatLonFromString(latlon string) (float64, float64, error) {
	location := strings.Split(latlon, ",")
	var latitude float64
	var longitude float64
	var err error

	latitude, err = strconv.ParseFloat(location[0], 64)
	if err == nil {
		longitude, err = strconv.ParseFloat(location[1], 64)
		if err == nil {
			return latitude, longitude, err
		}
	}

	return 0, 0, err
}

func init() {
	flag.StringVar(&latLon, "latlon", "47.6124759,-122.3395052", "Lat,Lon to compare")
	flag.StringVar(&filename, "filename", "", "Filename to get lat/lon")
	flag.StringVar(&llColsString, "cols", "1,2", "Cols that contain the lat,lon (respectively also only when using input file)")
	flag.IntVar(&sleepSeconds, "sleep", 1, "Time to wait before attempting another rev geo request (sec)")
	flag.BoolVar(&append, "append", false, "Append data to end of row in file")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
}

func main() {
	if filename == "" {
		latitude, longitude, err := getLatLonFromString(latLon)
		if err != nil {
			log.Fatal(err)
		}
		wgRGs.Add(1)
		go retrieveLocation(latitude, longitude, "")
		wgRGs.Wait()
	} else {
		file, err := os.Open(filename)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		llCols := strings.Split(llColsString, ",")
		latCol, err := strconv.ParseInt(llCols[0], 10, 32)
		if err != nil {
			log.Fatal("Invalid latitude column value - %s\n", llCols[0])
		}
		lonCol, err := strconv.ParseInt(llCols[1], 10, 32)
		if err != nil {
			log.Fatal("Invalid latitude column value - %s\n", llCols[1])
		}

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			cols := strings.Split(line, ",")
			if len(cols) < 2 {
				fmt.Printf("Line does not have enough columns\n")
				continue
			}
			latitude, err := strconv.ParseFloat(cols[latCol-1], 64)
			if err != nil {
				fmt.Printf("Can't parse latitude - %s\n", cols[latCol-1])
				continue
			}
			longitude, err := strconv.ParseFloat(cols[lonCol-1], 64)
			if err != nil {
				fmt.Printf("Can't parse longitude - %s\n", cols[lonCol-1])
				continue
			}
			wgRGs.Add(1)
			go retrieveLocation(latitude, longitude, line)
			time.Sleep(time.Duration(sleepSeconds) * time.Second)
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
		wgRGs.Wait()
	}
}
