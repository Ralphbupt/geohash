package geohash

import (
	"fmt"
	"testing"
)

func BenchmarkDis(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Distance(35.0, 114.1, 36.0, 114.1)
	}
}
func BenchmarkDisSimle(b *testing.B) {
	for i := 0; i < b.N; i++ {
		//(35.0, 114.1, 36.0, 114.1)
	}
}

func TestDeinterleave(t *testing.T) {

	x, y := deinterleave(0x5555555555555555)
	if x != 0x00000000ffffffff || y != 0 {
		t.Errorf("2222")
	}
	x, y = deinterleave(0xffffffff00000000)
	if x != 0x00000000ffff0000 || y != 0x00000000ffff0000 {
		t.Errorf("errtest111   %b %b", x, y)
	}

	var q, w uint32 = 23, 66
	fmt.Printf("%b %b\n", q, w)
	num := interleave(q, w)
	fmt.Printf("%b\n", num)

	a, s := deinterleave(num)

	fmt.Println(a, s)
}

func TestEnt(t *testing.T) {
	lat, lon := 57.64911, 10.40744

	bits := Encode2(lat, lon,13)
	fmt.Printf("%b\b", bits)

	bits1, _ := EncodeToBits(lat, lon, 26)
	fmt.Printf("%b\n", bits1.Bits)
	//t.Errorf(")

}
