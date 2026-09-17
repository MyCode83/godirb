package assemble

import (
	"context"
	"sync"

	"github.com/MyCode83/godirb/internal/cli"
	"github.com/MyCode83/godirb/internal/core"
	"github.com/MyCode83/godirb/internal/transport"
)

func BuildEngine(
	contextCancel context.Context,
	cancel context.CancelFunc,
	cfg cli.Config,
	mode core.Mode,
	client *transport.Client,
	method transport.Method,
	methodMode transport.MethodMode,
	auth string,
	wl []string,
	wg *sync.WaitGroup,
) *core.Core {
	limiter := make(chan struct{}, cfg.Threads)
	var dirsChan chan core.DirTask
	if mode == core.ModeDir {
		dirsChan = make(chan core.DirTask, cfg.Threads*50)
	}

	return &core.Core{
		// Mode
		Mode: mode,

		// Bools
		Recursive: cfg.Recursive,

		// Context
		Ctx:    contextCancel,
		Cancel: cancel,
		// Config
		Timeout: cfg.Timeout,
		Delay:   cfg.Delay,
		Depth:   cfg.Depth,
		Quiet:   cfg.Quiet,
		Debug:   cfg.Debug,

		// HTTP
		Client:     client,
		Method:     method,
		MethodMode: methodMode,
		UserAgents: cfg.UserAgent,
		AuthHeader: auth,
		Header:     cfg.Header,
		// cfg.Placeholder
		Placeholder: cfg.Placeholder,
		// Control
		IgnoreCodes: cfg.IgnoreCode,
		Exts:        cfg.Exts,

		// Concurrency
		Limiter:  limiter,
		DirsChan: dirsChan,

		// WG
		WG: wg,

		// WordList
		WL: wl,

		// State
		VisitedDirs: make(map[string]bool),
	}
}
