package main

import (
	"fmt"
	"sync"
	"os"


	"github.com/MyCode83/godirb/internal/assemble"
	"github.com/MyCode83/godirb/internal/cli"
	"github.com/MyCode83/godirb/internal/core"
	"github.com/MyCode83/godirb/internal/debug"
)

var version = "dev"

func reportIssue() {
	fmt.Fprintf(os.Stderr, "Please report it at https://github.com/MyCode83/godirb/issues")
}

func main() {
	assemble.ConfigureProcess()
	contextCancel, cancel := assemble.SetupSignals()

	cfg, wd := cli.ParseFlags()
	if cfg.Version {
		fmt.Println(version)
		return
	}
	assemble.LogParsedFlags(cfg, wd)
	cli.ValidateFlags(&cfg)
	assemble.LogValidatedFlags(cfg)
	mode := cli.SelectMode(core.ModeDir, cfg)
	assemble.LogSelectedMode(mode)

	// wd = instance
	// wl = wordlist slice
	method, methodMode, err := assemble.BuildMethod(cfg)
	if err != nil {
		debug.Error("method parse", err)
		fmt.Fprintf(os.Stderr, "[X] Error: invalid method '%s'\n", cfg.Method)
		os.Exit(2)
	}
	assemble.ApplyModeDefaults(&cfg, &wd, mode)

	client := assemble.BuildClient(cfg)
	assemble.ValidateURL(cfg, client, mode, method, methodMode)

	wl := assemble.BuildWordlist(wd)
	auth := assemble.BuildAuth(cfg)

	outputFormat, quietOutput := assemble.BuildOutputConfig(cfg)
	assemble.PrintBanner(cfg, mode, quietOutput)

	var wg sync.WaitGroup
	engine := assemble.BuildEngine(
		contextCancel,
		cancel,
		cfg,
		mode,
		client,
		method,
		methodMode,
		auth,
		wl,
		&wg,
	)

	err = assemble.AttachSignatures(engine)
	if err != nil {
		debug.Error("signature matcher", err)
		fmt.Fprintf(os.Stderr, "[X] Error loading response signatures: %v\n", err)
		reportIssue()
		os.Exit(2)
	}

	err = assemble.BuildCalibration(engine, &cfg, mode, client, &method, &methodMode)
	if err != nil {
		debug.Error("calibration build", err)
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(2)
	}

	stream, err := assemble.BuildOutputStream(outputFormat, cfg)
	if err != nil {
		debug.Error("output stream open", err)
		fmt.Fprintf(os.Stderr, "[X] Error writing output: %v\n", err)
		reportIssue()
		os.Exit(1)
	}

	err = assemble.Run(engine, stream, cfg.BaseURL, cancel)
	if err != nil {
		debug.Error("writing output err at Run", err)
		fmt.Fprintf(os.Stderr, "[X] Error writing output: %v\n", err)
		reportIssue()
		os.Exit(1)
	}
}
