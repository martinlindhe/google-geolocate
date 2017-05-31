package geolocate

import "math"

// Point represents a Physical Point in geographic notation [lat, lng]
// The different possible LocationTypes are documented here:
// https://developers.google.com/maps/documentation/geocoding/intro#Results
// It can for example be "ROOFTOP".
type Point struct {
	Lat          float64
	Lng          float64
	Address      string
	LocationType string
}

const (
	// EarthRadius is about 6,371km according to Wikipedia
	EarthRadius = 6371
)

// NewPoint returns a new Point
func NewPoint(lat, lng float64) *Point {
	return &Point{Lat: lat, Lng: lng}
}

// PointAtDistanceAndBearing returns a Point populated with the lat and lng coordinates
// by transposing the origin point the passed in distance (in kilometers)
// by the passed in compass bearing (in degrees).
// Original Implementation from: http://www.movable-type.co.uk/scripts/latlong.html
func (p *Point) PointAtDistanceAndBearing(dist float64, bearing float64) *Point {

	dr := dist / EarthRadius
	bearing = (bearing * (math.Pi / 180.0))

	lat1 := (p.Lat * (math.Pi / 180.0))
	lng1 := (p.Lng * (math.Pi / 180.0))

	lat2part1 := math.Sin(lat1) * math.Cos(dr)
	lat2part2 := math.Cos(lat1) * math.Sin(dr) * math.Cos(bearing)
	lat2 := math.Asin(lat2part1 + lat2part2)

	lng2part1 := math.Sin(bearing) * math.Sin(dr) * math.Cos(lat1)
	lng2part2 := math.Cos(dr) - (math.Sin(lat1) * math.Sin(lat2))
	lng2 := lng1 + math.Atan2(lng2part1, lng2part2)

	lng2 = math.Mod((lng2+3*math.Pi), (2*math.Pi)) - math.Pi

	lat2 = lat2 * (180.0 / math.Pi)
	lng2 = lng2 * (180.0 / math.Pi)

	return &Point{Lat: lat2, Lng: lng2}
}

// GreatCircleDistance calculates the Haversine distance between two points in kilometers.
// Original Implementation from: http://www.movable-type.co.uk/scripts/latlong.html
func (p *Point) GreatCircleDistance(p2 *Point) float64 {
	dLat := (p2.Lat - p.Lat) * (math.Pi / 180.0)
	dLon := (p2.Lng - p.Lng) * (math.Pi / 180.0)

	lat1 := p.Lat * (math.Pi / 180.0)
	lat2 := p2.Lat * (math.Pi / 180.0)

	a1 := math.Sin(dLat/2) * math.Sin(dLat/2)
	a2 := math.Sin(dLon/2) * math.Sin(dLon/2) * math.Cos(lat1) * math.Cos(lat2)

	a := a1 + a2
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return EarthRadius * c
}

// BearingTo calculates the initial bearing (sometimes referred to as forward azimuth)
// Original Implementation from: http://www.movable-type.co.uk/scripts/latlong.html
func (p *Point) BearingTo(p2 *Point) float64 {

	dLon := (p2.Lng - p.Lng) * math.Pi / 180.0

	lat1 := p.Lat * math.Pi / 180.0
	lat2 := p2.Lat * math.Pi / 180.0

	y := math.Sin(dLon) * math.Cos(lat2)
	x := math.Cos(lat1)*math.Sin(lat2) -
		math.Sin(lat1)*math.Cos(lat2)*math.Cos(dLon)
	brng := math.Atan2(y, x) * 180.0 / math.Pi

	return brng
}

// MidpointTo calculates the midpoint between point `p` and the supplied point.
// Original implementation from http://www.movable-type.co.uk/scripts/latlong.html
func (p *Point) MidpointTo(p2 *Point) *Point {
	lat1 := p.Lat * math.Pi / 180.0
	lat2 := p2.Lat * math.Pi / 180.0

	lon1 := p.Lng * math.Pi / 180.0
	dLon := (p2.Lng - p.Lng) * math.Pi / 180.0

	bx := math.Cos(lat2) * math.Cos(dLon)
	by := math.Cos(lat2) * math.Sin(dLon)

	lat3Rad := math.Atan2(
		math.Sin(lat1)+math.Sin(lat2),
		math.Sqrt(math.Pow(math.Cos(lat1)+bx, 2)+math.Pow(by, 2)),
	)
	lon3Rad := lon1 + math.Atan2(by, math.Cos(lat1)+bx)

	return &Point{
		Lat: lat3Rad * 180.0 / math.Pi,
		Lng: lon3Rad * 180.0 / math.Pi,
	}
}
