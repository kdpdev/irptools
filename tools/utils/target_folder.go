package utils

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"irptools/utils/errs"
	"irptools/utils/fs"
)

type TargetFolder struct {
	Path                string `json:"path"`
	CleanupIfExists     bool   `json:"cleanupIfExists"`
	WithoutCreationTime bool   `json:"withoutCreationTime"`
}

func (this TargetFolder) Validate() error {
	return errs.CheckRequiredString(this.Path, "path")
}

func (this TargetFolder) ValidateSourcePath(sourcePath string) error {
	if strings.Index(this.Path, sourcePath) == 0 && len(this.Path) >= len(sourcePath) {
		return errs.Errorf("source contains target: source = '%s', target = '%s'", sourcePath, this.Path)
	}
	return nil
}

func (this TargetFolder) Adjust() (TargetFolder, error) {
	var err error
	this.Path, err = filepath.Abs(this.Path)
	if err != nil {
		return this, errs.Wrap(err)
	}
	return this, nil
}

func (this TargetFolder) Join(elems ...string) TargetFolder {
	this.Path = filepath.Join(append([]string{this.Path}, elems...)...)
	return this
}

func (this TargetFolder) PrepareTarget() (TargetFolder, error) {
	if !this.WithoutCreationTime {
		now := time.Now()
		nowStr := now.Format("2006-01-02_15-04-05") + fmt.Sprintf("_%03d", now.UnixMilli()%1000)
		this = this.Join(nowStr)
		this.WithoutCreationTime = false
	}

	created, err := fs.EnsureDirExists(this.Path)
	if err != nil {
		return this, errs.Errorf("failed to ensure target dir exists: %w", err)
	}

	if !created {
		if this.CleanupIfExists {
			pattern := filepath.Join(this.Path, "*")
			err = fs.RemoveGlob(pattern)
			if err != nil {
				return this, errs.Errorf("failed to cleanup target path: %w", err)
			}
			this.CleanupIfExists = true
		}
	}

	return this, nil
}
