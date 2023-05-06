package helper

import (
	"fmt"

	model "github.com/codeindex2937/msi-helper/model"
)

func (h Helper) AddService(id, name string, comp model.Component) error {
	if err := h.DB.Insert(model.ServiceControl{
		ServiceControl: fmt.Sprintf("INSTALL_SERVICE_%v", id),
		Name:           id,
		Event:          model.SERVICECTRL_EVENT_START | model.SERVICECTRL_EVENT_STOP | model.SERVICECTRL_EVENT_UNINSTALL_STOP,
		Wait:           model.OptInt16(1),
		Component:      comp.Component,
	}); err != nil {
		return err
	}

	if err := h.DB.Insert(model.ServiceControl{
		ServiceControl: fmt.Sprintf("REMOVE_OLD_SERVICE_%v", id),
		Name:           id,
		Event:          model.SERVICECTRL_EVENT_STOP | model.SERVICECTRL_EVENT_DELETE | model.SERVICECTRL_EVENT_UNINSTALL_STOP | model.SERVICECTRL_EVENT_UNINSTALL_DELETE,
		Wait:           model.OptInt16(1),
		Component:      comp.Component,
	}); err != nil {
		return err
	}

	if err := h.DB.Insert(model.ServiceInstall{
		ServiceInstall: fmt.Sprintf("SERVICE_%v", id),
		Name:           id,
		DisplayName:    model.OptString(name),
		ServiceType:    model.SERVICE_TYPE_OWN_PROCESS,
		StartType:      model.SERVICE_START_TYPE_AUTO,
		ErrorControl:   model.SERVICE_ERROR_CONTROL_NORMAL,
		Component:      comp.Component,
	}); err != nil {
		return err
	}
	return nil
}
