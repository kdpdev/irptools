package errs

import (
	"errors"
	"runtime"
	"strings"
)

func NewPackageError(text string) error {
	if text == "" {
		text = "error"
	}
	var b strings.Builder
	b.WriteString(getPackageName())
	b.WriteString(": ")
	b.WriteString(text)
	text = b.String()
	return errors.New(text)
}

func getPackageName() string {
	pc, _, _, _ := runtime.Caller(2)
	parts := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	pl := len(parts)
	packageName := ""

	if parts[pl-2][0] == '(' {
		packageName = strings.Join(parts[0:pl-2], ".")
	} else {
		packageName = strings.Join(parts[0:pl-1], ".")
	}

	return packageName
}
