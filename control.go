package helper

import model "github.com/codeindex2937/msi-helper/model"

func (h Helper) AddConditionalControl(control model.Control, conditions ...string) error {
	if err := h.DB.Insert(control); err != nil {
		return err
	}

	for _, cond := range conditions {
		if err := h.DB.Insert(model.ControlCondition{
			Dialog:    control.Dialog,
			Control:   control.Control,
			Action:    "Hide",
			Condition: cond,
		}); err != nil {
			return err
		}
	}
	return nil
}
