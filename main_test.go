package hdf5

import (
	"fmt"
	"os"
	"testing"

	"github.com/ready-steady/assert"
	"github.com/ready-steady/fixture"
)

func TestOpenClose(t *testing.T) {
	path := fixture.MakeFile()
	defer os.Remove(path)

	file, err := Open(path)
	assert.Success(err, t)

	err = file.Close()
	assert.Success(err, t)
}

func TestPut(t *testing.T) {
	path := fixture.MakeFile()
	defer os.Remove(path)

	file, _ := Create(path)
	defer file.Close()

	for i, o := range fixtureObjects {
		assert.Success(file.Put(fmt.Sprintf("%c", 'A'+i), o), t)
	}
}
