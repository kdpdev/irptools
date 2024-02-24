package misc

import (
	"encoding/json"
	"os"

	"irptools/utils/errs"
)

func LoadJsonConfigFromFile[T errs.Validatable](filePath string, setup func(cfg T) (T, error)) (T, error) {
	var cfg T

	cfgData, err := os.ReadFile(filePath)
	if err != nil {
		return cfg, errs.Wrap(err)
	}

	err = json.Unmarshal(cfgData, &cfg)
	if err != nil {
		return cfg, errs.Wrap(err)
	}

	err = errs.CheckValid(cfg, "config")
	if err != nil {
		return cfg, errs.Wrap(err)
	}

	cfg, err = setup(cfg)
	if err != nil {
		return cfg, errs.Wrap(err)
	}

	return cfg, errs.CheckValid(cfg, "config")
}
