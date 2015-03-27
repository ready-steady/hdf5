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

// Put writes an object into the file.
func (f *File) Put(name string, something interface{}, dimensions ...uint) error {
	object, err := createObject(reflect.ValueOf(something), dimensions...)
	if err != nil {
		return err
	}
	defer object.free()

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

func createObject(value reflect.Value, dimensions ...uint) (*object, error) {
	switch value.Kind() {
	case reflect.Slice:
		return createArray(value, dimensions...)
	case reflect.Struct:
		return createStruct(value)
	default:
		return createScalar(value)
	}
}

func createArray(value reflect.Value, dimensions ...uint) (*object, error) {
	object := newObject()

	object.data = unsafe.Pointer(value.Pointer())
	object.flag |= flagReference

	bid, ok := kindTypeMapping[value.Type().Elem().Kind()]
	if !ok {
		return nil, errors.New("encountered an unsupported data type")
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
		object.free()
		return nil, errors.New("the dimensions do not match")
	}

	// NOTE: The C version of HDF5 adheres to the row-major order. This
	// packages, however, favors the column-major order.
	//
	// http://www.hdfgroup.org/HDF5/doc/UG/UG_frame12Dataspaces.html
	for i := 0; i < nd/2; i++ {
		dimensions[i], dimensions[nd-1-i] = dimensions[nd-1-i], dimensions[i]
	}

	one := C.hsize_t(1)

	object.sid = C.H5Screate_simple(1, (*C.hsize_t)(unsafe.Pointer(&one)), nil)
	if object.sid < 0 {
		object.free()
		return nil, errors.New("cannot create a data space")
	}

	object.tid = C.H5Tarray_create2(bid, C.uint(nd), (*C.hsize_t)(unsafe.Pointer(&dimensions[0])))
	if object.tid < 0 {
		object.free()
		return nil, errors.New("cannot create a data type")
	}

	return object, nil
}

func createScalar(value reflect.Value) (*object, error) {
	object := newObject()

	pointer := reflect.New(value.Type())
	reflect.Indirect(pointer).Set(value)

	object.data = unsafe.Pointer(pointer.Pointer())

	bid, ok := kindTypeMapping[value.Kind()]
	if !ok {
		object.free()
		return nil, errors.New("encountered an unsupported data type")
	}

	one := C.hsize_t(1)

	object.sid = C.H5Screate_simple(1, (*C.hsize_t)(unsafe.Pointer(&one)), nil)
	if object.sid < 0 {
		object.free()
		return nil, errors.New("cannot create a data space")
	}

	object.tid = C.H5Tarray_create2(bid, 1, (*C.hsize_t)(unsafe.Pointer(&one)))
	if object.tid < 0 {
		object.free()
		return nil, errors.New("cannot create a data type")
	}

	return object, nil
}

func createStruct(value reflect.Value) (*object, error) {
	object := newObject()

	typo := value.Type()
	size := C.size_t(typo.Size())

	object.data = C.malloc(size)
	if object.data == nil {
		return nil, errors.New("cannot allocate memory")
	}
	object.flag |= flagOwned

	pointer := reflect.New(typo)
	reflect.Indirect(pointer).Set(value)
	C.memcpy(object.data, unsafe.Pointer(pointer.Pointer()), size)

	one := C.hsize_t(1)

	object.sid = C.H5Screate_simple(1, (*C.hsize_t)(unsafe.Pointer(&one)), nil)
	if object.sid < 0 {
		object.free()
		return nil, errors.New("cannot create a data space")
	}

	object.tid = C.H5Tcreate(C.H5T_COMPOUND, size)
	if object.tid < 0 {
		object.free()
		return nil, errors.New("cannot create a compound data type")
	}

	count := typo.NumField()

	for i := 0; i < count; i++ {
		field := typo.Field(i)

		o, err := createObject(value.Field(i))
		if err != nil {
			object.free()
			return nil, err
		}
		object.deps = append(object.deps, o)

		tid := o.tid
		offset := C.size_t(field.Offset)

		if o.flag&flagReference != 0 {
			tid = C.H5Tvlen_create(tid)
			if tid < 0 {
				object.free()
				return nil, errors.New("cannnot create a variable-length data type")
			}

			h := (*C.hvl_t)(unsafe.Pointer(uintptr(object.data) + uintptr(offset)))
			h.len, h.p = 1, o.data

			o = newObject()
			o.tid = tid
			object.deps = append(object.deps, o)
		}

		cname := C.CString(field.Name)
		defer C.free(unsafe.Pointer(cname))

		if C.H5Tinsert(object.tid, cname, offset, tid) < 0 {
			object.free()
			return nil, errors.New("cannot construct a compound data type")
		}
	}

	return object, nil
}
