package cli

import (
	"time"
)

var PreUserAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36",
}

type Config struct {
	Test string

	// Flags
	URL        string
	BaseURL    string
	Threads    int
	Depth      int
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
	Debug     bool
	Version   bool
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
