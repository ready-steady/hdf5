package hdf5

// #include <stdlib.h>
// #include <hdf5.h>
import "C"

import (
	"unsafe"
)

type object struct {
	data unsafe.Pointer
	sid  C.hid_t
	tid  C.hid_t
}

func (o *object) free() {
	_ = C.H5Tclose(o.tid)
	_ = C.H5Sclose(o.sid)
}
