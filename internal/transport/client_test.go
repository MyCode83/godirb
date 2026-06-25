package transport

import (
	"net"
	"testing"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttputil"
)

func newTestClient(t *testing.T, handler fasthttp.RequestHandler) (*Client, string, func()) {
	t.Helper()

	ln := fasthttputil.NewInmemoryListener()
	server := &fasthttp.Server{Handler: handler}
	done := make(chan error, 1)

	go func() {
		done <- server.Serve(ln)
	}()

	raw := &fasthttp.Client{
		Dial: func(addr string) (net.Conn, error) {
			return ln.Dial()
		},
	}

	cleanup := func() {
		if err := server.Shutdown(); err != nil {
			t.Fatalf("shutdown test server: %v", err)
		}
		if err := <-done; err != nil {
			t.Fatalf("serve test server: %v", err)
		}
	}

	return New(raw), "http://test.local", cleanup
}

func TestDoSendsRequestWithHeaders(t *testing.T) {
	client, url, cleanup := newTestClient(t, func(ctx *fasthttp.RequestCtx) {
		if got := string(ctx.Request.Header.Peek("X-Test")); got != "ok" {
			t.Fatalf("X-Test header = %q, want %q", got, "ok")
		}

		ctx.SetStatusCode(fasthttp.StatusTeapot)
	})
	defer cleanup()

	resp, err := client.Do(RequestOptions{
		URL:     url,
		Method:  MethodGET,
		Headers: []string{"X-Test: ok"},
	})
	if err != nil {
		t.Fatalf("Do returned error: %v", err)
	}
	if resp.StatusCode != fasthttp.StatusTeapot {
		t.Fatalf("status = %d, want %d", resp.StatusCode, fasthttp.StatusTeapot)
	}
}
