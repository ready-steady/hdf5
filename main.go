// Package hdf5 provides a reader and writer of HDF5 files.
//
// https://en.wikipedia.org/wiki/Hierarchical_Data_Format
package hdf5

// #cgo CFLAGS: -Ihdf5/install/include
// #cgo LDFLAGS: -lm -lz
//
// #include <stdlib.h>
// #include <string.h>
//
// #include <hdf5.h>
// #include <hdf5_hl.h>
import "C"

import (
	"errors"
	"unsafe"
)

// File represents a file.
type File struct {
	id C.hid_t
}

// Open opens a file for reading and writing.
func Open(path string) (*File, error) {
	const (
		// https://github.com/copies/hdf5/blob/master/src/H5Fpublic.h#L46
		F_ACC_RDWR = 0x0001
	)

	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	id := C.H5Fopen(cpath, F_ACC_RDWR, C.H5P_DEFAULT)
	if id < 0 {
		return nil, errors.New("failed to open the file")
	}

	file := &File{
		id: id,
	}

	return file, nil
}

// Close closes the file.
func (f *File) Close() error {
	if err := C.H5Fclose(f.id); err != 0 {
		return errors.New("failed to close the file")
	}
	return nil
}
