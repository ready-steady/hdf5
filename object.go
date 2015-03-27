package hdf5

// #include <stdlib.h>
// #include <hdf5.h>
import "C"

import (
	"unsafe"
)

const (
	flagReference = 1 << iota
	flagOwned
)

type object struct {
	data unsafe.Pointer
	flag uint8

	sid C.hid_t
	tid C.hid_t

	deps []*object
}

func newObject() *object {
	return &object{
		sid: -1,
		tid: -1,
	}
}

func (o *object) free() {
	for i := range o.deps {
		o.deps[i].free()
	}
	if o.tid >= 0 {
		_ = C.H5Tclose(o.tid)
	}
	if o.sid >= 0 {
		_ = C.H5Sclose(o.sid)
	}
	if o.data != nil && o.flag&flagOwned != 0 {
		C.free(o.data)
	}
}
