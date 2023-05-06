package helper

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	model "github.com/codeindex2937/msi-helper/model"
)

func (h Helper) PackFiles(target, wdir string, comps ...model.Component) error {
	ddfName := "cabinet.ddf"
	ddfPath := filepath.Join(wdir, ddfName)

	if err := os.MkdirAll(wdir, 0755); err != nil {
		return err
	}

	f, err := os.Create(ddfPath)
	if err != nil {
		return err
	}

	_, _ = f.WriteString(".Set Cabinet=on\n")
	_, _ = f.WriteString(".Set Compress=on\n")
	_, _ = f.WriteString(".Set UniqueFiles=on\n")
	_, _ = f.WriteString(".Set RptFileName=nul\n")
	_, _ = f.WriteString(".Set InfFileName=nul\n")
	_, _ = f.WriteString(".Set CompressionType=MSZIP\n")
	_, _ = f.WriteString(".Set MaxDiskSize=0\n")               // do not split into multiple cabinets
	_, _ = f.WriteString(".Set DiskDirectoryTemplate=Disk*\n") // all cabinets go into a single directory
	_, _ = f.WriteString(fmt.Sprintf(".Set CabinetNameTemplate=%v\n", target))
	for _, comp := range comps {
		for _, cred := range comp.Files {
			abs, err := filepath.Abs(cred.Path)
			if err != nil {
				return err
			}
			_, _ = f.WriteString(fmt.Sprintf("%v %v\n", abs, cred.Key))
		}
	}
	f.Close()

	cmd := exec.Command("makecab", "/F", ddfName)
	cmd.Dir = wdir
	if err := cmd.Run(); err != nil {
		return nil
	}

	if err := h.DB.Insert(model.Streams{
		Name: target,
		Data: model.Stream(filepath.Join(wdir, "Disk1", target)),
	}); err != nil {
		return err
	}
	if err := h.DB.Insert(model.Media{
		DiskId:       1,
		LastSequence: int16(h.fileSequence - 1),
		Cabinet:      model.OptString("#" + target),
	}); err != nil {
		return err
	}

	return nil
}
