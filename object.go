package hdf5

// #include <stdlib.h>
// #include <hdf5.h>
import "C"

import (
	"unsafe"
)

const (
	flagVariableLength = 1 << iota
	flagOwnedMemory
)

type object struct {
	data unsafe.Pointer
	flag uint8
	tid  C.hid_t
	deps []*object
}

func newObject() *object {
	return &object{
		tid: -1,
	}
}

func (o *object) new() *object {
	object := newObject()
	o.deps = append(o.deps, object)
	return object
}

func (o *object) free() {
	for i := range o.deps {
		o.deps[i].free()
	}
	if o.tid >= 0 {
		C.H5Tclose(o.tid)
	}
	if o.data != nil && o.flag&flagOwnedMemory != 0 {
		C.free(o.data)
	}
}
