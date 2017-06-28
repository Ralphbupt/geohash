package geohash

import (
	"fmt"
	"math"
)

type HashBits struct {
	Bits uint64
	Step uint
}
type HashRadious struct {
	hash     HashBits
	area     Area
	neibours Neighbors
}

type Neighbors struct {
	North     HashBits
	NorthEast HashBits
	East      HashBits
	SouthEast HashBits
	South     HashBits
	SouthWest HashBits
	West      HashBits
	NorthWest HashBits
}

const (
	MERCATOR_MAX = 20037726.37
	MERCATOR_MIN = -20037726.37
)

func GetDefaultArea() Area {
	area := Area{}
	area.Latitude.Min, area.Latitude.Max = MIN_LATITUDE, MAX_LATITUDE
	area.Longitude.Min, area.Longitude.Max = MIN_LONGITUDE, MAX_LONGITUDE
	return area
}

/*func GetNeighborsByRadio(latitude float64, longitude float64, radious_meters float64) HashRadious {
	//
	var long_range, lat_range Range
	var radious HashRadious
	var neighbors Neighbors
	bounds := BoundInBox(latitude, longitude, radious_meters)


	area := GetDefaultArea()
	lat_range.Min,	lat_range.Max = -88.0,88.0
	long_range.Min, long_range.Max = -180.0, 180.0

	step := estimateStepsByRadious(radious_meters, latitude)
	hash, _ := EncodeToBits(latitude, longitude, step)
	neighbors = GetNeighbors(*hash)
	//Decode()
	return HashRadious{}
}
*/
func GetNeighbors(bits HashBits) Neighbors {

	neighbors := Neighbors{}
	neighbors.East = bits
	neighbors.West = bits
	neighbors.South = bits
	neighbors.North = bits
	neighbors.NorthEast = bits
	neighbors.NorthWest = bits
	neighbors.SouthEast = bits
	neighbors.SouthWest = bits

	moveX(&neighbors.East, 1)

	moveX(&neighbors.West, -1)

	moveY(&neighbors.South, -1)

	moveY(&neighbors.North, 1)

	moveX(&neighbors.NorthWest, -1)
	moveY(&neighbors.NorthWest, 1)

	moveX(&neighbors.NorthEast, 1)
	moveY(&neighbors.NorthEast, 1)

	moveX(&neighbors.SouthEast, 1)
	moveY(&neighbors.SouthEast, -1)

	moveX(&neighbors.SouthWest, -1)
	moveY(&neighbors.SouthWest, -1)

	return neighbors
}

func moveY(bits *HashBits, move int) {
	var mask uint64 = 0x0000000000000003
	val := bits.Bits & mask
	switch val {
	case 0, 1:
		if move > 0 {
			bits.Bits += 2
		} else {
			bits.Bits -= 6
		}

	case 2, 3:
		if move > 0 {
			bits.Bits += 6
		} else {
			bits.Bits -= 2
		}

	default:
		fmt.Printf("error occured when sdsfd")
	}

}

func moveX(bits *HashBits, move int) {
	var mask uint64 = 0x0000000000000003
	val := bits.Bits & mask

	switch val {
	case 1, 3:
		if move > 0 {
			bits.Bits += 3
		} else {
			bits.Bits -= 1
		}
	case 0, 2:
		if move > 0 {
			bits.Bits += 1
		} else {
			bits.Bits += 3
		}
	default:
		fmt.Printf("error occured when sdsfd")
	}
	fmt.Println(bits.Bits, move)

}

// estimate the step of the area box, the  hashbit length of longitude and latude
func estimateStepsByRadious(meters float64, latitude float64) uint {
	if meters-0.0 <= 0.0000001 {
		return 26
	}
	step := 1
	for meters < MERCATOR_MAX {
		meters *= 2
		step++
	}

	step -= 2

	// while range torwards to poles, range a latger area
	if latitude > 66 || latitude < -66 {
		step--
		if latitude > 80 || latitude < -80 {
			step--
		}
	}

	if step < 1 {
		step = 1
	}
	if step > 26 {
		step = 26
	}
	return uint(step)
}

func BoundInBox(latitude float64, longitude float64, radius_meters float64) []float64 {
	bounds := make([]float64, 4)
	long_diff := rad2deg(radius_meters / EARTH_RADIU / math.Cos(deg2rad(latitude)))
	lat_diff := rad2deg(radius_meters / EARTH_RADIU)

	bounds[0] = longitude - long_diff
	bounds[2] = longitude + long_diff
	bounds[1] = latitude - lat_diff
	bounds[3] = latitude + lat_diff
	return bounds
}
func Distance(lat1d, long1d, lat2d, long2d float64) float64 {
	lat1r, long1r, lat2r, long2r := deg2rad(lat1d), deg2rad(long1d), deg2rad(lat2d), deg2rad(long2d)
	u := math.Sin((lat2r - lat1r) / 2)
	v := math.Sin((long1r - long2r) / 2)
	// accroding to Haversine formula: https://en.wikipedia.org/wiki/Haversine_formula

	return 2.0 * EARTH_RADIU * math.Asin(math.Sqrt(u*u+math.Cos(lat1r)*math.Cos(lat2r)*v*v))
}

func deg2rad(degree float64) (rad float64) {
	return degree * math.Pi / 180
}

func rad2deg(rad float64) (degree float64) {
	return rad * 180 / math.Pi
}
