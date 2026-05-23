package debug

import (
	"fmt"
	"log"
	"os"
	"sync/atomic"

	"github.com/valyala/fasthttp"
)

var enabled atomic.Bool

func Set(value bool) {
	enabled.Store(value)
	if value {
		log.SetOutput(os.Stderr)
		log.SetFlags(log.LstdFlags | log.Lmicroseconds)
		Printf("debug enabled")
	}
}

func Enabled() bool {
	return enabled.Load()
}

func Printf(format string, args ...any) {
	if !Enabled() {
		return
	}
	log.Printf("[debug] "+format, args...)
}

func Request(label string, request *fasthttp.Request) {
	if !Enabled() {
		return
	}
	headers := 0
	request.Header.VisitAll(func(_, _ []byte) {
		headers++
	})
	Printf("%s request method=%s url=%s ua=%q headers=%d body=%d",
		label,
		request.Header.Method(),
		request.URI().FullURI(),
		request.Header.UserAgent(),
		headers,
		len(request.Body()),
	)
}

func Response(label string, response *fasthttp.Response) {
	if !Enabled() {
		return
	}
	headers := 0
	response.Header.VisitAll(func(_, _ []byte) {
		headers++
	})
	Printf("%s response status=%d body=%d headers=%d",
		label,
		response.StatusCode(),
		len(response.Body()),
		headers,
	)
}

func Error(label string, err error) {
	if !Enabled() || err == nil {
		return
	}
	Printf("%s error=%s", label, err)
}

func Value(label string, value any) {
	if !Enabled() {
		return
	}
	Printf("%s=%s", label, fmt.Sprint(value))
}
