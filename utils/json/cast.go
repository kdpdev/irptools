package jsonutils

import (
	"encoding/json"

	"irptools/utils/errs"
)

func Cast(from interface{}, to interface{}) error {
	jsonData, err := json.Marshal(from)
	if err != nil {
		return errs.Wrap(err)
	}
	err = json.Unmarshal(jsonData, to)
	if err != nil {
		return errs.Wrap(err)
	}
	return nil
}
