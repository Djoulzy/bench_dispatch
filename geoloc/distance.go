package geoloc

import (
	"math"
)

var coeffPyth float64 = 111120  // 1852 * 60
var coeffHav float64 = 12756280 // 6378100 * 2

var maisonLat = 43.32942
var maisonLng = 5.49234

func degreesToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

func hsin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}

// DistanceSimple : Calcul de la distance en metre entre 2 points GPS en utilisant Pythagore
func DistanceSimple(lat1, long1, lat2, long2 float64) float64 {
	deltaY := lat2 - lat1
	deltaX := (long1 - long2) * math.Cos((lat1+lat2)/2)
	dist := math.Sqrt(deltaX*deltaX + deltaY*deltaY)

	return coeffPyth * dist
}

// DistanceAccurate : Haversin
func DistanceAccurate(lat1, lon1, lat2, lon2 float64) float64 {
	la1 := degreesToRadians(lat1)
	la2 := degreesToRadians(lat2)
	long := degreesToRadians(lon2 - lon1)

	h := hsin(la2-la1) + math.Cos(la1)*math.Cos(la2)*hsin(long)

	return coeffHav * math.Asin(math.Sqrt(h))
}

// DistanceFromHome :
func DistanceFromHome(lat float64, long float64) float64 {
	return DistanceSimple(lat, long, maisonLat, maisonLng)
}
