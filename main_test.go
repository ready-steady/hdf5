package hdf5

import (
	"os"
	"testing"

	"github.com/ready-steady/assert"
	"github.com/ready-steady/fixture"
)

func TestOpen(t *testing.T) {
	path := fixture.MakeFile()
	defer os.Remove(path)

	file, err := Open(path)
	assert.Success(err, t)

	err = file.Close()
	assert.Success(err, t)
}
