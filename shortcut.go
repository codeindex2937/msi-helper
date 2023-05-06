package helper

import (
	"fmt"

	model "github.com/codeindex2937/msi-helper/model"
)

func (h Helper) AddDesktopShortcut(
	name,
	guid string,
	comp model.Component,
	icon *model.Icon,
) error {
	signature := "DESKTOP_SHORTCUT"
	shortcutComponent := "DESKTOP_SHORTCUT"
	directory := "DesktopFolder"

	shortcut := model.Shortcut{
		Target:    comp.Feature,
		Component: comp.Component,
		WkDir:     model.OptString(comp.Directory),
		Shortcut:  "desktop_shortcut",
		Directory: directory,
		Name:      name,
		Arguments: model.OptString(""),
		IconIndex: model.OptInt16(0),
		ShowCmd:   model.OptInt16(1),
	}
	if icon != nil {
		shortcut.Icon = model.OptString(icon.Name)
	}
	if err := h.DB.Insert(shortcut); err != nil {
		return err
	}

	if err := h.DB.Insert(model.Signature{
		Signature: signature,
		FileName:  fmt.Sprintf("%v.lnk", name),
	}); err != nil {
		return err
	}
	if err := h.DB.Insert(model.DrLocator{
		Signature: signature,
		Path:      model.OptString(fmt.Sprintf("[%v]", directory)),
	}); err != nil {
		return err
	}
	if err := h.DB.Insert(model.AppSearch{
		Property:  "DESKTOP_SHORTCUT_EXIST",
		Signature: signature,
	}); err != nil {
		return err
	}

	if err := h.DB.Insert(model.Component{
		Component:   shortcutComponent,
		ComponentId: model.OptString(guid),
		Directory:   directory,
		Attributes:  int16(model.COMPONENT_ATTRIBUTES_64BIT),
	}); err != nil {
		return err
	}
	if err := h.DB.Insert(model.FeatureComponents{
		Feature:   comp.Feature,
		Component: shortcutComponent,
	}); err != nil {
		return err
	}

	if icon != nil {
		if err := h.DB.Update(model.Property{Property: "ARPPRODUCTICON", Value: icon.Name}, nil, nil); err != nil {
			return err
		}
	}

	return nil
}

func (h Helper) AddMenuShortcut(
	linkName,
	folderName,
	guid string,
	comp model.Component,
	icon *model.Icon,
) error {
	signature := "STARTMENU_SHORTCUT"
	shortcutComponent := "MENU_SHORTCUT"
	directory := "StartMenu"

	shortcut := model.Shortcut{
		Shortcut:  "menu_shortcut",
		Target:    comp.Feature,
		Component: comp.Component,
		WkDir:     model.OptString(comp.Directory),
		Directory: directory,
		Name:      linkName,
		Arguments: model.OptString(""),
		IconIndex: model.OptInt16(0),
		ShowCmd:   model.OptInt16(1),
	}
	if icon != nil {
		shortcut.Icon = model.OptString(icon.Name)
	}
	if err := h.DB.Insert(shortcut); err != nil {
		return err
	}

	if err := h.DB.Insert(model.Signature{
		Signature: signature,
		FileName:  fmt.Sprintf("%v.lnk", linkName),
	}); err != nil {
		return err
	}
	if err := h.DB.Insert(model.DrLocator{
		Signature: signature,
		Path:      model.OptString(`[ProgramMenuFolder]\Example`),
	}); err != nil {
		return err
	}
	if err := h.DB.Insert(model.AppSearch{
		Signature: signature,
		Property:  "STARTMENU_SHORTCUT_EXIST",
	}); err != nil {
		return err
	}

	if err := h.DB.Insert(model.Component{
		Component:   shortcutComponent,
		ComponentId: model.OptString(guid),
		Directory:   directory,
		Attributes:  int16(model.COMPONENT_ATTRIBUTES_64BIT),
	}); err != nil {
		return err
	}
	if err := h.DB.Insert(model.FeatureComponents{
		Feature:   comp.Feature,
		Component: shortcutComponent,
	}); err != nil {
		return err
	}

	programMenuFolder := "ProgramMenuFolder"
	if err := h.DB.Insert(model.Directory{
		Directory:       programMenuFolder,
		DirectoryParent: model.OptString("TARGETDIR"),
		DefaultDir:      ".",
	}); err != nil {
		return err
	}
	if err := h.DB.Insert(model.Directory{
		Directory:       directory,
		DirectoryParent: model.OptString(programMenuFolder),
		DefaultDir:      h.getMSIFileName(folderName),
	}); err != nil {
		return err
	}
	if err := h.DB.Insert(model.CreateFolder{
		Component: shortcutComponent,
		Directory: directory,
	}); err != nil {
		return err
	}
	return nil
}
