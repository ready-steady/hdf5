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

// Put writes data into the file.
func (f *File) Put(name string, something interface{}, dimensions ...uint) error {
	object := newObject()
	defer object.free()

	if err := initializeToPut(object, reflect.ValueOf(something), dimensions...); err != nil {
		return err
	}

	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	one := C.hsize_t(1)
	sid := C.H5Screate_simple(1, (*C.hsize_t)(unsafe.Pointer(&one)), nil)
	if sid < 0 {
		return errors.New("cannot create a data space")
	}
	defer C.H5Sclose(sid)

	if result := C.H5Lexists(f.fid, cname, C.H5P_DEFAULT); result < 0 {
		return errors.New("cannot check if the name already exists")
	} else if result > 0 && C.H5Ldelete(f.fid, cname, C.H5P_DEFAULT) < 0 {
		return errors.New("cannot overwrite an existing dataset")
	}

	did := C.H5Dcreate2(f.fid, cname, object.tid, sid, C.H5P_DEFAULT, C.H5P_DEFAULT, C.H5P_DEFAULT)
	if did < 0 {
		return errors.New("cannot create a dataset")
	}
	defer C.H5Dclose(did)

	if C.H5Dwrite(did, object.tid, C.H5S_ALL, C.H5S_ALL, C.H5P_DEFAULT, object.data) < 0 {
		return errors.New("cannot write the dataset into the file")
	}

	return nil
}

func initializeToPut(object *object, value reflect.Value, dimensions ...uint) error {
	switch value.Kind() {
	case reflect.Slice:
		return initializeSliceToPut(object, value, dimensions...)
	case reflect.Struct:
		return initializeStructToPut(object, value)
	default:
		return initializeScalarToPut(object, value)
	}
}

func initializeScalarToPut(object *object, value reflect.Value) error {
	pointer := reflect.New(value.Type())
	reflect.Indirect(pointer).Set(value)

	object.data = unsafe.Pointer(pointer.Pointer())

	bid, ok := kindTypeMapping[value.Kind()]
	if !ok {
		return errors.New("encountered an unsupported datatype")
	}

	one := C.hsize_t(1)
	object.tid = C.H5Tarray_create2(bid, 1, (*C.hsize_t)(unsafe.Pointer(&one)))
	if object.tid < 0 {
		return errors.New("cannot create an array datatype")
	}

	return nil
}

func initializeSliceToPut(object *object, value reflect.Value, dimensions ...uint) error {
	object.data = unsafe.Pointer(value.Pointer())
	object.flag |= flagVariableLength

	bid, ok := kindTypeMapping[value.Type().Elem().Kind()]
	if !ok {
		return errors.New("encountered an unsupported datatype")
	}

	nd := len(dimensions)

	if nd == 0 {
		nd, dimensions = 1, []uint{uint(value.Len())}
	}

	length := uint(1)
	for i := range dimensions {
		length *= dimensions[i]
	}
	if length != uint(value.Len()) {
		return errors.New("the dimensions do not match")
	}

	// NOTE: The C version of HDF5 adheres to the row-major order. This
	// packages, however, favors the column-major order.
	//
	// http://www.hdfgroup.org/HDF5/doc/UG/UG_frame12Dataspaces.html
	for i := 0; i < nd/2; i++ {
		dimensions[i], dimensions[nd-1-i] = dimensions[nd-1-i], dimensions[i]
	}

	object.tid = C.H5Tarray_create2(bid, C.uint(nd), (*C.hsize_t)(unsafe.Pointer(&dimensions[0])))
	if object.tid < 0 {
		return errors.New("cannot create an array datatype")
	}

	return nil
}

func initializeStructToPut(object *object, value reflect.Value) error {
	typo := value.Type()
	size := C.size_t(typo.Size())

	object.data = C.malloc(size)
	if object.data == nil {
		return errors.New("cannot allocate memory")
	}
	object.flag |= flagOwnedMemory

	object.tid = C.H5Tcreate(C.H5T_COMPOUND, size)
	if object.tid < 0 {
		return errors.New("cannot create a compound datatype")
	}

	count := typo.NumField()

	for i := 0; i < count; i++ {
		field := typo.Field(i)
		if len(field.PkgPath) > 0 {
			continue
		}

		o := object.new()
		if err := initializeToPut(o, value.Field(i)); err != nil {
			return err
		}

		address := unsafe.Pointer(uintptr(object.data) + uintptr(field.Offset))

		if o.flag&flagVariableLength != 0 {
			tid := C.H5Tvlen_create(o.tid)
			if tid < 0 {
				return errors.New("cannnot create a variable-length datatype")
			}

			// NOTE: It is assumed here that sizeof(hvl_t) <= v.Type().Size().
			h := (*C.hvl_t)(address)
			h.len, h.p = 1, o.data

			o = object.new()
			o.tid = tid
		} else {
			C.memcpy(address, o.data, C.size_t(value.Field(i).Type().Size()))
		}

		cname := C.CString(field.Name)
		defer C.free(unsafe.Pointer(cname))

		if C.H5Tinsert(object.tid, cname, C.size_t(field.Offset), o.tid) < 0 {
			return errors.New("cannot construct a compound datatype")
		}
	}

	return nil
}
