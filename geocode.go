package geolocate

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// GoogleGeo contains all the functionality of interacting with the Google Maps Geocoding Service
type GoogleGeo struct {
	client *http.Client
	apiKey string
	region string
}

// NewGoogleGeo returns a new GoogleGeo instance
func NewGoogleGeo(apiKey string) *GoogleGeo {
	return &GoogleGeo{
		client: &http.Client{Timeout: time.Second * 10},
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
			LocationType string `json:"location_type"`
		}
	}
}

// GoogleReverseGeocodeResponse is the return format of a reverse geocode response
type GoogleReverseGeocodeResponse struct {
	Results []reverseGeocodeResult
}

type reverseGeocodeResult struct {
	AddressComponents []reverseGeocodeAddressComponent `json:"address_components"`
	FormattedAddress  string                           `json:"formatted_address"`
	Types             []string                         `json:"types"`
}

type reverseGeocodeAddressComponent struct {
	LongName  string   `json:"long_name"`
	ShortName string   `json:"short_name"`
	Types     []string `json:"types"`
}

// This is the error that consumers receive when there
// are no results from the geocoding request.
var errGoogleZeroResults = errors.New("ZERO_RESULTS")

// SetGoogleAPIKey sets the API key to use
func (g *GoogleGeo) SetGoogleAPIKey(newAPIKey string) {
	g.apiKey = newAPIKey
}

// GeocodeWithRegion geocodes the passed in query string and returns a pointer to a new Point struct.
// Returns an error if the underlying request cannot complete.
// Read about region usage here: https://developers.google.com/maps/documentation/geocoding/intro#RegionCodes
func (g *GoogleGeo) GeocodeWithRegion(address, region string) (*Point, error) {
	return g.geocode(address, region)
}

// Geocode geocodes the passed in query string and returns a pointer to a new Point struct.
// Returns an error if the underlying request cannot complete.
func (g *GoogleGeo) Geocode(address string) (*Point, error) {
	return g.geocode(address, "")
}

func (g *GoogleGeo) geocode(address, region string) (*Point, error) {
	queryStr, err := g.geocodeQueryStr(address, region)
	if err != nil {
		return nil, err
	}
	data, err := g.request(queryStr)
	if err != nil {
		return nil, err
	}
	res := &googleGeocodeResponse{}
	if err := json.Unmarshal(data, res); err != nil {
		return nil, err
	}
	if len(res.Results) == 0 {
		return nil, errGoogleZeroResults
	}

	lat := res.Results[0].Geometry.Location.Lat
	lng := res.Results[0].Geometry.Location.Lng
	loc := res.Results[0].Geometry.LocationType
	point := &Point{
		Address:      res.Results[0].FormattedAddress,
		Lat:          lat,
		Lng:          lng,
		LocationType: loc,
	}
	return point, nil
}

// ReverseGeocode geocodes the pointer to a Point struct and returns the first address that matches
// or returns an error if the underlying request cannot complete.
func (g *GoogleGeo) ReverseGeocode(p *Point) (string, error) {
	res, err := g.ReverseGeocodeDetailed(p)
	if err != nil {
		return "", err
	}
	return res.Results[0].FormattedAddress, err
}

// ReverseGeocodeDetailed geocodes the pointer to a Point struct and returns detailed address information.
func (g *GoogleGeo) ReverseGeocodeDetailed(p *Point) (*GoogleReverseGeocodeResponse, error) {
	queryStr, err := g.reverseGeocodeQueryStr(p)
	if err != nil {
		return nil, err
	}
	data, err := g.request(queryStr)
	if err != nil {
		return nil, err
	}
	res := &GoogleReverseGeocodeResponse{}
	if err := json.Unmarshal(data, res); err != nil {
		return nil, err
	}
	if len(res.Results) == 0 {
		return nil, errGoogleZeroResults
	}
	return res, nil
}

// request issues a request to the google geocoding service and forwards the passed in params string
// as a URL-encoded entity.  Returns an array of byes as a result, or an error if one occurs during the process.
// Note: Since this is an arbitrary request, you are responsible for passing in your API key if you want one.
func (g *GoogleGeo) request(params string) ([]byte, error) {
	if g.client == nil {
		g.client = &http.Client{Timeout: time.Second * 10}
	}

	baseURL := "https://maps.googleapis.com/maps/api/geocode/json?sensor=false&"
	if g.region != "" {
		baseURL += "region=" + g.region + "&"
	}

	fullURL := baseURL + params
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, err
	}
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

func (g *GoogleGeo) geocodeQueryStr(address, region string) (string, error) {
	safeQuery := url.QueryEscape(address)
	var queryStr = bytes.NewBufferString("")
	if _, err := queryStr.WriteString(fmt.Sprintf("address=%s", safeQuery)); err != nil {
		return "", err
	}
	if region != "" {
		if _, err := queryStr.WriteString(fmt.Sprintf("&region=%s", region)); err != nil {
			return "", err
		}
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

// DetailsToAddress parses a result into an Address object
func DetailsToAddress(det *reverseGeocodeResult) *Address {
	adr := Address{}

	for _, d := range det.AddressComponents {
		for _, t := range d.Types {
			switch t {
			case "street_number":
				adr.StreetNumber = d.LongName
			case "locality":
				adr.Locality = d.LongName
			case "sublocality":
				adr.Sublocality = d.LongName
			case "country":
				adr.Country = d.LongName
			case "postal_code":
				adr.PostalCode = d.LongName
			case "neighborhood":
				adr.Neighborhood = d.LongName
			case "administrative_area_level_1":
				adr.AdministrativeAreaLevel1 = d.LongName
			case "administrative_area_level_2":
				adr.AdministrativeAreaLevel2 = d.LongName
			default:
				//log.Println("unknown type", t)
			}
		}
	}

	return &adr
}

// Address holds a parsed address result
type Address struct {
	StreetNumber             string
	Locality                 string
	Sublocality              string
	Neighborhood             string
	Route                    string
	PostalCode               string
	Country                  string
	AdministrativeAreaLevel1 string
	AdministrativeAreaLevel2 string
}
