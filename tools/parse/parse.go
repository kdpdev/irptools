package parse

import (
	"context"
	"fmt"
	"strings"

	"irptools/signals/signal"
	"irptools/signals/sources/fz"
	"irptools/signals/sources/visio"
	signalutils "irptools/signals/utils"
	"irptools/tools/stat"
	"irptools/tools/utils"
	"irptools/utils/errs"
	"irptools/utils/logs"
)

func Main(ctx context.Context, cfg Config) error {
	return utils.DecorateExecution(ctx, "PARSE", func(ctx context.Context) error {
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
	logs.L(ctx).I("target ->: %s", cfg.Target.Folder.Path)

	execCfg := cfg
	execCfg.Target.Folder = execCfg.Target.Folder.Join("result")
	err = execParse(ctx, execCfg)
	if err != nil {
		return errs.Wrap(err)
	}

	if cfg.Target.WithStat {
		statCfg := stat.Config{
			Target: cfg.Target.Folder.Join("stat"),
			Source: execCfg.Target.Folder.Path,
		}
		err = stat.Main(ctx, statCfg)
		if err != nil {
			return errs.Wrap(err)
		}
	}

	return nil
}

func execParse(ctx context.Context, cfg Config) error {
	l := logs.L(ctx)

	var err error
	parsedSignalsCount := 0
	for k, sourceCfg := range cfg.Sources {
		scopedCtx, _ := logs.WithScope(ctx, "  ")
		scope := fmt.Sprintf("%s.%s", k, sourceCfg.Type)
		err = utils.DecorateExecution(scopedCtx, scope, func(ctx context.Context) error {
			l := logs.L(ctx)
			if sourceCfg.Skip {
				l.I("skipped")
				return nil
			}

			targetCfg := cfg.Target
			targetCfg.Folder = targetCfg.Folder.Join(k)
			count, err := parseSource(ctx, sourceCfg, targetCfg)
			l.I("parsed: %v -> %s", count, targetCfg.Folder.Path)
			parsedSignalsCount += count
			if err != nil {
				return errs.Errorf("failed to parse %s: %w", scope, err)
			}
			return nil
		})

		if err != nil {
			break
		}
	}

	l.I("signals count = %v", parsedSignalsCount)
	l.I("results -> %s", cfg.Target.Folder.Path)

	return errs.Wrap(err)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func parseSource(ctx context.Context, sourceCfg SourceConfig, targetCfg TargetConfig) (int, error) {
	if sourceCfg.Skip {
		return 0, nil
	}

	parse, ok := getAdaptedParsers()[sourceCfg.Type]
	if !ok {
		return 0, errs.Errorf("unknown source type: '%s'", sourceCfg.Type)
	}

	consumers, err := newSignalConsumersFactory(sourceCfg, targetCfg, func(filePath string) (signalutils.ClosableSignalConsumer, error) {
		return signalutils.NewJsonFileWriter(filePath, targetCfg.PrettyJsonPrint)
	})
	if err != nil {
		return 0, errs.Wrap(err)
	}

	count, err := parse(ctx, sourceCfg.Path, sourceCfg.Options, consumers)
	if err != nil {
		return count, errs.Errorf("failed to parse '%s' source: %w", sourceCfg.Type, err)
	}

	return count, nil
}

type AdaptedParseFn = func(ctx context.Context, path string, opts SourceOptions, consumerFactory *signalutils.SignalsToFileConsumersFactory) (int, error)

func getAdaptedParsers() map[string]AdaptedParseFn {
	parsers := map[string]AdaptedParseFn{
		"fz":    adaptFzParser(),
		"visio": adaptVisioParser(),
	}
	return parsers
}

func adaptFzParser() AdaptedParseFn {
	return func(ctx context.Context, path string, opts SourceOptions, consumerFactory *signalutils.SignalsToFileConsumersFactory) (int, error) {
		return fz.ParseIrFiles(ctx, path, opts, func(filePath string) (fz.ClosableSignalConsumer, error) {
			return consumerFactory.NewConsumer(filePath)
		})
	}
}

func adaptVisioParser() AdaptedParseFn {
	return func(ctx context.Context, path string, opts SourceOptions, consumerFactory *signalutils.SignalsToFileConsumersFactory) (int, error) {
		return visio.ParseCsvFiles(ctx, path, opts, func(filePath string) (visio.ClosableSignalConsumer, error) {
			return consumerFactory.NewConsumer(filePath)
		})
	}
}

func newSignalConsumersFactory(
	sourceCfg SourceConfig,
	targetCfg TargetConfig,
	getConsumer signalutils.SignalsToFileConsumerSourceFn) (*signalutils.SignalsToFileConsumersFactory, error) {

	var trs []signalutils.TransformSignalFn
	if !targetCfg.KeepSourceField {
		trs = append(trs, func(signal signal.Signal) (signal.Signal, error) {
			signal.Source = ""
			return signal, nil
		})
	}

	if targetCfg.FieldsToLower {
		toLower := func(str *string) {
			*str = strings.ToLower(*str)
		}
		trs = append(trs, func(signal signal.Signal) (signal.Signal, error) {
			toLower(&signal.Source)
			toLower(&signal.Brand)
			toLower(&signal.Device)
			toLower(&signal.Function)
			toLower(&signal.Protocol)
			return signal, nil
		})
	}

	if len(trs) != 0 {
		prev := getConsumer
		getConsumer = func(filePath string) (signalutils.ClosableSignalConsumer, error) {
			encoder, err := prev(filePath)
			if err != nil {
				return nil, errs.Wrap(err)
			}
			return signalutils.NewTransformingSignalConsumer(encoder, trs), nil
		}
	}

	getTargetFilePath := signalutils.RepeatSourceTreeTargetFilePathStrategy(sourceCfg.Path, targetCfg.Folder.Path)

	return signalutils.NewSignalsToFileConsumersFactory(getConsumer, getTargetFilePath), nil
}
