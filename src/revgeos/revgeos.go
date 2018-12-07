package revgeos

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	. "bitbucket.org/happyhourpal/erg/src/models"
	"github.com/joho/godotenv"
)

var (
	openCageKey      string // https://opencagedata.com/
	fourSquareID     string // https://developer.foursquare.com/
	fourSquareSecret string
	yelpID           string // https://www.yelp.com/developers
	yelpKey          string
	faceBookKey      string // https://developers.facebook.com/products/places
	hereAppID        string // https://www.here.com/en/products-services/map-content/here-map-data
	hereAppCode      string
	mapBoxKey        string // https://www.mapbox.com/search/
	osmKey           string // https://wiki.openstreetmap.org/wiki/Nominatim
	factualKey       string // -
)

func getRGData(url string) ([]byte, error) {
	var body []byte
	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err == nil {
		body, err = ioutil.ReadAll(resp.Body)
	}

	return body, err
}

func getRGDataViaBearer(url string, bearer string) ([]byte, error) {
	var body []byte
	// Create a Bearer string by appending string access token
	var bearerHeader = "Bearer " + bearer

	// Create a new request using http
	req, err := http.NewRequest("GET", url, nil)
	// add authorization header to the req
	req.Header.Add("Authorization", bearerHeader)

	// Send req using http Client
	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err == nil {
		body, err = ioutil.ReadAll(resp.Body)
	}

	return body, err
}

func GetOpenCage(ch chan string, rgStats *RGStats) {
	provider := "OpenCage"

	url := fmt.Sprintf("https://api.opencagedata.com/geocode/v1/json?q=%f+%f&key=%s&min_confidence=10&abbrv", rgStats.Latitude, rgStats.Longitude, openCageKey)
	data, err := getRGData(url)
	if err != nil {
		ch <- fmt.Sprintf("Err wih %s Resp", provider)
	} else {
		var result OpenCage
		err = json.Unmarshal(data, &result)
		if err != nil || result.Status.Code != 200 {
			ch <- fmt.Sprintf("Err with %s Data - %d", provider, result.Status.Code)
		} else {
			if len(result.Results) > 0 {
				rgStats.OpenCage = strings.Split(result.Results[0].Name, ",")[0]
			}
			ch <- provider
		}
	}
}

func GetFourSquare(ch chan string, rgStats *RGStats) {
	provider := "FourSquare"

	url := fmt.Sprintf("https://api.foursquare.com/v2/venues/search?ll=%f,%f&v=20181201&client_id=%s&client_secret=%s", rgStats.Latitude, rgStats.Longitude, fourSquareID, fourSquareSecret)
	data, err := getRGData(url)
	if err != nil {
		ch <- fmt.Sprintf("Err wih %s Resp", provider)
	} else {
		var result FourSquare
		err = json.Unmarshal(data, &result)
		if err != nil || result.Meta.Code != 200 {
			ch <- fmt.Sprintf("Err with %s Data - %d", provider, result.Meta.Code)
		} else {
			if len(result.Response.Venues) > 0 {
				rgStats.Fourquare = result.Response.Venues[0].Name
			}
			ch <- provider
		}
	}
}

func GetYelp(ch chan string, rgStats *RGStats) {
	provider := "Yelp"

	url := fmt.Sprintf("https://api.yelp.com/v3/businesses/search?latitude=%f&longitude=%f&radius=%d", rgStats.Latitude, rgStats.Longitude, rgStats.RadiusMeters)
	data, err := getRGDataViaBearer(url, yelpKey)
	if err != nil {
		ch <- fmt.Sprintf("Err wih %s Resp", provider)
	} else {
		var result Yelp
		err = json.Unmarshal(data, &result)
		if err != nil {
			ch <- fmt.Sprintf("Err with %s Data", provider)
		} else {
			if len(result.Businesses) > 0 {
				rgStats.Yelp = result.Businesses[0].Name
			}
			ch <- provider
		}
	}
}

func GetMapbox(ch chan string, rgStats *RGStats) {
	provider := "MapBox"

	url := fmt.Sprintf("https://api.mapbox.com/geocoding/v5/mapbox.places/%f,%f.json?access_token=%s", rgStats.Longitude, rgStats.Latitude, mapBoxKey)
	data, err := getRGDataViaBearer(url, yelpKey)
	if err != nil {
		ch <- fmt.Sprintf("Err wih %s Resp", provider)
	} else {
		var result MapBox
		err = json.Unmarshal(data, &result)
		if err != nil {
			ch <- fmt.Sprintf("Err with %s Data", provider)
		} else {
			if len(result.Features) > 0 {
				rgStats.MapBox = result.Features[0].Name
			}
			ch <- provider
		}
	}
}

func GetOSM(ch chan string, rgStats *RGStats) {
	provider := "OSM"

	url := fmt.Sprintf("https://nominatim.openstreetmap.org/reverse?format=jsonv2&lat=%f&lon=%f", rgStats.Latitude, rgStats.Longitude)
	data, err := getRGData(url)
	if err != nil {
		ch <- fmt.Sprintf("Err wih %s Resp", provider)
	} else {
		var result OSM
		//		fmt.Printf("OSM = %s\n", data)
		err = json.Unmarshal(data, &result)
		if err != nil {
			ch <- fmt.Sprintf("Err with %s Data", provider)
		} else {
			rgStats.OSM = result.Name
			ch <- provider
		}
	}
}

// Return the absolute path of any path
func getAbsPath(path string) string {
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		log.Printf("Unable to get absolute path %s\n", err.Error())
	}
	return absolutePath
}

func setupVars() {
	openCageKey = os.Getenv("OPENCAGE_KEY")
	fourSquareID = os.Getenv("FOURSQUARE_ID")
	fourSquareSecret = os.Getenv("FOURSQUARE_SECRET")
	yelpID = os.Getenv("YELP_ID")
	yelpKey = os.Getenv("YELP_KEY")
	hereAppID = os.Getenv("HERE_APP_ID")
	hereAppCode = os.Getenv("HERE_APP_CODE")
	mapBoxKey = os.Getenv("MAPBOX_TOKEN")
}

func init() {
	envp := getAbsPath("./.env")
	err := godotenv.Load(envp)
	if err != nil {
		log.Println("Could not load .env file. Checking for necessary environmental variables...")
	}
	setupVars()
}
