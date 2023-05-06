package msi

import (
	"fmt"
	"reflect"
	"strings"
	"unsafe"

	"github.com/google/uuid"
	"golang.org/x/sys/windows"
)

const (
	MSIDBOPEN_READONLY = iota
	MSIDBOPEN_TRANSACT
	MSIDBOPEN_DIRECT
	MSIDBOPEN_CREATE
	MSIDBOPEN_CREATEDIRECT

	UPGRADE_ATTRIBUTES_VERSION_ONLY_DETECT   = 0x002
	UPGRADE_ATTRIBUTES_VERSION_MAX_INCLUSIVE = 0x200

	COMPONENT_ATTRIBUTES_64BIT = 0x0100

	FILE_ATTRIBUTES_COMPRESSED     = 0x004000
	FILE_ATTRIBUTES_NON_COMPRESSED = 0x002000

	CUSTOM_ACTION_TYPE_DLL      = 1
	CUSTOM_ACTION_TYPE_VBSCRIPT = 6

	CUSTOM_ACTION_TYPE_EXE_IN_FILE            = 18
	CUSTOM_ACTION_TYPE_ALERT_AND_EXIT         = 19
	CUSTOM_ACTION_TYPE_SET_FORMATTED_PROPERTY = 51

	CUSTOM_ACTION_TYPE_CONTINUE       = 0x00000040
	CUSTOM_ACTION_TYPE_ASYNC          = 0x00000080
	CUSTOM_ACTION_TYPE_IN_SCRIPT      = 0x00000400
	CUSTOM_ACTION_TYPE_NO_IMPERSONATE = 0x00000800
	CUSTOM_ACTION_TYPE_HIDE_TARGET    = 0x00002000
	CUSTOM_ACTION_TYPE_FIRST_SEQUENCE = 0x00000100

	CONTROL_ATTRIBUTES_ENABLED     = 0x00000001
	CONTROL_ATTRIBUTES_VISIBLE     = 0x00000002
	CONTROL_ATTRIBUTES_TRANSPARENT = 0x00010000
	CONTROL_ATTRIBUTES_NO_PREFIX   = 0x00020000

	SERVICECTRL_EVENT_START            = 0x001
	SERVICECTRL_EVENT_STOP             = 0x002
	SERVICECTRL_EVENT_DELETE           = 0x008
	SERVICECTRL_EVENT_UNINSTALL_STOP   = 0x020
	SERVICECTRL_EVENT_UNINSTALL_DELETE = 0x080

	SERVICE_TYPE_OWN_PROCESS = 0x00000010

	SERVICE_START_TYPE_AUTO = 0x00000002

	SERVICE_ERROR_CONTROL_NORMAL = 0x00000001

	CONDITION_PRE_INSTALLED      = "UPGRADE_FOUND <> \"\" OR NEWER_INSTALLED <> \"\""
	CONDITION_PRE_CLEAN_INSTALL  = "UPGRADE_FOUND = \"\" AND NEWER_INSTALLED = \"\""
	CONDITION_POST_INSTALLED     = "Installed OR UPGRADE_FOUND OR NEWER_INSTALLED"
	CONDITION_POST_CLEAN_INSTALL = "NOT Installed AND NOT UPGRADE_FOUND AND NOT NEWER_INSTALLED"
	CONDITION_INSTALL_OR_UPGRADE = "(NOT REMOVE) OR (UPGRADINGPRODUCTCODE)"
	CONDITION_REMOVE             = "(REMOVE=\"ALL\") AND (NOT UPGRADINGPRODUCTCODE)"
	CONDITION_POST_UPGRADE       = "UPGRADE_FOUND"
)

var (
	modMsi                            = windows.NewLazyDLL("msi.dll")
	procMsiOpenDatabase               = modMsi.NewProc("MsiOpenDatabaseA")
	procMsiCloseHandle                = modMsi.NewProc("MsiCloseHandle")
	procMsiDatabaseCommit             = modMsi.NewProc("MsiDatabaseCommit")
	procMsiGetSummaryInformation      = modMsi.NewProc("MsiGetSummaryInformationA")
	procMsiDatabaseOpenView           = modMsi.NewProc("MsiDatabaseOpenViewA")
	procMsiDatabaseImport             = modMsi.NewProc("MsiDatabaseImportA")
	procMsiDatabaseGenerateTransform  = modMsi.NewProc("MsiDatabaseGenerateTransformA")
	procMsiCreateTransformSummaryInfo = modMsi.NewProc("MsiCreateTransformSummaryInfoA")
)

