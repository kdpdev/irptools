package fs

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"irptools/utils/errs"
)

type FilesConsumer = func(filePath string) (next bool, err error)
type PathFilter func(filePath string) bool

func AdjustPathSlash(path string) string {
	return strings.ReplaceAll(path, "\\", "/")
}

func AdjustPathsSlash(paths []string) []string {
	result := make([]string, 0, len(paths))
	for _, p := range paths {
		result = append(result, AdjustPathSlash(p))
	}
	return paths
}

func EnumFilePaths(root string, consume FilesConsumer) error {
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		next, err := consume(AdjustPathSlash(path))
		if err != nil {
			return errs.Wrap(err)
		}
		if !next {
			return context.Canceled
		}
		return nil
	})

	return errs.Wrap(err)
}

func EnumFilePathsWithFilter(root string, pass PathFilter, consume FilesConsumer) error {
	return EnumFilePaths(root, func(filePath string) (bool, error) {
		if pass(filePath) {
			next, err := consume(filePath)
			if err != nil {
				return false, errs.Wrap(err)
			}
			return next, nil
		}
		return true, nil
	})
}

func EnumFilePathsWithExt(root string, ext string, consume FilesConsumer) error {
	isExt := func(filePath string) bool {
		fileExt := strings.ToLower(filepath.Ext(filePath))
		return fileExt == ext
	}
	return EnumFilePathsWithFilter(root, isExt, consume)
}

func OpenReadOnlyFile(filePath string) (*os.File, error) {
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func CreateWriteOnlyFile(filePath string) (*os.File, error) {
	return os.OpenFile(filePath, os.O_CREATE|os.O_EXCL|os.O_APPEND|os.O_WRONLY, 0644)
}

func GetFileSize(file *os.File) (uint64, error) {
	stat, err := file.Stat()
	if err != nil {
		return 0, err
	}
	return uint64(stat.Size()), nil
}

func EnsureDirExists(dirPath string) (created bool, err error) {
	fi, err := os.Stat(dirPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return false, errs.Errorf("os.Stat failed: %w", err)
		}

		err = os.MkdirAll(dirPath, 0644)
		if err != nil {
			return false, errs.Errorf("os.Mkdir failed: %w", err)
		}

		return true, nil
	}

	fileMode := fi.Mode()
	if !fileMode.IsDir() {
		return false, errs.Errorf("existing '%v' fs entry is not a dir", dirPath)
	}

	return false, nil
}

func EnsureDirCreated(dirPath string) error {
	created, err := EnsureDirExists(dirPath)
	if err != nil {
		return errs.Errorf("failed to ensure the '%v' dir exists: %w", dirPath, err)
	}

	if !created {
		return errs.Errorf("the '%v' is already existed", dirPath)
	}

	return nil
}

func RemoveGlob(pathPattern string) error {
	matches, err := filepath.Glob(pathPattern)
	if err != nil {
		return errs.Wrap(err)
	}

	for _, item := range matches {
		err = os.RemoveAll(item)
		if err != nil {
			return errs.Errorf("failed to remove item: %w", err)
		}
	}

	return nil
}
