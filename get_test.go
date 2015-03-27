package hdf5

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/ready-steady/assert"
)

func TestGet(t *testing.T) {
	path := findFixture("data.h5")

	file, _ := Open(path)
	defer file.Close()

	for i, o := range fixtureObjects {
		v := reflect.New(reflect.TypeOf(o))
		p := v.Interface()
		assert.Success(file.Get(fmt.Sprintf("%c", 'A'+i), p), t)
		assert.Equal(reflect.Indirect(v).Interface(), o, t)
	}
}
