package cli

import (
	"godirb/internal/core"
	"sync"
	"time"

	"github.com/valyala/fasthttp"
)

var useColors = true

const banner string = (`		                     
   ____ _  ____   ____/ /   (_)   _____   / /_
  / __  / / __ \ / __  /   / /   / ___/  / __ \
 / /_/ / / /_/ // /_/ /   / /   / /     / /_/ /
 \__  /  \____/ \____/   /_/   /_/     /_____/
/____/
`)

var (
	client       *fasthttp.Client
	wg           sync.WaitGroup
	tasksWG      sync.WaitGroup
	visitedMutex sync.Mutex
	mode         core.Mode = core.ModeDir
)

const version = "0.9.0"

var preUserAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36",
}

type Config struct {
	Test string

	// Flags
	URL        string
	BaseURL    string
	Threads    int
	IgnoreCode []int
	Exts       []string
	RawTimeout string
	RawDelay   string
	Timeout    time.Duration
	Delay      time.Duration
	UserAgent  []string
	// bool
	NoColor   bool
	Recursive bool
	//forces
	ForceHead    bool
	ForceThreads bool
	ForceProxy   bool

	Proxy  string
	Method string
	// Basic Auth
	Username string
	Password string
	// Placeholder
	Placeholder string

	Header []string
	// TLS
	Insecure bool

	Quiet bool

	JSON   bool
	CSV    bool
	Output string
}
