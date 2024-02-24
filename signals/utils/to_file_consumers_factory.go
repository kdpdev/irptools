package utils

import (
	"path/filepath"
	"strings"

	"irptools/utils/errs"
)

type TargetFilePathStrategyFn func(filePath string) (string, error)

func NewSignalsToFileConsumersFactory(getConsumer SignalsToFileConsumerSourceFn, getTargetFilePath TargetFilePathStrategyFn) *SignalsToFileConsumersFactory {
	return &SignalsToFileConsumersFactory{
		getConsumer:       getConsumer,
		getTargetFilePath: getTargetFilePath,
	}
}

type SignalsToFileConsumersFactory struct {
	getConsumer       SignalsToFileConsumerSourceFn
	getTargetFilePath TargetFilePathStrategyFn
}

func (this *SignalsToFileConsumersFactory) NewConsumer(filePath string) (ClosableSignalConsumer, error) {
	targetFilePath, err := this.getTargetFilePath(filePath)
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return this.getConsumer(targetFilePath)
}

func RepeatSourceTreeTargetFilePathStrategy(sourcePath string, targetPath string) TargetFilePathStrategyFn {
	return func(filePath string) (string, error) {
		filePath, err := filepath.Abs(filePath)
		if err != nil {
			return "", errs.Wrap(err)
		}

		if strings.Index(filePath, sourcePath) != 0 {
			return "", errs.Errorf("unexpected file path: '%s'", filePath)
		}

		relFilePath := filePath[len(sourcePath):]
		if relFilePath == "" {
			_, relFilePath = filepath.Split(filePath)
		}
		resultsFilePath := filepath.Join(targetPath, relFilePath)
		return resultsFilePath, nil
	}
}

func ToOneFolderTargetFilePathStrategy(sourcePath string, targetPath string) TargetFilePathStrategyFn {
	return func(filePath string) (string, error) {
		filePath, err := filepath.Abs(filePath)
		if err != nil {
			return "", errs.Wrap(err)
		}

		if strings.Index(filePath, sourcePath) != 0 {
			return "", errs.Errorf("unexpected file path: '%s'", filePath)
		}

		relFilePath := filePath[len(sourcePath):]
		if relFilePath == "" {
			_, relFilePath = filepath.Split(filePath)
		} else {
			separatorBytes := [...]byte{filepath.Separator}
			separator := string(separatorBytes[:])
			relFilePath = strings.Trim(relFilePath, separator)
			parts := strings.Split(relFilePath, separator)
			relFilePath = strings.Join(parts, "__")
		}
		resultsFilePath := filepath.Join(targetPath, relFilePath)
		return resultsFilePath, nil
	}
}
