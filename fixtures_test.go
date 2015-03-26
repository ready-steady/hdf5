package hdf5

import (
	"path"
)

const (
	fixturePath = "fixtures"
)

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

	struct {
		A []float64
		B []float64
	}{
		[]float64{1, 2, 3},
		[]float64{4, 5, 6},
	},
}
