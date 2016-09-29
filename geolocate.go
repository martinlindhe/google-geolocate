package geolocate

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/kellydunn/golang-geo"
)

// Geolocate returns a rough location based on your IP
func Geolocate(apiKey string) (*geo.Point, error) {
	data, err := request(apiKey)
	if err != nil {
		return nil, err
	}

	res := &geolocateResponse{}
	json.Unmarshal(data, res)
	if res.Error.Code != 0 {
		e := res.Error.Errors[0]
		return nil, fmt.Errorf(e.Domain + "." + e.Reason + "." + e.Message)
	}

	point := geo.NewPoint(res.Location.Lat, res.Location.Lng)
	return point, nil
}

// This struct contains selected fields from Google's Geocoding Service response
type geolocateResponse struct {
	Error struct {
		Code    int
		Message string
		Errors  []struct {
			Domain  string
			Reason  string
			Message string
		}
	}
	Location struct {
		Lat float64
		Lng float64
	}
}

func request(apiKey string) ([]byte, error) {
	client := &http.Client{}

	dst := "https://www.googleapis.com/geolocation/v1/geolocate?key=" + apiKey

	form := url.Values{}
	// form.Add("considerIp", "true")

	req, _ := http.NewRequest("POST", dst, strings.NewReader(form.Encode()))
	resp, requestErr := client.Do(req)

	if requestErr != nil {
		return nil, requestErr
	}

	data, dataReadErr := ioutil.ReadAll(resp.Body)

	if dataReadErr != nil {
		return nil, dataReadErr
	}

	return data, nil
}
