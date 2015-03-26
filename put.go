package hdf5

// #include <stdlib.h>
// #include <hdf5.h>
import "C"

import (
	"errors"
	"reflect"
	"unsafe"
)

// Put writes an object into the file.
func (f *File) Put(name string, something interface{}, dimensions ...uint) error {
	object, err := f.createObject(reflect.ValueOf(something), dimensions...)
	if err != nil {
		return err
	}
	defer object.free()

	return f.putObject(name, object)
}

func (f *File) createObject(value reflect.Value, dimensions ...uint) (*object, error) {
	switch value.Kind() {
	case reflect.Slice:
		return f.createArray(value, dimensions...)
	case reflect.Struct:
		return f.createStruct(value)
	default:
		return f.createScalar(value)
	}
}

func (f *File) createArray(value reflect.Value, dimensions ...uint) (*object, error) {
	data := unsafe.Pointer(value.Pointer())

	bid, ok := kindTypeMapping[value.Type().Elem().Kind()]
	if !ok {
		return nil, errors.New("encountered an unsupported data type")
	}

	if len(dimensions) == 0 {
		dimensions = []uint{uint(value.Len())}
	}

	length := uint(1)
	for i := range dimensions {
		length *= dimensions[i]
	}
	if length != uint(value.Len()) {
		return nil, errors.New("the dimensions do not match")
	}

	one := C.hsize_t(1)

	sid := C.H5Screate_simple(1, (*C.hsize_t)(unsafe.Pointer(&one)), nil)
	if sid < 0 {
		return nil, errors.New("cannot create a data space")
	}

	tid := C.H5Tarray_create2(bid, C.uint(len(dimensions)),
		(*C.hsize_t)(unsafe.Pointer(&dimensions[0])))
	if tid < 0 {
		_ = C.H5Sclose(sid)
		return nil, errors.New("cannot create a data type")
	}

	object := &object{
		data: data,
		sid:  sid,
		tid:  tid,
	}

	return object, nil
}

func (f *File) createScalar(value reflect.Value) (*object, error) {
	pointer := reflect.New(value.Type())
	reflect.Indirect(pointer).Set(value)

	data := unsafe.Pointer(pointer.Pointer())

	bid, ok := kindTypeMapping[value.Kind()]
	if !ok {
		return nil, errors.New("encountered an unsupported data type")
	}

	one := C.hsize_t(1)

	sid := C.H5Screate_simple(1, (*C.hsize_t)(unsafe.Pointer(&one)), nil)
	if sid < 0 {
		return nil, errors.New("cannot create a data space")
	}

	tid := C.H5Tarray_create2(bid, 1, (*C.hsize_t)(unsafe.Pointer(&one)))
	if tid < 0 {
		_ = C.H5Sclose(sid)
		return nil, errors.New("cannot create a data type")
	}

	object := &object{
		data: data,
		sid:  sid,
		tid:  tid,
	}

	return object, nil
}

func (f *File) createStruct(value reflect.Value) (*object, error) {
	typo := value.Type()
	pointer := reflect.New(typo)
	reflect.Indirect(pointer).Set(value)

	data := unsafe.Pointer(pointer.Pointer())

	one := C.hsize_t(1)

	sid := C.H5Screate_simple(1, (*C.hsize_t)(unsafe.Pointer(&one)), nil)
	if sid < 0 {
		return nil, errors.New("cannot create a data space")
	}

	tid := C.H5Tcreate(C.H5T_COMPOUND, C.size_t(typo.Size()))
	if tid < 0 {
		_ = C.H5Sclose(sid)
		return nil, errors.New("cannot create a compound data type")
	}

	count := typo.NumField()

	for i := 0; i < count; i++ {
		field := typo.Field(i)

		object, err := f.createObject(value.Field(i))
		if err != nil {
			_ = C.H5Tclose(tid)
			_ = C.H5Sclose(sid)
			return nil, err
		}

		cname := C.CString(field.Name)
		defer C.free(unsafe.Pointer(cname))

		if C.H5Tinsert(tid, cname, C.size_t(field.Offset), object.tid) != 0 {
			_ = C.H5Tclose(tid)
			_ = C.H5Sclose(sid)
			return nil, errors.New("cannot construct a compound data type")
		}
	}

	object := &object{
		data: data,
		sid:  sid,
		tid:  tid,
	}

	return object, nil
}

func (f *File) putObject(name string, object *object) error {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	did := C.H5Dcreate2(f.fid, cname, object.tid, object.sid,
		C.H5P_DEFAULT, C.H5P_DEFAULT, C.H5P_DEFAULT)
	if did < 0 {
		return errors.New("cannot create a dataset")
	}
	defer C.H5Dclose(did)

	if C.H5Dwrite(did, object.tid, C.H5S_ALL, C.H5S_ALL, C.H5P_DEFAULT, object.data) != 0 {
		return errors.New("cannot write a dataset into the file")
	}

	return nil
}
