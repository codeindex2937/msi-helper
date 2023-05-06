package msi

import (
	"reflect"
	"unsafe"
)

var (
	procMsiCreateRecord     = modMsi.NewProc("MsiCreateRecord")
	procMsiRecordSetInteger = modMsi.NewProc("MsiRecordSetInteger")
	procMsiRecordSetString  = modMsi.NewProc("MsiRecordSetStringA")
	procMsiRecordSetStream  = modMsi.NewProc("MsiRecordSetStreamA")
)

type Record struct {
	handle  uintptr
	numCols int
}

func NewRecord(numCols int) (Record, error) {
	handle := uintptr(0)
	handle, _, err := procMsiCreateRecord.Call(
		uintptr(numCols),
	)
	if handle == 0 {
		return Record{}, err
	}

	return Record{handle: handle, numCols: numCols}, nil
}

func (r *Record) Close() error {
	if r.handle == 0 {
		return nil
	}

	if ret, _, err := procMsiCloseHandle.Call(r.handle); ret != 0 {
		return err
	}

	r.handle = 0
	return nil
}

func (r *Record) setRecord(col int, t reflect.StructField, v reflect.Value) error {
	i := v.Interface()
	switch t.Type {
	case nullStringType:
		nullStr := i.(NullString)
		if nullStr.Valid {
			if err := r.setString(col, nullStr.String); err != nil {
				return err
			}
		}
	case nullInt16Type:
		nullInt16 := i.(NullInt16)
		if nullInt16.Valid {
			if err := r.setInteger(col, int(nullInt16.Int16)); err != nil {
				return err
			}
		}
	case nullInt32Type:
		nullInt32 := i.(NullInt32)
		if nullInt32.Valid {
			if err := r.setInteger(col, int(nullInt32.Int32)); err != nil {
				return err
			}
		}
	case streamType:
		if err := r.setStream(col, i.(Stream)); err != nil {
			return err
		}
	default:
		switch t.Type.Kind() {
		case reflect.String:
			if err := r.setString(col, i.(string)); err != nil {
				return err
			}
		case reflect.Int16:
			if err := r.setInteger(col, int(i.(int16))); err != nil {
				return err
			}
		case reflect.Int32:
			if err := r.setInteger(col, int(i.(int32))); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *Record) setString(col int, value string) error {
	v := append([]byte(value), 0)
	ret, _, err := procMsiRecordSetString.Call(
		r.handle,
		uintptr(col),
		uintptr(unsafe.Pointer(&v[0])),
	)
	if ret != 0 {
		return err
	}

	return nil
}

func (r *Record) setInteger(col int, value int) error {
	ret, _, err := procMsiRecordSetInteger.Call(
		r.handle,
		uintptr(col),
		uintptr(value),
	)
	if ret != 0 {
		return err
	}

	return nil
}

func (r *Record) setStream(col int, path Stream) error {
	ret, _, err := procMsiRecordSetStream.Call(
		r.handle,
		uintptr(col),
		uintptr(unsafe.Pointer(&[]byte(path)[0])),
	)
	if ret != 0 {
		return err
	}

	return nil
}
