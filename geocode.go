package geolocate

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// GoogleGeo contains all the functionality of interacting with the Google Maps Geocoding Service
type GoogleGeo struct {
	client *http.Client
	apiKey string
}

// NewGoogleGeo returns a new GoogleGeo instance
func NewGoogleGeo(apiKey string) *GoogleGeo {
	return &GoogleGeo{
		client: &http.Client{},
		apiKey: apiKey,
	}
}

// This struct contains selected fields from Google's Geocoding Service response
type googleGeocodeResponse struct {
	Results []struct {
		FormattedAddress string `json:"formatted_address"`
		Geometry         struct {
			Location struct {
				Lat float64
				Lng float64
			}
		}
	}
}

type googleReverseGeocodeResponse struct {
	Results []struct {
		FormattedAddress string `json:"formatted_address"`
	}
}

// This is the error that consumers receive when there
// are no results from the geocoding request.
var errGoogleZeroResults = errors.New("ZERO_RESULTS")

// SetGoogleAPIKey sets the API key to use
func (g *GoogleGeo) SetGoogleAPIKey(newAPIKey string) {
	g.apiKey = newAPIKey
}

// Geocode geocodes the passed in query string and returns a pointer to a new Point struct.
// Returns an error if the underlying request cannot complete.
func (g *GoogleGeo) Geocode(address string) (*Point, error) {
	queryStr, err := g.geocodeQueryStr(address)
	if err != nil {
		return nil, err
	}
	data, err := g.request(queryStr)
	if err != nil {
		return nil, err
	}
	res := &googleGeocodeResponse{}
	json.Unmarshal(data, res)
	if len(res.Results) == 0 {
		return nil, errGoogleZeroResults
	}

	lat := res.Results[0].Geometry.Location.Lat
	lng := res.Results[0].Geometry.Location.Lng
	point := &Point{
		Address: res.Results[0].FormattedAddress,
		Lat:     lat,
		Lng:     lng,
	}
	return point, nil
}

// ReverseGeocode geocodes the pointer to a Point struct and returns the first address that matches
// or returns an error if the underlying request cannot complete.
func (g *GoogleGeo) ReverseGeocode(p *Point) (string, error) {
	queryStr, err := g.reverseGeocodeQueryStr(p)
	if err != nil {
		return "", err
	}
	data, err := g.request(queryStr)
	if err != nil {
		return "", err
	}
	res := &googleReverseGeocodeResponse{}
	err = json.Unmarshal(data, res)
	if err != nil {
		return "", err
	}
	if len(res.Results) == 0 {
		return "", errGoogleZeroResults
	}
	return res.Results[0].FormattedAddress, err
}

// request issues a request to the google geocoding service and forwards the passed in params string
// as a URL-encoded entity.  Returns an array of byes as a result, or an error if one occurs during the process.
// Note: Since this is an arbitrary request, you are responsible for passing in your API key if you want one.
func (g *GoogleGeo) request(params string) ([]byte, error) {
	if g.client == nil {
		g.client = &http.Client{}
	}

	fullURL := "https://maps.googleapis.com/maps/api/geocode/json?sensor=false&" + params
	req, _ := http.NewRequest("GET", fullURL, nil)
	resp, requestErr := g.client.Do(req)
	if requestErr != nil {
		return nil, requestErr
	}

	data, dataReadErr := ioutil.ReadAll(resp.Body)
	if dataReadErr != nil {
		return nil, dataReadErr
	}
	return data, nil
}

func (g *GoogleGeo) geocodeQueryStr(address string) (string, error) {
	safeQuery := url.QueryEscape(address)
	var queryStr = bytes.NewBufferString("")
	if _, err := queryStr.WriteString(fmt.Sprintf("address=%s", safeQuery)); err != nil {
		return "", err
	}
	if _, err := queryStr.WriteString(fmt.Sprintf("&key=%s", g.apiKey)); err != nil {
		return "", err
	}
	return queryStr.String(), nil
}

func (g *GoogleGeo) reverseGeocodeQueryStr(p *Point) (string, error) {
	var queryStr = bytes.NewBufferString("")
	if _, err := queryStr.WriteString(fmt.Sprintf("latlng=%f,%f", p.Lat, p.Lng)); err != nil {
		return "", err
	}
	if _, err := queryStr.WriteString(fmt.Sprintf("&key=%s", g.apiKey)); err != nil {
		return "", err
	}
	return queryStr.String(), nil
}
