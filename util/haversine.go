package util

import "math"

func toRad(deg float64) float64 {
	return deg * math.Pi / 180
}

func DistanceKm(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371.0 // km

	dLat := toRad(lat2 - lat1)
	dLon := toRad(lon2 - lon1)

	rLat1 := toRad(lat1)
	rLat2 := toRad(lat2)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(rLat1)*math.Cos(rLat2)*math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}
