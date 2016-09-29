package geolocate

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	apiKey = "XXX" // provide a valid api key to run tests
)

func TestNoKey(t *testing.T) {
	point, err := Geolocate("")
	if point != nil {
		t.Error("unexpected result")
	}
	assert.Equal(t, "usageLimits.keyInvalid.Bad Request", err.Error())
}

func TestWithKey(t *testing.T) {
	point, err := Geolocate(apiKey)
	assert.Equal(t, nil, err)
	fmt.Println(point)
}
