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

	resp, err := client.Do(&RequestOptions{
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

func TestDoSendsHeaderValueContainingColon(t *testing.T) {
	client, url, cleanup := newTestClient(t, func(ctx *fasthttp.RequestCtx) {
		if got := string(ctx.Request.Header.Peek("Authorization")); got != "Bearer part:with:colons" {
			t.Fatalf("Authorization header = %q, want %q", got, "Bearer part:with:colons")
		}

		ctx.SetStatusCode(fasthttp.StatusOK)
	})
	defer cleanup()

	resp, err := client.Do(&RequestOptions{
		URL:     url,
		Method:  MethodGET,
		Headers: []string{"Authorization: Bearer part:with:colons"},
	})
	if err != nil {
		t.Fatalf("Do returned error: %v", err)
	}
	if resp.StatusCode != fasthttp.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, fasthttp.StatusOK)
	}
}

func TestDoRejectsInvalidHeader(t *testing.T) {
	var req fasthttp.Request

	if err := applyHeaders(&req, []string{"X-Test"}); err == nil {
		t.Fatal("applyHeaders returned nil error for invalid header")
	}
}

func TestDoSendsRepeatedHeaders(t *testing.T) {
	client, url, cleanup := newTestClient(t, func(ctx *fasthttp.RequestCtx) {
		if got := string(ctx.Request.Header.Peek("Authorization")); got != "Bearer TOKEN" {
			t.Fatalf("Authorization header = %q, want %q", got, "Bearer TOKEN")
		}
		if got := string(ctx.Request.Header.Peek("X-Test")); got != "value" {
			t.Fatalf("X-Test header = %q, want %q", got, "value")
		}

		ctx.SetStatusCode(fasthttp.StatusOK)
	})
	defer cleanup()

	resp, err := client.Do(&RequestOptions{
		URL:    url,
		Method: MethodGET,
		Headers: []string{
			"Authorization: Bearer TOKEN",
			"X-Test: value",
		},
	})
	if err != nil {
		t.Fatalf("Do returned error: %v", err)
	}
	if resp.StatusCode != fasthttp.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, fasthttp.StatusOK)
	}
}

func TestDoCopiesResponseBody(t *testing.T) {
	bodies := []string{"first", "second-body"}
	requests := 0
	client, url, cleanup := newTestClient(t, func(ctx *fasthttp.RequestCtx) {
		if requests >= len(bodies) {
			t.Fatalf("unexpected request %d", requests+1)
		}

		ctx.SetBodyString(bodies[requests])
		requests++
	})
	defer cleanup()

	first, err := client.Do(&RequestOptions{
		URL:    url,
		Method: MethodGET,
	})
	if err != nil {
		t.Fatalf("first Do returned error: %v", err)
	}

	_, err = client.Do(&RequestOptions{
		URL:    url,
		Method: MethodGET,
	})
	if err != nil {
		t.Fatalf("second Do returned error: %v", err)
	}

	if got := string(first.Body); got != "first" {
		t.Fatalf("first body after second request = %q, want %q", got, "first")
	}
}

func TestDoSwitchModeRotatesMethods(t *testing.T) {
	var got []string
	client, url, cleanup := newTestClient(t, func(ctx *fasthttp.RequestCtx) {
		got = append(got, string(ctx.Method()))
	})
	defer cleanup()

	method, mode, err := ParseMethod("SWITCH")
	if err != nil {
		t.Fatalf("ParseMethod returned error: %v", err)
	}

	opts := RequestOptions{
		URL:        url,
		Method:     method,
		MethodMode: mode,
	}
	for range 3 {
		if _, err := client.Do(&opts); err != nil {
			t.Fatalf("Do returned error: %v", err)
		}
	}

	want := []string{"GET", "HEAD", "GET"}
	if len(got) != len(want) {
		t.Fatalf("got %d requests, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("method %d = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestDoSwitchModeKeepsStateInRequestOptions(t *testing.T) {
	var got []string
	client, url, cleanup := newTestClient(t, func(ctx *fasthttp.RequestCtx) {
		got = append(got, string(ctx.Method()))
	})
	defer cleanup()

	optsA := RequestOptions{
		URL:        url,
		Method:     MethodHEAD,
		MethodMode: MethodModeSwitch,
	}
	optsB := RequestOptions{
		URL:        url,
		Method:     MethodGET,
		MethodMode: MethodModeSwitch,
	}

	if _, err := client.Do(&optsA); err != nil {
		t.Fatalf("Do optsA returned error: %v", err)
	}
	if _, err := client.Do(&optsB); err != nil {
		t.Fatalf("Do optsB returned error: %v", err)
	}

	want := []string{"GET", "HEAD"}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("method %d = %q, want %q", i, got[i], want[i])
		}
	}
}
