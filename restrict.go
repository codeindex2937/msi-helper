package helper

import model "github.com/codeindex2937/msi-helper/model"

func (h Helper) BlockOldVersion(productVersion string) error {
	if err := h.AddAction(Action{
		Name:      "EXIT_OLD_INSTALL",
		Sequence:  201,
		Defered:   false,
		Condition: model.OptString("NEWER_INSTALLED"),
	}, "", "A higher upgrade is already installed.",
		model.CUSTOM_ACTION_TYPE_ALERT_AND_EXIT,
	); err != nil {
		return err
	}

	if err := h.DB.Update(model.Control{
		Dialog:  "PrepareDlg",
		Control: "Description",
		Height:  60,
	}, []string{"Height"}, nil); err != nil {
		return err
	}

	if err := h.DB.Insert(model.Upgrade{
		UpgradeCode:    h.upgradeCode,
		VersionMin:     model.OptString(productVersion),
		Attributes:     model.UPGRADE_ATTRIBUTES_VERSION_ONLY_DETECT,
		ActionProperty: "NEWER_INSTALLED",
	}); err != nil {
		return err
	}

	if err := h.DB.Insert(model.Upgrade{
		UpgradeCode:    h.upgradeCode,
		VersionMax:     model.OptString(productVersion),
		Attributes:     model.UPGRADE_ATTRIBUTES_VERSION_MAX_INCLUSIVE,
		ActionProperty: "UPGRADE_FOUND",
	}); err != nil {
		return err
	}

	return nil
}
