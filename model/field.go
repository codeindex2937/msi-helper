package msi

import (
	"fmt"
	"reflect"
	"strings"

	"golang.org/x/exp/slices"
)

var tableMap = map[string]string{}
var fieldsMap = map[string]string{}
var placeholdersMap = map[string]string{}
var keyCountMap = map[string][2]int{}
var fieldMap = map[string]map[string]fieldInfo{}
var nullStringType reflect.Type
var nullInt16Type reflect.Type
var nullInt32Type reflect.Type
var streamType reflect.Type

type Stream string

type NullString struct {
	String string
	Valid  bool
}

type NullInt16 struct {
	Int16 int16
	Valid bool
}

type NullInt32 struct {
	Int32 int32
	Valid bool
}

type fieldInfo struct {
	name    string
	isField bool
	isKey   bool
}

func init() {
	nullStringType = reflect.TypeOf(NullString{})
	nullInt16Type = reflect.TypeOf(NullInt16{})
	nullInt32Type = reflect.TypeOf(NullInt32{})
	streamType = reflect.TypeOf(Stream(""))

	updateKind(ActionText{}, "ActionText")
	updateKind(AdminExecuteSequence{}, "AdminExecuteSequence")
	updateKind(AdminUISequence{}, "AdminUISequence")
	updateKind(AdvtExecuteSequence{}, "AdvtExecuteSequence")
	updateKind(AppSearch{}, "AppSearch")
	updateKind(BBControl{}, "BBControl")
	updateKind(Billboard{}, "Billboard")
	updateKind(Binary{}, "Binary")
	updateKind(CheckBox{}, "CheckBox")
	updateKind(ComboBox{}, "ComboBox")
	updateKind(Component{}, "Component")
	updateKind(Control{}, "Control")
	updateKind(ControlCondition{}, "ControlCondition")
	updateKind(ControlEvent{}, "ControlEvent")
	updateKind(CreateFolder{}, "CreateFolder")
	updateKind(CustomAction{}, "CustomAction")
	updateKind(Dialog{}, "Dialog")
	updateKind(Directory{}, "Directory")
	updateKind(DrLocator{}, "DrLocator")
	updateKind(Error{}, "Error")
	updateKind(EventMapping{}, "EventMapping")
	updateKind(Feature{}, "Feature")
	updateKind(FeatureComponents{}, "FeatureComponents")
	updateKind(File{}, "File")
	updateKind(Icon{}, "Icon")
	updateKind(InstallExecuteSequence{}, "InstallExecuteSequence")
	updateKind(InstallUISequence{}, "InstallUISequence")
	updateKind(LaunchCondition{}, "LaunchCondition")
	updateKind(ListBox{}, "ListBox")
	updateKind(ListView{}, "ListView")
	updateKind(Media{}, "Media")
	updateKind(Property{}, "Property")
	updateKind(RadioButton{}, "RadioButton")
	updateKind(RegLocator{}, "RegLocator")
	updateKind(Registry{}, "Registry")
	updateKind(RemoveFile{}, "RemoveFile")
	updateKind(RemoveRegistry{}, "RemoveRegistry")
	updateKind(ServiceControl{}, "ServiceControl")
	updateKind(ServiceInstall{}, "ServiceInstall")
	updateKind(Shortcut{}, "Shortcut")
	updateKind(Signature{}, "Signature")
	updateKind(TextStyle{}, "TextStyle")
	updateKind(UIText{}, "UIText")
	updateKind(Upgrade{}, "Upgrade")
	updateKind(Validation{}, "_Validation")
	updateKind(Streams{}, "_Streams")
}

func OptString(v string) NullString {
	return NullString{String: v, Valid: true}
}

func OptInt16(v int16) NullInt16 {
	return NullInt16{Int16: v, Valid: true}
}

func OptInt32(v int32) NullInt32 {
	return NullInt32{Int32: v, Valid: true}
}

func updateKind(inst interface{}, table string) {
	t := reflect.TypeOf(inst)

	fields := []string{}
	placeholders := []string{}
	keyCount := 0
	for i := 0; i < t.NumField(); i++ {
		field, isField, isKey := getFieldTag(t.Field(i))
		if !isField {
			continue
		}
		if isKey {
			keyCount++
		}

		if _, ok := fieldMap[t.Name()]; !ok {
			fieldMap[t.Name()] = map[string]fieldInfo{}
		}

		fieldMap[t.Name()][t.Field(i).Name] = fieldInfo{field, isField, isKey}
		fields = append(fields, fmt.Sprintf("`%v`", field))
		placeholders = append(placeholders, "?")
	}

	tableMap[t.Name()] = table
	keyCountMap[t.Name()] = [2]int{keyCount, len(fields)}
	fieldsMap[t.Name()] = strings.Join(fields, ",")
	placeholdersMap[t.Name()] = strings.Join(placeholders, ",")
}

func getFieldTag(f reflect.StructField) (string, bool, bool) {
	parts := strings.Split(f.Tag.Get("msi"), ",")
	var isKey, isField bool
	if len(parts[0]) > 0 {
		isField = true
	}
	if len(parts) > 1 && slices.Contains(parts[1:], "key") {
		isKey = true
	}
	return parts[0], isField, isKey
}
