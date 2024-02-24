package utils

import (
	"context"

	"irptools/utils/fp"
	"irptools/utils/logs"
)

func DecorateExecution(ctx context.Context, scope string, exec func(ctx context.Context) error) error {
	l := logs.L(ctx)
	l.I(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")

	ctx, scopedl := logs.WithScope(ctx, scope)
	dur, err := fp.MeasureCallE(func() error {
		return exec(ctx)
	})

	if err != nil {
		scopedl.E("FAILED: exec time: %v; error: %v", dur, err)
	} else {
		scopedl.I("OK: exec time: %v", dur)
	}

	l.I("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")

	return err
}
