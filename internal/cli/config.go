package cli
type Config struct {
	// Flags
	URL        string
	Threads    int
	IgnoreCode []int
	Exts       []string
	RawTimeout  string
	RawDelay	string
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

}