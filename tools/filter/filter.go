package filter

import (
	"context"
	"os"

	"irptools/signals/signal"
	signalutils "irptools/signals/utils"
	"irptools/tools/stat"
	"irptools/tools/utils"
	"irptools/utils/errs"
	jsonutils "irptools/utils/json"
	"irptools/utils/logs"
)

func Main(ctx context.Context, cfg Config) error {
	return utils.DecorateExecution(ctx, "FILTER", func(ctx context.Context) error {
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
	err = execFilter(ctx, execCfg)
	if err != nil {
		return errs.Wrap(err)
	}

	if cfg.Target.WithStat {
		if _, err = os.Stat(execCfg.Target.Folder.Path); !os.IsNotExist(err) {
			statCfg := stat.Config{
				Target: cfg.Target.Folder.Join("stat"),
				Source: execCfg.Target.Folder.Path,
			}
			err = stat.Main(ctx, statCfg)
			if err != nil {
				return errs.Wrap(err)
			}
		}
	}

	return nil
}

func execFilter(ctx context.Context, cfg Config) error {
	jsonPred, err := jsonutils.BuildPredicate(cfg.Filter, jsonutils.DefaultLogic())
	if err != nil {
		return errs.Errorf("failed to build predicate: %w", err)
	}

	var sPtr *signal.Signal
	sObj := newSignalObject(&sPtr)

	filter := func(s signal.Signal) (bool, error) {
		sPtr = &s
		res := jsonPred.Is(sObj)
		return res, nil
	}

	getConsumer := func(filePath string) (signalutils.ClosableSignalConsumer, error) {
		postponing := signalutils.NewPostponingConsumer(func() (signalutils.ClosableSignalConsumer, error) {
			return signalutils.NewJsonFileWriter(filePath, cfg.Target.PrettyJsonPrint)
		})
		filtering := signalutils.NewFilteringSignalConsumer(postponing, filter)
		return filtering, nil
	}

	getTargetFilePath := signalutils.RepeatSourceTreeTargetFilePathStrategy(cfg.Source, cfg.Target.Folder.Path)
	if cfg.Target.ToOneFolder {
		getTargetFilePath = signalutils.ToOneFolderTargetFilePathStrategy(cfg.Source, cfg.Target.Folder.Path)
	}
	factory := signalutils.NewSignalsToFileConsumersFactory(getConsumer, getTargetFilePath)
	err = signalutils.EnumSignals(ctx, cfg.Source, factory.NewConsumer)
	if err != nil {
		return errs.Wrap(err)
	}

	return nil
}

func newSignalObject(sr **signal.Signal) jsonutils.MappedObject {
	return jsonutils.NewMappedObject(map[string]func() (any, error){
		"id":        func() (any, error) { return (*sr).Id, nil },
		"source":    func() (any, error) { return (*sr).Source, nil },
		"brand":     func() (any, error) { return (*sr).Brand, nil },
		"device":    func() (any, error) { return (*sr).Device, nil },
		"protocol":  func() (any, error) { return (*sr).Protocol, nil },
		"function":  func() (any, error) { return (*sr).Function, nil },
		"frequency": func() (any, error) { return (*sr).Frequency, nil },
		"data":      func() (any, error) { return (*sr).Data, nil },
		"Source":    func() (any, error) { return (*sr).Source, nil },
		"Brand":     func() (any, error) { return (*sr).Brand, nil },
		"Device":    func() (any, error) { return (*sr).Device, nil },
		"Protocol":  func() (any, error) { return (*sr).Protocol, nil },
		"Function":  func() (any, error) { return (*sr).Function, nil },
		"Frequency": func() (any, error) { return (*sr).Frequency, nil },
		"Data":      func() (any, error) { return (*sr).Data, nil },
	})
}
