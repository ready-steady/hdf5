package hdf5

import (
	"fmt"
	"os"
	"testing"

	"github.com/ready-steady/assert"
	"github.com/ready-steady/fixture"
)

func TestPut(t *testing.T) {
	path := fixture.MakeFile()
	defer os.Remove(path)

	file, _ := Create(path)
	defer file.Close()

	for i, o := range fixtureObjects {
		assert.Success(file.Put(fmt.Sprintf("%c", 'A'+i), o), t)
	}
}

func TestPutMatrix(t *testing.T) {
	path := fixture.MakeFile()
	defer os.Remove(path)

	file, _ := Create(path)
	defer file.Close()

	data := []float64{
		0.01, 0.02,
		0.03, 0.04,
		0.05, 0.06,

		0.07, 0.08,
		0.09, 0.10,
		0.11, 0.12,

		0.13, 0.14,
		0.15, 0.16,
		0.17, 0.18,

		0.19, 0.20,
		0.21, 0.22,
		0.23, 0.24,
	}

	assert.Success(file.Put("M", data, 2, 3, 4), t)
}
