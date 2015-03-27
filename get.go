package hdf5

// #include <stdlib.h>
// #include <string.h>
// #include <hdf5.h>
import "C"

import (
	"errors"
	"reflect"
	"unsafe"
)

// Get reads data from the file.
func (f *File) Get(name string, something interface{}) error {
	value := reflect.ValueOf(something)
	if value.Kind() != reflect.Ptr {
		return errors.New("expected a pointer")
	}

	ivalue := reflect.Indirect(value)

	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	did := C.H5Dopen2(f.fid, cname, C.H5P_DEFAULT)
	if did < 0 {
		return errors.New("cannot find the dataset")
	}
	defer C.H5Dclose(did)

	object := newObject()
	defer object.free()

	object.tid = C.H5Dget_type(did)
	if object.tid < 0 {
		return errors.New("cannot get the datatype of the dataset")
	}

	if err := initializeObject(object, ivalue); err != nil {
		return err
	}

	if C.H5Dread(did, object.tid, C.H5S_ALL, C.H5S_ALL, C.H5P_DEFAULT, object.data) != 0 {
		return errors.New("cannot read the dataset from the file")
	}

	if err := finalizeObject(object, ivalue); err != nil {
		return err
	}

	return nil
}

func initializeObject(object *object, value reflect.Value) error {
	switch value.Kind() {
	case reflect.Slice:
		return initializeSlice(object, value)
	case reflect.Struct:
		return initializeStruct(object, value)
	default:
		return initializeScalar(object, value)
	}
}

func finalizeObject(object *object, value reflect.Value) error {
	switch value.Kind() {
	case reflect.Struct:
		return finalizeStruct(object, value)
	default:
		return nil
	}
}

func initializeScalar(object *object, value reflect.Value) error {
	bid, ok := kindTypeMapping[value.Kind()]
	if !ok {
		return errors.New("encountered an unsupported datatype")
	}

	if err := checkArrayType(object.tid, bid); err != nil {
		return err
	}

	if length, err := getArrayLength(object.tid); err != nil {
		return err
	} else if length != 1 {
		return errors.New("expected an array with a single element")
	}

	object.data = unsafe.Pointer(value.Addr().Pointer())

	return nil
}

func initializeSlice(object *object, value reflect.Value) error {
	typo := value.Type()

	bid, ok := kindTypeMapping[typo.Elem().Kind()]
	if !ok {
		return errors.New("encountered an unsupported datatype")
	}

	if err := checkArrayType(object.tid, bid); err != nil {
		return err
	}

	length, err := getArrayLength(object.tid)
	if err != nil {
		return err
	}

	buffer := reflect.MakeSlice(typo, int(length), int(length))
	shadow := reflect.Indirect(reflect.New(typo))
	shadow.Set(buffer)

	src := (*reflect.SliceHeader)(unsafe.Pointer(shadow.UnsafeAddr()))
	dst := (*reflect.SliceHeader)(unsafe.Pointer(value.UnsafeAddr()))

	dst.Data, src.Data = src.Data, dst.Data
	dst.Cap, src.Cap = src.Cap, dst.Cap
	dst.Len, src.Len = src.Len, dst.Len

	object.data = unsafe.Pointer(dst.Data)

	return nil
}

func initializeStruct(object *object, value reflect.Value) error {
	return nil
}

func finalizeStruct(object *object, value reflect.Value) error {
	return nil
}

func checkArrayType(tid C.hid_t, bid C.hid_t) error {
	if id := C.H5Tget_class(tid); id < 0 {
		return errors.New("cannot get the data class of a datatype")
	} else if id != C.H5T_ARRAY {
		return errors.New("expected an array")
	}

	if id := C.H5Tget_super(tid); id < 0 {
		return errors.New("cannot get the base type of a datatype")
	} else if C.H5Tequal(bid, id) == 0 {
		return errors.New("the types do not match")
	}

	return nil
}

func getArrayLength(tid C.hid_t) (C.hsize_t, error) {
	nd := C.H5Tget_array_ndims(tid)
	if nd < 0 {
		return 0, errors.New("cannot get the dimensionality of an array")
	}

	dimensions := make([]C.hsize_t, nd)
	if C.H5Tget_array_dims2(tid, (*C.hsize_t)(unsafe.Pointer(&dimensions[0]))) != nd {
		return 0, errors.New("cannot get the dimensions of an array")
	}

	length := C.hsize_t(1)
	for i := range dimensions {
		length *= dimensions[i]
	}

	return length, nil
}
