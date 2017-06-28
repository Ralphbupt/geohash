package geohash

import (
	"errors"

	"strings"
	"sync"
	"fmt"

)

type GeoPoint struct {
	Longitude float64
	Latitude  float64
	Hashstr   string
	Score     float64
}

type LocationMap struct {
	Map map[string]*GeoPoint
	mutex *sync.RWMutex
}


func NewMap() LocationMap{
	m := make(map[string]*GeoPoint)

	return LocationMap{
		Map:m,
		mutex:new(sync.RWMutex),
	}
}

/*
func ( l * LocationMap)Add(latitude float64, longitude float64) {
	l.mutex.Lock()
	l.mutex.Unlock()

	bits := Encode2(latitude, longitude, 26)
	hashstr := Encode(latitude, longitude,11 )

}
*/


type Area struct {
	Latitude  Range
	Longitude Range
}

type Range struct {
	Min float64
	Max float64
}

const (
	GEO_LENGTH            = 11
	EARTH_RADIU           = 6372797.560856
	BASE_32               = "0123456789bcdefghjkmnpqrstuvwxyz"
	MAX_LATITUDE  float64 = 90.0
	MIN_LATITUDE  float64 = -90.0
	MAX_LONGITUDE float64 = 180.0
	MIN_LONGITUDE float64 = -180.0
)

var (
	ErrInvalidLocation  = errors.New("invalid input location!")
	ErrInvalidAccuracy  = errors.New("input  invalid  accuricy")
	ErrInvalidcharacter = errors.New("invalid input charater")
)

var (
	base32map = map[int32]int{
		48: 0, 49: 1, 50: 2, 51: 3, 52: 4, 53: 5, 54: 6, 55: 7,
		56: 8, 57: 9, 98: 10, 99: 11, 100: 12, 101: 13, 102: 14, 103: 15,
		104: 16, 106: 17, 107: 18, 109: 19, 110: 20, 112: 21, 113: 22, 114: 23,
		115: 24, 116: 25, 117: 26, 118: 27, 119: 28, 120: 29, 121: 30, 122: 31,
	}
	defaultArea = Area{Range{-90.0, 90.0}, Range{-180.0, 180.0}}
)

// geohash accuricy
/*
// *******		length		lat bits	lng bits	lat error	lng error	km error
// *******		1			2			3			±23			±23			±2500
// *******		2			5			5			±2.8	 	±5.6	 	±630
// *******		3			7			8	 		±0.70	 	±0.70	  	±78
// *******		4			10			10	 		±0.087	 	±0.18	  	±20
// *******		5			12			13	 		±0.022	 	±0.022	  	±2.4
// *******		6			15			15	 		±0.0027	 	±0.0055	   	±0.61
// *******		7			17			18	 		±0.00068	±0.00068	±0.076
// *******		8			20			20	 		±0.000085	±0.00017	±0.019
// *******		9  			22			23			±0.000021	±0.000021 	±0.0024
// *******		10  		25			25			±0.0000027	±0.0000054	±0.00060
// *******		11  		27			28			±0.00000067	±0.00000067	±0.000074
*/
func Encode(latitude, longitude float64, accuracy int) (string, error) {

	if latitude < MIN_LATITUDE || latitude > MAX_LATITUDE || longitude < MIN_LONGITUDE || longitude > MAX_LONGITUDE {
		return "", ErrInvalidLocation
	}

	if accuracy <= 0 || accuracy > GEO_LENGTH {
		return "", ErrInvalidAccuracy
	}

	step, bitlength := 0, 0
	min_lon, max_lon := MIN_LONGITUDE, MAX_LONGITUDE
	min_lat, max_lat := MIN_LATITUDE, MAX_LATITUDE

	even := true
	var bits uint64 = 0
	result := make([]byte, 0)

	for step < accuracy {
		bits = bits << 1
		if even {
			if m := (max_lon + min_lon) / 2; longitude > m {
				bits = bits | 0x01
				min_lon = m
			} else {
				max_lon = m
			}
		} else {
			if m := (max_lat + min_lat) / 2; latitude > m {
				bits = bits | 0x01
				min_lat = m
			} else {
				max_lat = m
			}
		}
		even = !even
		if bitlength == 4 {
			step++
			bitlength = 0
			result = append(result, BASE_32[bits])
			bits = 0
		} else {
			bitlength++
		}
	}
	return string(result), nil
}


