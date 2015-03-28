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

	value = reflect.Indirect(value)

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

	if err := initializeToGet(object, value); err != nil {
		return err
	}

	if C.H5Dread(did, object.tid, C.H5S_ALL, C.H5S_ALL, C.H5P_DEFAULT, object.data) != 0 {
		return errors.New("cannot read the dataset from the file")
	}

	if err := finalizeToGet(object, value); err != nil {
		return err
	}

	return nil
}

func initializeToGet(object *object, value reflect.Value) error {
	switch value.Kind() {
	case reflect.Slice:
		return initializeSliceToGet(object, value)
	case reflect.Struct:
		return initializeStructToGet(object, value)
	default:
		return initializeScalarToGet(object, value)
	}
}

func finalizeToGet(object *object, value reflect.Value) error {
	switch value.Kind() {
	case reflect.Struct:
		return finalizeStructToGet(object, value)
	default:
		return nil
	}
}

func initializeScalarToGet(object *object, value reflect.Value) error {
	bid, ok := kindTypeMapping[value.Kind()]
	if !ok {
		return errors.New("encountered an unsupported datatype")
	}

	if err := checkArrayType(object.tid, bid); err != nil {
		return err
	}

	if length, err := computeArrayLength(object.tid); err != nil {
		return err
	} else if length != 1 {
		return errors.New("expected an array with a single element")
	}

	object.data = unsafe.Pointer(value.Addr().Pointer())

	return nil
}

func initializeSliceToGet(object *object, value reflect.Value) error {
	typo := value.Type()

	bid, ok := kindTypeMapping[typo.Elem().Kind()]
	if !ok {
		return errors.New("encountered an unsupported datatype")
	}

	if err := checkArrayType(object.tid, bid); err != nil {
		return err
	}

	length, err := computeArrayLength(object.tid)
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
	object.flag |= flagVariableLength

	return nil
}

func initializeStructToGet(object *object, value reflect.Value) error {
	if tid := C.H5Tget_class(object.tid); tid < 0 {
		return errors.New("cannot get a data class")
	} else if tid != C.H5T_COMPOUND {
		return errors.New("expected a compound datatype")
	}

	size := C.H5Tget_size(object.tid)
	if size < 0 {
		return errors.New("cannot get the size of a compound datatype")
	}

	object.data = C.malloc(size)
	if object.data == nil {
		return errors.New("cannot allocate memory")
	}
	object.flag |= flagOwnedMemory

	return nil
}

func finalizeStructToGet(object *object, value reflect.Value) error {
	typo := value.Type()
	count := typo.NumField()

	for i := 0; i < count; i++ {
		field := typo.Field(i)

		cname := C.CString(field.Name)
		defer C.free(unsafe.Pointer(cname))

		j := C.H5Tget_member_index(object.tid, cname)
		if j < 0 {
			continue
		}

		o := object.new()

		o.tid = C.H5Tget_member_type(object.tid, C.uint(j))
		if o.tid < 0 {
			return errors.New("cannot get the datatype of a field")
		}

		if cid := C.H5Tget_class(o.tid); cid < 0 {
			return errors.New("cannot get the data class of a field")
		} else if cid == C.H5T_VLEN {
			if tid := C.H5Tget_super(o.tid); tid < 0 { // Close?
				return errors.New("cannot get the base type of a field")
			} else {
				o = object.new()
				o.tid = tid
			}
		}

		if err := initializeToGet(o, value.Field(i)); err != nil {
			return err
		}

		size := C.H5Tget_size(o.tid)
		if size < 0 {
			return errors.New("cannot get a size")
		}

		offset := C.H5Tget_member_offset(object.tid, C.uint(j))
		address := unsafe.Pointer(uintptr(object.data) + uintptr(offset))

		if o.flag&flagVariableLength != 0 {
			h := (*C.hvl_t)(address)
			if h.len != 1 {
				return errors.New("expected a variable-length datatype with a single element")
			}
			address = h.p
		}

		C.memcpy(o.data, address, size)

		if err := finalizeToGet(o, value.Field(i)); err != nil {
			return err
		}
	}

	return nil
}

func checkArrayType(tid C.hid_t, bid C.hid_t) error {
	if cid := C.H5Tget_class(tid); cid < 0 {
		return errors.New("cannot get a data class")
	} else if cid != C.H5T_ARRAY {
		return errors.New("expected an array datatype")
	}

	if tid := C.H5Tget_super(tid); tid < 0 { // Close?
		return errors.New("cannot get the base type of a datatype")
	} else if C.H5Tequal(bid, tid) == 0 {
		return errors.New("the types do not match")
	}

	return nil
}

func computeArrayLength(tid C.hid_t) (C.hsize_t, error) {
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
