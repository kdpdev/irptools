package export_fz

import (
	"context"

	signalutils "irptools/signals/utils"
	"irptools/tools/utils"
	"irptools/utils/errs"
	"irptools/utils/logs"
)

func Main(ctx context.Context, cfg Config) error {
	return utils.DecorateExecution(ctx, "EXPORT FZ", func(ctx context.Context) error {
		return execMain(ctx, cfg)
	})
}

func execMain(ctx context.Context, cfg Config) error {
	err := errs.CheckValid(cfg, "config")
	if err != nil {
		return err
	}

	cfg.Target.Folder, err = cfg.Target.Folder.PrepareTarget()
	if err != nil {
		return errs.Errorf("failed to prepare target: %w", err)
	}

	l := logs.L(ctx)
	l.I("source <-: %s", cfg.Source)
	l.I("target ->: %s", cfg.Target.Folder.Path)

	execCfg := cfg
	execCfg.Target.Folder = execCfg.Target.Folder.Join("result")
	err = execExportFz(ctx, execCfg)
	if err != nil {
		return errs.Wrap(err)
	}

	return nil
}

func execExportFz(ctx context.Context, cfg Config) error {
	getConsumer := func(filePath string) (signalutils.ClosableSignalConsumer, error) {
		return signalutils.NewIrFileWriter(filePath)
	}

	getTargetFilePath := signalutils.RepeatSourceTreeTargetFilePathStrategy(cfg.Source, cfg.Target.Folder.Path)
	if cfg.Target.ToOneFolder {
		getTargetFilePath = signalutils.ToOneFolderTargetFilePathStrategy(cfg.Source, cfg.Target.Folder.Path)
	}

	factory := signalutils.NewSignalsToFileConsumersFactory(getConsumer, getTargetFilePath)

	err := signalutils.EnumSignals(ctx, cfg.Source, factory.NewConsumer)

	return errs.Wrap(err)
}
