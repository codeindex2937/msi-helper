package helper

import model "github.com/codeindex2937/msi-helper/model"

type Helper struct {
	DB               model.Database
	fileSequence     int
	fileNameConflict map[string]int
	upgradeCode      string
}

func New(db model.Database, upgradeCode string) Helper {
	return Helper{
		DB:               db,
		fileSequence:     1,
		fileNameConflict: map[string]int{},
		upgradeCode:      upgradeCode,
	}
}

func SerializeDirectories(dir *model.Directory) []*model.Directory {
	dirs := []*model.Directory{dir}
	for _, subdir := range dir.Dirs {
		subdir.DirectoryParent = model.OptString(dir.Directory)
		dirs = append(dirs, SerializeDirectories(subdir)...)
	}
	return dirs
}
