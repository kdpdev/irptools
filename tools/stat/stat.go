package stat

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strconv"

	"irptools/signals/signal"
	signalutils "irptools/signals/utils"
	"irptools/tools/utils"
	"irptools/utils/alg"
	"irptools/utils/errs"
	"irptools/utils/logs"
)

func Main(ctx context.Context, cfg Config) error {
	return utils.DecorateExecution(ctx, "STAT", func(ctx context.Context) error {
		return execMain(ctx, cfg)
	})
}

func execMain(ctx context.Context, cfg Config) error {
	err := errs.CheckValid(cfg, "config")
	if err != nil {
		return err
	}

	l := logs.L(ctx)
	l.I("source <-: %s", cfg.Source)
	l.I("target ->: %s", cfg.Target.Path)

	cfg.Target, err = cfg.Target.PrepareTarget()
	if err != nil {
		return errs.Wrap(err)
	}

	stat, err := collectSignalsStat(ctx, cfg.Source)
	if err != nil {
		return errs.Wrap(err)
	}

	err = storeData(filepath.Join(cfg.Target.Path, "all.json"), stat)
	if err != nil {
		return errs.Errorf("failed to store all: %w", err)
	}

	err = storeData(filepath.Join(cfg.Target.Path, "ints.json"), stat.Ints)
	if err != nil {
		return errs.Errorf("failed to store ints: %w", err)
	}

	for k, v := range stat.Strs {
		fileName := k + ".json"
		list := alg.MapKeys(v)
		sort.Strings(list)
		err = storeData(filepath.Join(cfg.Target.Path, fileName), []interface{}{
			map[string]int{"count": len(v)},
			map[string]interface{}{k: v},
			map[string]interface{}{"list": list},
		})

		if err != nil {
			return errs.Errorf("failed to store %s: %w", k, err)
		}
	}

	logs.L(ctx).I("results -> %s", cfg.Target.Path)

	return nil
}

func storeData(filePath string, data interface{}) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return errs.Wrap(err)
	}

	err = os.WriteFile(filePath, jsonData, 0644)
	if err != nil {
		return errs.Wrap(err)
	}

	return nil
}

func addSignalToStatItem(item *Item, s signal.Signal) {
	item.IncInt("Signals", 1)
	item.AddStr("Brands", s.Brand)
	item.AddStr("Devices", s.Device)
	item.AddStr("Functions", s.Function)
	item.AddStr("Protocols", s.Protocol)
	item.AddStr("Frequencies", strconv.Itoa(int(s.Frequency)))
}

type statCollector struct {
	rootStat *Item
	stat     Item
}

func (this *statCollector) Consume(s signal.Signal) error {
	this.stat.IncInt("Signals", 1)
	this.stat.AddStr("Brands", s.Brand)
	this.stat.AddStr("Devices", s.Device)
	this.stat.AddStr("Functions", s.Function)
	this.stat.AddStr("Protocols", s.Protocol)
	this.stat.AddStr("Frequencies", strconv.Itoa(int(s.Frequency)))
	return nil
}

func (this *statCollector) Close() error {
	this.stat.IncInt("Files", 1)
	this.rootStat.AddItem(this.stat)
	return nil
}

func collectSignalsStat(ctx context.Context, sourcePath string) (Item, error) {
	rootStat := NewItem()
	err := signalutils.EnumSignals(ctx, sourcePath, func(filePath string) (signalutils.ClosableSignalConsumer, error) {
		return &statCollector{rootStat: &rootStat, stat: NewItem()}, nil
	})
	if err != nil {
		return rootStat, errs.Wrap(err)
	}
	return rootStat, nil
}
