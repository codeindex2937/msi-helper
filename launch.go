package helper

import (
	"fmt"

	model "github.com/codeindex2937/msi-helper/model"
)

func (h Helper) LaunchAppPostInstall(script Script, comp model.Component, exec string) error {
	hideCondition := []string{"Installed", "VersionHandler < \"5.00\""}

	if err := h.DB.Insert(script.Binary); err != nil {
		return err
	}
	if err := h.DB.Insert(model.CheckBox{
		Property: "LaunchApp",
		Value:    model.OptString("1"),
	}); err != nil {
		return err
	}

	if err := h.DB.Insert(model.ControlEvent{
		Dialog:    "ExitDialog",
		Control:   "Finish",
		Event:     "DoAction",
		Argument:  "LaunchApp",
		Condition: model.OptString("LaunchApp = 1 AND NOT Installed"),
		Ordering:  model.OptInt16(1),
	}); err != nil {
		return err
	}

	if err := h.DB.Insert(model.Property{
		Property: "LaunchApp",
		Value:    "1",
	}); err != nil {
		return err
	}

	if err := h.AddConditionalControl(model.Control{
		Dialog:      "ExitDialog",
		Control:     "LaunchAppBox",
		Type:        "CheckBox",
		X:           135,
		Y:           125,
		Width:       11,
		Height:      11,
		Attributes:  model.OptInt32(model.CONTROL_ATTRIBUTES_ENABLED | model.CONTROL_ATTRIBUTES_VISIBLE),
		Property:    model.OptString("LaunchApp"),
		ControlNext: model.OptString("Back"),
	}, hideCondition...); err != nil {
		return err
	}

	if err := h.DB.Update(model.Control{
		Dialog:      "ExitDialog",
		Control:     "Bitmap",
		ControlNext: model.OptString("LaunchAppBox"),
	}, []string{"ControlNext"}, nil); err != nil {
		return err
	}

	if err := h.AddConditionalControl(model.Control{
		Dialog:     "ExitDialog",
		Control:    "LaunchAppBoxText",
		Type:       "Text",
		X:          150,
		Y:          124,
		Width:      200,
		Height:     20,
		Attributes: model.OptInt32(model.CONTROL_ATTRIBUTES_ENABLED | model.CONTROL_ATTRIBUTES_VISIBLE | model.CONTROL_ATTRIBUTES_TRANSPARENT | model.CONTROL_ATTRIBUTES_NO_PREFIX),
		Text:       model.OptString("Launch application when finished."),
	}, hideCondition...); err != nil {
		return err
	}

	name := "LaunchApp"
	if err := h.DB.Insert(model.CustomAction{
		Action: name,
		Type:   int16(script.Type | model.CUSTOM_ACTION_TYPE_CONTINUE | model.CUSTOM_ACTION_TYPE_ASYNC),
		Source: model.OptString(script.Name),
		Target: model.OptString("LaunchApp"),
	}); err != nil {
		return err
	}
	if err := h.DB.Insert(model.InstallExecuteSequence{
		Action: name,
	}); err != nil {
		return err
	}

	initAction := "INIT_" + name
	if err := h.DB.Insert(model.CustomAction{
		Action: initAction,
		Type:   model.CUSTOM_ACTION_TYPE_SET_FORMATTED_PROPERTY,
		Source: model.OptString("LaunchAppPath"),
		Target: model.OptString(fmt.Sprintf("[#%v]", comp.GetFileKey(exec))),
	}); err != nil {
		return err
	}
	if err := h.DB.Insert(model.InstallUISequence{
		Action:   initAction,
		Sequence: model.OptInt16(1001),
	}); err != nil {
		return err
	}
	return nil
}
