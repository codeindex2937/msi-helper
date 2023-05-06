package helper

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	model "github.com/codeindex2937/msi-helper/model"
)

func NewComponent(feature, id, name, dir string) model.Component {
	return model.Component{
		Feature:     feature,
		Component:   name,
		ComponentId: model.OptString(id),
		Directory:   dir,
		Files:       map[string]model.FileCredential{},
	}
}

func (h *Helper) AddComponent(comp *model.Component, srcdir, keyfile string) error {
	if len(srcdir) > 0 {
		if err := h.addFiles(comp, srcdir); err != nil {
			return err
		}
	}

	for _, f := range comp.Files {
		stat, err := os.Stat(f.Path)
		if err != nil {
			return err
		}

		if err := h.DB.Insert(model.File{
			File:       f.Key,
			Component:  comp.Component,
			FileName:   f.ShortName,
			FileSize:   int32(stat.Size()),
			Attributes: model.OptInt16(model.FILE_ATTRIBUTES_COMPRESSED),
			Sequence:   int16(h.fileSequence),
		}); err != nil {
			return err
		}

		h.fileSequence++
	}

	if len(keyfile) > 0 {
		if !comp.SetKeyPath(keyfile) {
			return fmt.Errorf("invalid key path")
		}
	}

	if err := h.DB.Insert(*comp); err != nil {
		return err
	}
	if err := h.DB.Insert(model.FeatureComponents{
		Feature:   comp.Feature,
		Component: comp.Component,
	}); err != nil {
		return err
	}

	return nil
}

func (h *Helper) addFiles(c *model.Component, path string) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		re := regexp.MustCompile(`[^\w.]`)
		extFixed := strings.ReplaceAll(e.Name(), ".json", ".jso")
		partialKey := re.ReplaceAllString(extFixed, "_")
		c.Files[e.Name()] = model.FileCredential{
			Key:       c.Component + partialKey,
			Path:      filepath.Join(path, e.Name()),
			ShortName: h.getMSIFileName(e.Name()),
		}
	}
	return nil
}

func (h *Helper) getMSIFileName(name string) string {
	initName := getShortfileName(name, 0)
	count, ok := h.fileNameConflict[initName]
	if !ok {
		h.fileNameConflict[initName] = 0
		return fmt.Sprintf("%v|%v", initName, name)
	} else {
		return fmt.Sprintf("%v|%v", getShortfileName(name, count), name)
	}
}

func getShortfileName(name string, count int) string {
	ext := filepath.Ext(name)

	pureName := name[:len(name)-len(ext)]
	if len(pureName) > 8 {
		pureName = fmt.Sprintf("%v~%v", pureName[:6], count+1)
	}

	if ext == ".json" {
		ext = ".jso"
	}

	return pureName + ext
}
