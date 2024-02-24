package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	"irptools/tools/filter"
	"irptools/tools/parse"
	"irptools/tools/stat"
	"irptools/utils/alg"
	"irptools/utils/errs"
)

func main() {

	const defaultCmd = ""
	const defaultCfg = "cfg_$cmd$.json"

	cmds := map[string]func(ctx context.Context, cfg string) error{
		"parse":  makeExecCmdFn(parse.Main, parse.LoadConfig),
		"stat":   makeExecCmdFn(stat.Main, stat.LoadConfig),
		"filter": makeExecCmdFn(filter.Main, filter.LoadConfig),
	}

	var cmdLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	cmd := cmdLine.String("cmd", defaultCmd, fmt.Sprintf("command: %s", alg.MapKeys(cmds)))
	cfg := cmdLine.String("cfg", defaultCfg, "config")
	err := cmdLine.Parse(os.Args[1:])
	if err != nil {
		log.Fatal(err)
		return
	}

	*cfg = strings.ReplaceAll(*cfg, "$cmd$", *cmd)

	l := log.Default()
	l.Println("cmd =", *cmd)
	l.Println("cfg =", *cfg)

	execute, ok := cmds[*cmd]
	if !ok {
		cmdLine.Usage()
		l.Fatalf("FAILED: Unknown command: '%s'", *cmd)
		return
	}

	ctx := context.Background()
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()
	err = execute(ctx, *cfg)
	if err != nil {
		l.Fatalf("FAILED: cmd = '%s': %v", *cmd, err)
		return
	}

	fmt.Print("DONE")
}

type execFn[T any] func(ctx context.Context, cfg T) error
type loadCfgFn[T any] func(cfgPath string) (T, error)

func makeExecCmdFn[T any](exec execFn[T], loadCfg loadCfgFn[T]) func(ctx context.Context, cfg string) error {
	return func(ctx context.Context, cfgPath string) error {
		cfg, err := loadCfg(cfgPath)
		if err != nil {
			return errs.Errorf("failed to load config: %w", err)
		}
		return errs.Wrap(exec(ctx, cfg))
	}
}
