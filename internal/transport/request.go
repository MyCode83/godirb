package transport

import "sync"

type RequestOptions struct {
	mu sync.Mutex

	URL        string
	Method     Method
	MethodMode MethodMode

	Headers []string
}

func (opts *RequestOptions) methodForRequest() Method {
	if opts.MethodMode != MethodModeSwitch {
		return opts.Method
	}

	opts.mu.Lock()
	defer opts.mu.Unlock()

	opts.Method.Toggle()

	return opts.Method
}
