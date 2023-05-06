package msi

import (
	"reflect"
	"unsafe"
)

var (
	procMsiSummaryInfoGetProperty = modMsi.NewProc("MsiSummaryInfoGetPropertyA")
	procMsiSummaryInfoSetProperty = modMsi.NewProc("MsiSummaryInfoSetPropertyA")
	procMsiSummaryInfoPersist     = modMsi.NewProc("MsiSummaryInfoPersist")
)

const (
	PID_CODEPAGE     = 1  // VT_I2
	PID_TITLE        = 2  // VT_LPSTR
	PID_SUBJECT      = 3  // VT_LPSTR
	PID_AUTHOR       = 4  // VT_LPSTR
	PID_KEYWORDS     = 5  // VT_LPSTR
	PID_COMMENTS     = 6  // VT_LPSTR
	PID_TEMPLATE     = 7  // VT_LPSTR
	PID_LASTAUTHOR   = 8  // VT_LPSTR
	PID_REVNUMBER    = 9  // VT_LPSTR
	PID_LASTPRINTED  = 11 // VT_FILETIME
	PID_CREATE_DTM   = 12 // VT_FILETIME
	PID_LASTSAVE_DTM = 13 // VT_FILETIME
	PID_PAGECOUNT    = 14 // VT_I4
	PID_WORDCOUNT    = 15 // VT_I4
	PID_CHARCOUNT    = 16 // VT_I4
	PID_APPNAME      = 18 // VT_LPSTR
	PID_SECURITY     = 19 // VT_I4
)

const (
	VT_I2       = 2
	VT_I4       = 3
	VT_LPSTR    = 30
	VT_FILETIME = 64
)

var typeMap = map[int]int{
	PID_CODEPAGE:     VT_I2,
	PID_TITLE:        VT_LPSTR,
	PID_SUBJECT:      VT_LPSTR,
	PID_AUTHOR:       VT_LPSTR,
	PID_KEYWORDS:     VT_LPSTR,
	PID_COMMENTS:     VT_LPSTR,
	PID_TEMPLATE:     VT_LPSTR,
	PID_LASTAUTHOR:   VT_LPSTR,
	PID_REVNUMBER:    VT_LPSTR,
	PID_LASTPRINTED:  VT_FILETIME,
	PID_CREATE_DTM:   VT_FILETIME,
	PID_LASTSAVE_DTM: VT_FILETIME,
	PID_PAGECOUNT:    VT_I4,
	PID_WORDCOUNT:    VT_I4,
	PID_CHARCOUNT:    VT_I4,
	PID_APPNAME:      VT_LPSTR,
	PID_SECURITY:     VT_I4,
}

type SummaryInformation struct {
	handle uintptr
}

func (info SummaryInformation) Persist() error {
	ret, _, err := procMsiSummaryInfoPersist.Call(
		info.handle,
	)
	if ret != 0 {
		return err
	}
	return nil
}

func (info *SummaryInformation) Close() error {
	if info.handle == 0 {
		return nil
	}

	if ret, _, err := procMsiCloseHandle.Call(info.handle); ret != 0 {
		return err
	}

	info.handle = 0
	return nil
}

func (info SummaryInformation) GetString(prop int) (string, error) {
	var valueBuf [256]byte
	valueBufLen := uintptr(len(valueBuf))
	ret, _, err := procMsiSummaryInfoGetProperty.Call(
		info.handle,
		uintptr(prop),
		0,
		uintptr(unsafe.Pointer(nil)),
		uintptr(unsafe.Pointer(nil)),
		uintptr(unsafe.Pointer(&valueBuf[0])),
		uintptr(unsafe.Pointer(&valueBufLen)))
	if ret != 0 {
		return "", err
	}
	return string(valueBuf[:valueBufLen]), nil
}

func (info SummaryInformation) SetProperty(prop int, value any) error {
	typeId := typeMap[prop]
	switch typeId {
	case VT_LPSTR:
		b := append([]byte(value.(string)), 0)
		ret, _, err := procMsiSummaryInfoSetProperty.Call(
			info.handle,
			uintptr(prop),
			VT_LPSTR,
			0,
			uintptr(unsafe.Pointer(nil)),
			uintptr(unsafe.Pointer(&b[0])),
		)
		if ret == 0 {
			return err
		}
	case VT_I2, VT_I4:
		var v int32
		switch reflect.TypeOf(value).Kind() {
		case reflect.Int:
			v = int32(value.(int))
		case reflect.Int16:
			v = int32(value.(int16))
		case reflect.Int32:
			v = int32(value.(int32))
		}
		ret, _, err := procMsiSummaryInfoSetProperty.Call(
			info.handle,
			uintptr(prop),
			uintptr(typeId),
			uintptr(v),
			uintptr(unsafe.Pointer(nil)),
			uintptr(unsafe.Pointer(nil)),
		)
		if ret == 0 {
			return err
		}
	case VT_FILETIME:
		v := value.(int64)
		ret, _, err := procMsiSummaryInfoSetProperty.Call(
			info.handle,
			uintptr(prop),
			VT_FILETIME,
			0,
			uintptr(unsafe.Pointer(&v)),
			uintptr(unsafe.Pointer(nil)),
		)
		if ret == 0 {
			return err
		}
	}
	return nil
}
