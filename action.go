package helper

import model "github.com/codeindex2937/msi-helper/model"

type Action struct {
	Name      string
	Method    string
	Sequence  int16
	Defered   bool
	Condition model.NullString
	Parameter string
}

type Script struct {
	model.Binary
	Type int
}

func (h Helper) AddScript(script Script, actions []Action) error {
	if err := h.DB.Insert(script.Binary); err != nil {
		return err
	}
	for _, action := range actions {
		action.Name = script.Name + action.Name
		if err := h.AddAction(
			action,
			script.Name,
			action.Method,
			script.Type,
		); err != nil {
			return err
		}
	}
	return nil
}

func (h Helper) AddAction(action Action, source, target string, attr int) error {
	if action.Defered {
		attr |= model.CUSTOM_ACTION_TYPE_IN_SCRIPT | model.CUSTOM_ACTION_TYPE_NO_IMPERSONATE | model.CUSTOM_ACTION_TYPE_HIDE_TARGET
	} else {
		attr |= model.CUSTOM_ACTION_TYPE_FIRST_SEQUENCE
	}

	if err := h.DB.Insert(model.CustomAction{
		Action: action.Name,
		Type:   int16(attr),
		Source: model.OptString(source),
		Target: model.OptString(target),
	}); err != nil {
		return err
	}

	if err := h.DB.Insert(model.InstallExecuteSequence{
		Action:    action.Name,
		Condition: action.Condition,
		Sequence:  model.OptInt16(action.Sequence),
	}); err != nil {
		return err
	}
	if !action.Defered {
		if err := h.DB.Insert(model.InstallUISequence{
			Action:    action.Name,
			Condition: action.Condition,
			Sequence:  model.OptInt16(action.Sequence),
		}); err != nil {
			return err
		}
	}
	if len(action.Parameter) > 0 {
		initAction := "INIT_" + action.Name
		if err := h.DB.Insert(model.CustomAction{
			Action: initAction,
			Type:   model.CUSTOM_ACTION_TYPE_SET_FORMATTED_PROPERTY,
			Source: model.OptString(action.Name),
			Target: model.OptString(action.Parameter),
		}); err != nil {
			return err
		}
		if err := h.DB.Insert(model.InstallExecuteSequence{
			Action:    initAction,
			Condition: action.Condition,
			Sequence:  model.OptInt16(action.Sequence - 1),
		}); err != nil {
			return err
		}
	}
	return nil
}
