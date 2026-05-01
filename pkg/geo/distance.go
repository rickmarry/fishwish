package geo

import (
	"math"
)

const earthRadiusKM = 6371.0

type Point struct {
	Lat float64
	Lon float64
}

func Haversine(a, b Point) float64 {
	dLat := deg2rad(b.Lat - a.Lat)
	dLon := deg2rad(b.Lon - a.Lon)

	x := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(deg2rad(a.Lat))*math.Cos(deg2rad(b.Lat))*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(x), math.Sqrt(1-x))
	return earthRadiusKM * c
}

func HaversineMiles(a, b Point) float64 {
	return Haversine(a, b) * 0.621371
}

func deg2rad(d float64) float64 {
	return d * math.Pi / 180.0
}

func ClampLat(lat float64) float64 {
	if lat < -90 {
		return -90
	}
	if lat > 90 {
		return 90
	}
	return lat
}

func ClampLon(lon float64) float64 {
	if lon < -180 {
		return -180
	}
	if lon > 180 {
		return 180
	}
	return lon
}
