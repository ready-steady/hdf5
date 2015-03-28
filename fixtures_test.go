package hdf5

import (
	"path"
)

const (
	fixturePath = "fixtures"
)

type dummy1 struct {
	A uint
	B []uint
	C dummy2
}

type dummy2 struct {
	D uint
	E []uint
}

func findFixture(name string) string {
	return path.Join(fixturePath, name)
}

var fixtureObjects = []interface{}{
	int(-1),
	[]int{-1, -1, -1},

	uint(2),
	[]uint{2, 2, 2},

	int8(-3),
	[]int8{-3, -3, -3},

	uint8(4),
	[]uint8{4, 4, 4},

	int16(-5),
	[]int16{-5, -5, -5},

	uint16(6),
	[]uint16{6, 6, 6},

	int32(-7),
	[]int32{-7, -7, -7},

	uint32(8),
	[]uint32{8, 8, 8},

	int64(-9),
	[]int64{-9, -9, -9},

	uint64(10),
	[]uint64{10, 10, 10},

	float32(11),
	[]float32{11, 11, 11},

	float64(12),
	[]float64{12, 12, 12},

	dummy1{
		A: 1,
		B: []uint{2, 3},
		C: dummy2{
			D: 4,
			E: []uint{5, 6},
		},
	},
}