type Database struct {
	handle uintptr
}

func NewUUID() string {
	return fmt.Sprintf("{%v}", strings.ToUpper(uuid.NewString()))
}

func Open(name string, flag int) (Database, error) {
	dbHandle := uintptr(0)
	ret, _, err := procMsiOpenDatabase.Call(
		uintptr(unsafe.Pointer(&[]byte(name)[0])),
		uintptr(flag),
		uintptr(unsafe.Pointer(&dbHandle)))
	if ret != 0 {
		return Database{}, err
	}

	return Database{dbHandle}, nil
}

func (db *Database) Close() error {
	if db.handle == 0 {
		return nil
	}

	if ret, _, err := procMsiCloseHandle.Call(db.handle); ret != 0 {
		return err
	}

	db.handle = 0
	return nil
}

func (db Database) Commit() error {
	if ret, _, err := procMsiDatabaseCommit.Call(db.handle); ret != 0 {
		return err
	}
	return nil
}

func (db Database) OpenSummaryInformation(count int) (SummaryInformation, error) {
	handle := uintptr(0)

	ret, _, err := procMsiGetSummaryInformation.Call(
		db.handle,
		uintptr(unsafe.Pointer(nil)),
		uintptr(count),
		uintptr(unsafe.Pointer(&handle)),
	)
	if ret != 0 {
		return SummaryInformation{}, err
	}
	return SummaryInformation{handle}, nil
}

func (db Database) Insert(rows ...interface{}) error {
	for _, r := range rows {
		t := reflect.TypeOf(r)
		v := reflect.ValueOf(r)
		numField := keyCountMap[t.Name()][1]
		rec, err := NewRecord(numField)
		if err != nil {
			return err
		}
		defer rec.Close()

		for i := 0; i < t.NumField(); i++ {
			if info := fieldMap[t.Name()][t.Field(i).Name]; !info.isField {
				continue
			}
			if err := rec.setRecord(i+1, t.Field(i), v.Field(i)); err != nil {
				return err
			}
		}

		sql := fmt.Sprintf("INSERT INTO `%v` (%v) VALUES (%v)",
			tableMap[t.Name()],
			fieldsMap[t.Name()],
			placeholdersMap[t.Name()],
		)
		view, err := db.OpenView(sql)
		if err != nil {
			return err
		}
		defer view.Close()

		view.Execute(rec)
	}
	return nil
}

func (db Database) Update(data interface{}, fields []string, conds map[string]interface{}) error {
	var rec Record
	placeholders := []string{}
	t := reflect.TypeOf(data)
	v := reflect.ValueOf(data)

	numField := 0
	if len(fields) > 0 {
		numField += len(fields)
		if len(conds) == 0 {
			numField += keyCountMap[t.Name()][0]
		}
	} else {
		if len(conds) == 0 {
			numField = keyCountMap[t.Name()][1]
		} else {
			numField = keyCountMap[t.Name()][0]
		}
	}

	col := 1
	rec, err := NewRecord(numField)
	if err != nil {
		return err
	}
	defer rec.Close()

	if len(fields) > 0 {
		for _, f := range fields {
			sf, ok := t.FieldByName(f)
			if !ok {
				return fmt.Errorf("no field: %v", f)
			}

			placeholders = append(placeholders, fmt.Sprintf("`%v`=?", fieldMap[t.Name()][f].name))

			if err := rec.setRecord(col, sf, v.FieldByName(f)); err != nil {
				return err
			}
			col++
		}
	} else {
		for i := 0; i < t.NumField(); i++ {
			sf := t.Field(i)
			info := fieldMap[t.Name()][sf.Name]

			if !info.isField {
				continue
			}
			if info.isKey {
				continue
			}
			placeholders = append(placeholders, fmt.Sprintf("`%v`=?", info.name))
			if err := rec.setRecord(col, sf, v.Field(i)); err != nil {
				return err
			}
			col++
		}
	}

	criterias := []string{}
	if len(conds) > 0 {
		criterias = buildCondition(conds)
	} else {
		for i := 0; i < t.NumField(); i++ {
			sf := t.Field(i)
			info := fieldMap[t.Name()][sf.Name]

			if !info.isField {
				continue
			}
			if !info.isKey {
				continue
			}
			criterias = append(criterias, fmt.Sprintf("`%v`=?", info.name))
			if err := rec.setRecord(col, sf, v.Field(i)); err != nil {
				return err
			}
			col++
		}
	}

	sql := fmt.Sprintf("UPDATE `%v` SET %v WHERE %v",
		tableMap[t.Name()],
		strings.Join(placeholders, ","),
		strings.Join(criterias, " AND "),
	)
	view, err := db.OpenView(sql)
	if err != nil {
		return err
	}
	defer view.Close()

	view.Execute(rec)
	return nil
}

