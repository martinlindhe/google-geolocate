package geolocate

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExampleGeolocate() {
	client := NewGoogleGeo("")
	res, _ := client.Geolocate()
	fmt.Println(res)
	// Output: <nil>
}

func TestGeolocateNoKey(t *testing.T) {
	// Empty API Key
	c := NewGoogleGeo("")

	point, err := c.Geolocate()
	if point != nil {
		t.Error("unexpected result")
	}
	assert.Equal(t, "Google API key not provided", err.Error())
}