func EncodeToBits(latitude, longitude float64, step uint) (*HashBits, error) {
	if latitude < MIN_LATITUDE || latitude > MAX_LATITUDE || longitude < MIN_LONGITUDE || longitude > MAX_LONGITUDE {
		return nil, ErrInvalidLocation
	}

	var bits uint64 = 0
	min_lon, max_lon := MIN_LONGITUDE, MAX_LONGITUDE
	min_lat, max_lat := MIN_LATITUDE, MAX_LATITUDE
	even := true
	for i := 0; i < int(step)*2; i++ {
		bits = bits << 1
		if even {
			if m := (max_lon + min_lon) / 2; longitude > m {
				bits = bits | 0x01
				min_lon = m
			} else {
				max_lon = m
			}
		} else {
			if m := (max_lat + min_lat) / 2; latitude > m {
				bits = bits | 0x01
				min_lat = m
			} else {
				max_lat = m
			}
		}
		even = !even
	}
	hash := &HashBits{
		Bits: bits,
		Step: step * 2,
	}
	return hash, nil
}


func Decode(hashstr string) (area *Area, e error) {

	hashstr = strings.ToLower(hashstr)
	step := len(hashstr)
	if step == 0 || step > 11 {
		return nil, ErrInvalidAccuracy
	}
	area = &Area{}
	area.Latitude.Min, area.Latitude.Max = MIN_LATITUDE, MAX_LATITUDE
	area.Longitude.Min, area.Longitude.Max = MIN_LONGITUDE, MAX_LONGITUDE
	even := true

	for _, v := range hashstr {
		idx, exist := base32map[v]
		if !exist {
			return nil, ErrInvalidcharacter
		}
		for mask := 1 << 4; mask != 0; mask >>= 1 {
			var r *Range
			if even {
				r = &area.Longitude
			} else {
				r = &area.Latitude
			}
			if mid := r.Mid(); idx&mask != 0 {
				r.Min = mid
			} else {
				r.Max = mid
			}
			even = !even
		}
	}
	return area, nil
}

// reverse the interleave process
// origin from redis author and https://stackoverflow.com/questions/4909263/how-to-efficiently-de-interleave-bits-inverse-morton
func deinterleave(number uint64) (x uint64, y uint64) {
	B := []uint64{0x5555555555555555, 0x3333333333333333, 0x0f0f0f0f0f0f0f0f,
		0x00ff00ff00ff00ff, 0x0000fffff0000ffff, 0x00000000ffffffff}
	S := []uint8{0, 1, 2, 4, 8, 16}

	x, y = number, number>>1

	x = (x | (x >> S[0])) & B[0]
	y = (y | (y >> S[0])) & B[0]

	x = (x | (x >> S[1])) & B[1]
	y = (y | (y >> S[1])) & B[1]

	x = (x | (x >> 2)) & B[2]
	y = (y | (y >> 2)) & B[2]

	x = (x | (x >> 4)) & B[3]
	y = (y | (y >> 4)) & B[3]

	x = (x | (x >> 8)) & B[4]
	y = (y | (y >> 8)) & B[4]

	x = (x | (x >> 16)) & B[5]
	y = (y | (y >> 16)) & B[5]

	return
}

func Encode2(latitude float64, longitude float64, step uint8) uint64 {

	lat_offset := (latitude + 90.0) / 180.0
	lon_offset := (longitude + 180.0) / 360.0

	offset := 1 << step
	lat_offset *= float64(offset)
	lon_offset *= float64(offset)
	fmt.Printf("%b %b\n",uint32(lon_offset),uint32(lat_offset))

	return interleave(uint32(lat_offset), uint32(lon_offset))
}

func interleave(xlo uint32, ylo uint32) uint64 {

	B := []uint64{0x5555555555555555, 0x3333333333333333, 0x0f0f0f0f0f0f0f0f,
		0x00ff00ff00ff00ff, 0x0000fffff0000ffff}
	S := []uint8{1, 2, 4, 8, 16}

	var x, y uint64 = uint64(xlo), uint64(ylo)

	x = (x | (x << S[4])) & B[4]
	y = (y | (y << S[4])) & B[4]

	x = (x | (x << S[3])) & B[3]
	y = (y | (y << S[3])) & B[3]

	x = (x | (x << S[2])) & B[2]
	y = (y | (y << S[2])) & B[2]

	x = (x | (x << S[1])) & B[1]
	y = (y | (y << S[1])) & B[1]

	x = (x | (x << S[0])) & B[0]
	y = (y | (y << S[0])) & B[0]

	return x | (y << 1)

}

func (r *Range) Mid() float64 {
	return (r.Max + r.Min) / 2
}

func (a *Area) ToCoor() (latitude float64, longitude float64) {
	latitude = (a.Latitude.Min + a.Latitude.Max) / 2
	longitude = (a.Longitude.Min + a.Longitude.Max) / 2
	return
}
