package geolocate

import (
	"fmt"
	"testing"
)

func TestGreatCircleDistance(t *testing.T) {
	// Test that SEA and SFO are ~ 1091km apart, accurate to 100 meters.
	sea := &Point{Lat: 47.4489, Lng: -122.3094}
	sfo := &Point{Lat: 37.6160933, Lng: -122.3924223}
	sfoToSea := 1093.379199082169

	dist := sea.GreatCircleDistance(sfo)
	if !(dist < (sfoToSea+0.1) && dist > (sfoToSea-0.1)) {
		t.Error("Unnacceptable result.", dist)
	}
}

func TestPointAtDistanceAndBearing(t *testing.T) {
	sea := &Point{Lat: 47.44745785, Lng: -122.308065668024}
	p := sea.PointAtDistanceAndBearing(1090.7, 180)

	// Expected results of transposing point
	// ~1091km at bearing of 180 degrees
	resultLat := 37.638557
	resultLng := -122.308066

	withinLatBounds := p.Lat < resultLat+0.001 && p.Lat > resultLat-0.001
	withinLngBounds := p.Lng < resultLng+0.001 && p.Lng > resultLng-0.001
	if !(withinLatBounds && withinLngBounds) {
		t.Error("Unnacceptable result.", fmt.Sprintf("[%f, %f]", p.Lat, p.Lng))
	}
}

func TestBearingTo(t *testing.T) {
	p1 := &Point{Lat: 40.7486, Lng: -73.9864}
	p2 := &Point{Lat: 0.0, Lng: 0.0}
	bearing := p1.BearingTo(p2)

	// Expected bearing 60 degrees
	resultBearing := 100.610833

	withinBearingBounds := bearing < resultBearing+0.001 && bearing > resultBearing-0.001
	if !withinBearingBounds {
		t.Error("Unnacceptable result.", fmt.Sprintf("%f", bearing))
	}
}

func TestMidpointTo(t *testing.T) {
	p1 := &Point{Lat: 52.205, Lng: 0.119}
	p2 := &Point{Lat: 48.857, Lng: 2.351}

	p := p1.MidpointTo(p2)

	// Expected midpoint 50.5363°N, 001.2746°E
	resultLat := 50.53632
	resultLng := 1.274614

	withinLatBounds := p.Lat < resultLat+0.001 && p.Lat > resultLat-0.001
	withinLngBounds := p.Lng < resultLng+0.001 && p.Lng > resultLng-0.001
	if !(withinLatBounds && withinLngBounds) {
		t.Error("Unnacceptable result.", fmt.Sprintf("[%f, %f]", p.Lat, p.Lng))
	}
}
