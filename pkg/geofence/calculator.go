package geofence

import "math"

const earthRadius = 6371000

type Point struct {
	Latitude  float64
	Longitude float64
}

func CalculateDistance(p1, p2 Point) float64 {
	lat1Rad := toRadians(p1.Latitude)
	lat2Rad := toRadians(p2.Latitude)
	deltaLat := toRadians(p2.Latitude - p1.Latitude)
	deltaLon := toRadians(p2.Longitude - p1.Longitude)
	
	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
		math.Sin(deltaLon/2)*math.Sin(deltaLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	
	return earthRadius * c
}

func IsWithinRadius(center Point, target Point, radius float64) bool {
	distance := CalculateDistance(center, target)
	return distance <= radius
}

func toRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}