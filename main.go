// Package hdf5 provides a reader and writer of HDF5 files.
//
// https://en.wikipedia.org/wiki/Hierarchical_Data_Format
package hdf5

/*
#cgo CFLAGS: -Ihdf5/install/include
#cgo LDFLAGS: -ldl -lm -lz

#include <stdlib.h>
#include <hdf5.h>

uint _H5F_ACC_TRUNC() { return H5F_ACC_TRUNC; }
uint _H5F_ACC_RDWR() { return H5F_ACC_RDWR; }
*/
import "C"

import (
	"errors"
	"unsafe"
)

// File represents a file.
type File struct {
	fid C.hid_t
}

// Create creates a new file.
func Create(path string) (*File, error) {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	fid := C.H5Fcreate(cpath, C._H5F_ACC_TRUNC(), C.H5P_DEFAULT, C.H5P_DEFAULT)
	if fid < 0 {
		return nil, errors.New("failed to create a file")
	}

	file := &File{
		fid: fid,
	}

	return file, nil
}

// Open opens an existing file.
func Open(path string) (*File, error) {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	fid := C.H5Fopen(cpath, C._H5F_ACC_RDWR(), C.H5P_DEFAULT)
	if fid < 0 {
		return nil, errors.New("failed to open the file")
	}

	file := &File{
		fid: fid,
	}

	return file, nil
}

// Close closes the file.
func (f *File) Close() error {
	if err := C.H5Fclose(f.fid); err != 0 {
		return errors.New("failed to close the file")
	}
	return nil
}
