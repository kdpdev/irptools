package parse

import (
	"fmt"
	"path/filepath"

	"irptools/tools/utils"
	"irptools/utils/errs"
	"irptools/utils/misc"
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func LoadConfig(filePath string) (Config, error) {
	return misc.LoadJsonConfigFromFile(filePath, func(cfg Config) (Config, error) {
		return cfg.Adjust()
	})
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type Config struct {
	Target  TargetConfig            `json:"target"`
	Sources map[string]SourceConfig `json:"sources"`
}

func (this Config) Validate() error {
	return errs.Catch(func() {
		errs.ThrowCheckValid(this.Target, "target")
		errs.ThrowCheckNotZero(len(this.Sources), "len(sources)")
		for key, src := range this.Sources {
			errs.ThrowCheckValid(src, fmt.Sprintf("source[%s]", key))
			errs.ThrowIf(this.Target.Folder.ValidateSourcePath(src.Path))
		}
	})
}

func (this Config) Adjust() (Config, error) {
	var err error

	this.Target, err = this.Target.Adjust()
	if err != nil {
		return this, errs.Wrap(err)
	}

	for i, src := range this.Sources {
		this.Sources[i], err = src.Adjust()
		if err != nil {
			return this, errs.Wrap(err)
		}
	}

	return this, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type SourceOptions = interface{}

type SourceConfig struct {
	Skip    bool          `json:"skip"`
	Type    string        `json:"type"`
	Path    string        `json:"path"`
	Options SourceOptions `json:"options"`
}

func (this SourceConfig) Validate() error {
	return errs.Catch(func() {
		errs.ThrowCheckRequiredString(this.Type, "type")
		errs.ThrowCheckRequiredString(this.Path, "path")
	})
}

func (this SourceConfig) Adjust() (SourceConfig, error) {
	var err error
	this.Path, err = filepath.Abs(this.Path)
	if err != nil {
		return this, errs.Wrap(err)
	}
	return this, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type TargetConfig struct {
	Folder          utils.TargetFolder `json:"folder"`
	WithStat        bool               `json:"withStat"`
	PrettyJsonPrint bool               `json:"prettyJsonPrint"`
	KeepSourceField bool               `json:"keepSourceField"`
	FieldsToLower   bool               `json:"fieldsToLower"`
}

func (this TargetConfig) Validate() error {
	return errs.Catch(func() {
		errs.ThrowCheckValid(this.Folder, "folder")
	})
}

func (this TargetConfig) Adjust() (TargetConfig, error) {
	var err error
	this.Folder, err = this.Folder.Adjust()
	return this, errs.Wrap(err)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
