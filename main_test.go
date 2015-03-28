package hdf5

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/ready-steady/assert"
	"github.com/ready-steady/fixture"
)

func TestPutGet(t *testing.T) {
	path := fixture.MakeFile()
	defer os.Remove(path)

	file, err := Create(path)
	assert.Success(err, t)
	for i, o := range fixtureObjects {
		assert.Success(file.Put(fmt.Sprintf("%c", 'A'+i), o), t)
	}
	assert.Success(file.Close(), t)

	file, err = Open(path)
	assert.Success(err, t)
	for i, o := range fixtureObjects {
		v := reflect.New(reflect.TypeOf(o))
		p := v.Interface()
		assert.Success(file.Get(fmt.Sprintf("%c", 'A'+i), p), t)
		assert.Equal(reflect.Indirect(v).Interface(), o, t)
	}
	assert.Success(file.Close(), t)
}

func TestPutGetWithDimensions(t *testing.T) {
	path := fixture.MakeFile()
	defer os.Remove(path)

	data1 := []float64{
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

	file, err := Create(path)
	assert.Success(err, t)
	assert.Success(file.Put("data", data1, 2, 3, 4), t)
	assert.Success(file.Close(), t)

	data2 := []float64{}

	file, err = Open(path)
	assert.Success(err, t)
	assert.Success(file.Get("data", &data2), t)
	assert.Success(file.Close(), t)

	assert.Equal(data1, data2, t)
}

func TestPutTwice(t *testing.T) {
	path := fixture.MakeFile()
	defer os.Remove(path)

	file, err := Create(path)
	assert.Success(err, t)
	defer file.Close()

	assert.Success(file.Put("A", 42), t)
	assert.Success(file.Put("A", 42), t)
}

func ExampleFile() {
	put := func(path string) {
		file, _ := Create(path)
		defer file.Close()

		A := 42
		file.Put("A", A)

		B := []float64{1, 2, 3}
		file.Put("B", B)

		C := struct {
			D int
			E []float64
		}{
			D: 42,
			E: []float64{1, 2, 3},
		}
		file.Put("C", C)
	}

	get := func(path string) {
		file, _ := Open(path)
		defer file.Close()

		A := 0
		file.Get("A", &A)
		fmt.Println(A)

		B := []float64{}
		file.Get("B", &B)
		fmt.Println(B)

		C := struct {
			D int
			E []float64
		}{}
		file.Get("C", &C)
		fmt.Println(C)
	}

	path := fixture.MakeFile()
	defer os.Remove(path)

	put(path)
	get(path)
	// Output:
	// 42
	// [1 2 3]
	// {42 [1 2 3]}
}
