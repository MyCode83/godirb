package core

import (
	"context"
	"godirb/internal/baseline"
	"godirb/internal/wildcard"

	"github.com/fatih/color"
	"sync"
	"time"

	"github.com/valyala/fasthttp"
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
	Client *fasthttp.Client
	Ctx context.Context
	Cancel context.CancelFunc


	// Config

	Method	 string
	Placeholder string
	UserAgents	 []string
	IgnoreCodes	 []int
	Exts	[]string
	Header []string
	AuthHeader string
	Delay	time.Duration
	Timeout time.Duration
	Wildcard *wildcard.Wildcard
	Baseline *baseline.Baseline
	Quiet bool



	// Colors
	Others *color.Color
	File *color.Color
	
	// Concurrency
	Limiter chan struct{}
	DirsChan chan string

	WG *sync.WaitGroup
	WL []string


	// State
	VisitedDirs map[string]bool
	VisitedMutex sync.Mutex
}
