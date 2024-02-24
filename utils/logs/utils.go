package logs

import (
	"path/filepath"
	"runtime"
)

func callerFuncName() string {
	pc, _, _, _ := runtime.Caller(2)
	fn := runtime.FuncForPC(pc)
	fnName := filepath.Base(fn.Name())
	return fnName
}
