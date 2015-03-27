package hdf5

/*
#include <hdf5.h>

hid_t _H5T_NATIVE_INT8() { return H5T_NATIVE_INT8; }
hid_t _H5T_NATIVE_UINT8() { return H5T_NATIVE_UINT8; }
hid_t _H5T_NATIVE_INT16() { return H5T_NATIVE_INT16; }
hid_t _H5T_NATIVE_UINT16() { return H5T_NATIVE_UINT16; }
hid_t _H5T_NATIVE_INT32() { return H5T_NATIVE_INT32; }
hid_t _H5T_NATIVE_UINT32() { return H5T_NATIVE_UINT32; }
hid_t _H5T_NATIVE_INT64() { return H5T_NATIVE_INT64; }
hid_t _H5T_NATIVE_UINT64() { return H5T_NATIVE_UINT64; }
hid_t _H5T_NATIVE_FLOAT() { return H5T_NATIVE_FLOAT; }
hid_t _H5T_NATIVE_DOUBLE() { return H5T_NATIVE_DOUBLE; }
*/
import "C"

import (
	"reflect"
)

const (
	is64bit = uint64(^uint(0)) == ^uint64(0)
)

func init() {
	if is64bit {
		kindTypeMapping[reflect.Int] = C._H5T_NATIVE_INT64()
		kindTypeMapping[reflect.Uint] = C._H5T_NATIVE_UINT64()
	} else {
		kindTypeMapping[reflect.Int] = C._H5T_NATIVE_INT32()
		kindTypeMapping[reflect.Uint] = C._H5T_NATIVE_UINT32()
	}
}

var kindTypeMapping = map[reflect.Kind]C.hid_t{
	reflect.Int8:    C._H5T_NATIVE_INT8(),
	reflect.Uint8:   C._H5T_NATIVE_UINT8(),
	reflect.Int16:   C._H5T_NATIVE_INT16(),
	reflect.Uint16:  C._H5T_NATIVE_UINT16(),
	reflect.Int32:   C._H5T_NATIVE_INT32(),
	reflect.Uint32:  C._H5T_NATIVE_UINT32(),
	reflect.Int64:   C._H5T_NATIVE_INT64(),
	reflect.Uint64:  C._H5T_NATIVE_UINT64(),
	reflect.Float32: C._H5T_NATIVE_FLOAT(),
	reflect.Float64: C._H5T_NATIVE_DOUBLE(),
}
