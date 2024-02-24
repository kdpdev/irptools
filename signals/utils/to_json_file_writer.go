package utils

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"irptools/signals/signal"
	"irptools/utils/errs"
	"irptools/utils/fs"
)

func NewJsonFileWriter(filePath string) (*JsonFileWriter, error) {
	filePath, err := filepath.Abs(filePath)
	if err != nil {
		return nil, errs.Wrap(err)
	}

	dirPath, _ := filepath.Split(filePath)
	_, err = fs.EnsureDirExists(dirPath)
	if err != nil {
		return nil, errs.Wrap(err)
	}

	const jsonExt = ".json"
	if strings.LastIndex(filePath, jsonExt) != len(filePath)-len(jsonExt) {
		if filePath[len(filePath)-1] != '.' {
			filePath += "."
		}
		filePath += "json"
	}

	file, err := fs.CreateWriteOnlyFile(filePath)
	if err != nil {
		return nil, errs.Wrap(err)
	}

	return &JsonFileWriter{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

type JsonFileWriter struct {
	file    *os.File
	encoder *json.Encoder
}

func (this *JsonFileWriter) Consume(signal signal.Signal) error {
	return this.encoder.Encode(signal)
}

func (this *JsonFileWriter) Close() error {
	return this.file.Close()
}
