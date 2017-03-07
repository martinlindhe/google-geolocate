package geolocate

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExampleGeocode() {
	client := NewGoogleGeo("")
	res, _ := client.Geocode("New York City")
	fmt.Println(res)
	// Output: &{40.7127837 -74.0059413 New York, NY, USA}
}

func ExampleReverseGeocode() {
	client := NewGoogleGeo("")
	p := Point{Lat: 40.7127837, Lng: -74.0059413}
	res, _ := client.ReverseGeocode(&p)
	fmt.Println(res)
	// Output: New York City Hall, New York, NY 10007, USA
}

func ExampleReverseGeocodeDetailed() {
	client := NewGoogleGeo("")
	p := Point{Lat: 40.7127837, Lng: -74.0059413}
	res, _ := client.ReverseGeocodeDetailed(&p)

	address := DetailsToAddress(&res.Results[0])
	fmt.Printf("%#v", address)
	// Output: &geolocate.Address{StreetNumber:"", Locality:"New York", Sublocality:"Manhattan", Neighborhood:"Lower Manhattan", Route:"", PostalCode:"10007", Country:"United States", AdministrativeAreaLevel1:"New York", AdministrativeAreaLevel2:"New York County"}
}

func TestGoogleGeocoder(t *testing.T) {
	// Empty API Key
	c := NewGoogleGeo("")

	res, err := c.Geocode("Mora, Sweden")
	assert.Equal(t, nil, err)
	assert.Equal(t, "Mora, Sweden", res.Address)
}

func TestGoogleGeocoderQueryStr(t *testing.T) {

	// Empty API Key
	c := NewGoogleGeo("")

	address := "123 fake st"
	res, err := c.geocodeQueryStr(address)
	if err != nil {
		t.Errorf("Error creating query string: %v", err)
	}

	expected := "address=123+fake+st&key="
	if res != expected {
		t.Errorf(fmt.Sprintf("Mismatched query string.  Expected: %s.  Actual: %s", expected, res))
	}

	// Set api key to some value
	c.SetGoogleAPIKey("foo")
	res, err = c.geocodeQueryStr(address)
	if err != nil {
		t.Errorf("Error creating query string: %v", err)
	}

	expected = "address=123+fake+st&key=foo"
	if res != expected {
		t.Errorf(fmt.Sprintf("Mismatched query string.  Expected: %s.  Actual: %s", expected, res))
	}
}

func TestGoogleReverseGeocoderQueryStr(t *testing.T) {
	// Empty API Key
	c := NewGoogleGeo("")

	p := &Point{Lat: 123.45, Lng: 56.78}
	res, err := c.reverseGeocodeQueryStr(p)
	if err != nil {
		t.Errorf("Error creating query string: %v", err)
	}

	expected := "latlng=123.450000,56.780000&key="
	if res != expected {
		t.Errorf(fmt.Sprintf("Mismatched query string.  Expected: %s.  Actual: %s", expected, res))
	}

	// Set api key to some value
	c.SetGoogleAPIKey("foo")
	res, err = c.reverseGeocodeQueryStr(p)
	if err != nil {
		t.Errorf("Error creating query string: %v", err)
	}

	expected = "latlng=123.450000,56.780000&key=foo"
	if res != expected {
		t.Errorf(fmt.Sprintf("Mismatched query string.  Expected: %s.  Actual: %s", expected, res))
	}
}