func (db Database) Delete(data interface{}, conds map[string]interface{}) error {
	var rec Record
	t := reflect.TypeOf(data)
	v := reflect.ValueOf(data)

	col := 1
	criterias := []string{}
	if len(conds) > 0 {
		criterias = buildCondition(conds)
	} else {
		var err error
		rec, err = NewRecord(keyCountMap[t.Name()][0])
		if err != nil {
			return err
		}
		defer rec.Close()

		for i := 0; i < t.NumField(); i++ {
			sf := t.Field(i)
			info := fieldMap[t.Name()][sf.Name]

			if !info.isField {
				continue
			}
			if !info.isKey {
				continue
			}
			criterias = append(criterias, fmt.Sprintf("`%v`=?", info.name))
			if err := rec.setRecord(col, sf, v.Field(i)); err != nil {
				return err
			}
			col++
		}
	}

	criteria := ""
	if len(criterias) > 0 {
		criteria = "WHERE " + strings.Join(criterias, " AND ")
	}

	sql := fmt.Sprintf("DELETE FROM `%v` %v",
		tableMap[reflect.TypeOf(data).Name()],
		criteria,
	)
	view, err := db.OpenView(sql)
	if err != nil {
		return err
	}
	defer view.Close()

	view.Execute(rec)
	return nil
}

func (db Database) OpenView(query string) (View, error) {
	handle := uintptr(0)

	ret, _, err := procMsiDatabaseOpenView.Call(
		db.handle,
		uintptr(unsafe.Pointer(&[]byte(query)[0])),
		uintptr(unsafe.Pointer(&handle)),
	)
	if ret != 0 {
		return View{}, err
	}
	return View{handle}, nil
}

func (db Database) Import(folder, file string) error {
	ret, _, err := procMsiDatabaseImport.Call(
		db.handle,
		uintptr(unsafe.Pointer(&[]byte(folder)[0])),
		uintptr(unsafe.Pointer(&[]byte(file)[0])),
	)
	if ret != 0 {
		return err
	}
	return nil
}

func (db Database) CreateTransform(other Database, output string) error {
	ret, _, err := procMsiDatabaseGenerateTransform.Call(
		db.handle,
		other.handle,
		uintptr(unsafe.Pointer(&[]byte(output)[0])),
	)
	if ret != 0 {
		return err
	}

	ret, _, err = procMsiCreateTransformSummaryInfo.Call(
		db.handle,
		other.handle,
		uintptr(unsafe.Pointer(&[]byte(output)[0])),
		uintptr(0),
		uintptr(0),
	)
	if ret == 0 {
		return err
	}

	return nil
}

func buildCondition(conds map[string]interface{}) []string {
	criterias := []string{}
	for k, i := range conds {
		t := reflect.TypeOf(i)
		v := reflect.ValueOf(i)
		switch t.Kind() {
		case reflect.String:
			criterias = append(criterias, fmt.Sprintf("`%v` = '%v'", k, v.Interface()))
		default:
			criterias = append(criterias, fmt.Sprintf("`%v` = %v", k, v.Interface()))
		}
	}
	return criterias
}
