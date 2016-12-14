# About

[![GoDoc](https://godoc.org/github.com/martinlindhe/google-geolocate?status.svg)](https://godoc.org/github.com/martinlindhe/google-geolocate)
[![Travis-CI](https://api.travis-ci.org/martinlindhe/google-geolocate.svg)](https://travis-ci.org/martinlindhe/google-geolocate)

Golang client for the Google Maps Geocode and Geolocation API:s

https://developers.google.com/maps/documentation/geolocation/intro


## Usage
```go
import geo "github.com/martinlindhe/google-geolocate"

client := geo.NewGoogleGeo("api-key")
```

## Geocode
```go
res, _ := client.Geocode("New York City")
fmt.Println(res)
// Output: &{40.7127837 -74.0059413 New York, NY, USA}
```

## Reverse geocode
```go
p := geo.Point{Lat: 40.7127837, Lng: -74.0059413}
res, _ := client.ReverseGeocode(&p)
fmt.Println(res)
// Output: New York City Hall, New York, NY 10007, USA
```

## Geolocate

```go
res, _ := client.Geolocate()
fmt.Println(res)
// Output: &{ 40.7127837 -74.0059413}
```

### License

Under [MIT](LICENSE)

Parts of the code was based on [golang-geo](https://github.com/kellydunn/golang-geo)
