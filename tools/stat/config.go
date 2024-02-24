package stat

import (
	"path/filepath"

	"irptools/tools/utils"
	"irptools/utils/errs"
	"irptools/utils/misc"
)

func LoadConfig(filePath string) (Config, error) {
	return misc.LoadJsonConfigFromFile(filePath, func(cfg Config) (Config, error) {
		return cfg.Adjust()
	})
}

type Config struct {
	Source string             `json:"source"`
	Target utils.TargetFolder `json:"target"`
}

func (this Config) Validate() error {
	return errs.Catch(func() {
		errs.ThrowCheckValid(this.Target, "target")
		errs.ThrowCheckRequiredString(this.Source, "source")
		errs.ThrowIf(this.Target.ValidateSourcePath(this.Source))
	})
}

func (this Config) Adjust() (Config, error) {
	var err error

	this.Target, err = this.Target.Adjust()
	if err != nil {
		return this, errs.Wrap(err)
	}

	this.Source, err = filepath.Abs(this.Source)
	if err != nil {
		return this, errs.Wrap(err)
	}

	return this, nil
}
