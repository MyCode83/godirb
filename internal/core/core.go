package core

import (
	"context"
	"github.com/MyCode83/godirb/internal/baseline"
	"github.com/MyCode83/godirb/internal/transport"
	"github.com/MyCode83/godirb/internal/wildcard"

	"github.com/fatih/color"
	"sync"
	"time"
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
	Wildcard    *wildcard.Wildcard
	Baseline    *baseline.Baseline
	Quiet       bool
	Debug       bool

	// Colors
	Others *color.Color
	File   *color.Color

	// Concurrency
	Limiter  chan struct{}
	DirsChan chan string

	WG *sync.WaitGroup
	WL []string

	// State
	VisitedDirs  map[string]bool
	VisitedMutex sync.Mutex
	MethodMutex  sync.Mutex
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
