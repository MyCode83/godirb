package core

import (
	"context"
	"sync"
	"time"

	"github.com/MyCode83/godirb/internal/calibration"
	"github.com/MyCode83/godirb/internal/transport"
	"github.com/MyCode83/godirb/internal/signature"
)

type Mode int

const (
	ModeDir Mode = iota
	ModeFuzz
	ModePort
)

type Core struct {
	// Mode
	Mode Mode

	// Bools
	Recursive bool

	// infra
	Client *transport.Client
	Ctx    context.Context
	Cancel context.CancelFunc

	// Config

	Method      transport.Method
	MethodMode  transport.MethodMode
	Placeholder string
	UserAgents  []string
	IgnoreCodes []int
	Exts        []string
	Header      []string
	AuthHeader  string
	Delay       time.Duration
	Timeout     time.Duration
	Depth       int
	Calibration *calibration.Calibration
	Quiet       bool
	Debug       bool

	Signatures *signature.Matcher

	// Concurrency
	Limiter  chan struct{}
	DirsChan chan DirTask

	WG *sync.WaitGroup
	WL []string

	// State
	VisitedDirs  map[string]bool
	VisitedMutex sync.Mutex
	MethodMutex  sync.Mutex
}

type DirTask struct {
	URL   string
	Depth int
}

func (c *Core) nextRequestMethod() transport.Method {
	if c.MethodMode != transport.MethodModeSwitch {
		return c.Method
	}

	c.MethodMutex.Lock()
	defer c.MethodMutex.Unlock()

	c.Method.Toggle()

	return c.Method
}
