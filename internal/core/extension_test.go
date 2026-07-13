package core

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/MyCode83/godirb/internal/calibration"
	"github.com/MyCode83/godirb/internal/transport"
	"github.com/valyala/fasthttp"
)

func TestProcessExtensionsEmitsOnlyUnfilteredResults(t *testing.T) {
	var requestedPaths []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestedPaths = append(requestedPaths, r.URL.Path)

		switch r.URL.Path {
		case "/asset.txt":
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "found")
		case "/asset.skip":
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "missing")
		default:
			t.Fatalf("unexpected request path %q", r.URL.Path)
		}
	}))
	defer server.Close()

	c := &Core{
		Client:      transport.New(&fasthttp.Client{}),
		Ctx:         context.Background(),
		Method:      transport.MethodGET,
		MethodMode:  transport.MethodModeFixed,
		UserAgents:  []string{"godirb-test"},
		IgnoreCodes: []int{http.StatusNotFound},
		Exts:        []string{"txt", "bad-url", "skip"},
		Calibration: &calibration.Calibration{},
	}
	results := make(chan Result, len(c.Exts))
	request := &transport.RequestOptions{MethodMode: transport.MethodModeFixed}

	ok := c.processExtensions(
		request,
		results,
		"FILE",
		"test-ext",
		func(ext string) (string, error) {
			if ext == "bad-url" {
				return "", fmt.Errorf("bad extension")
			}
			return server.URL + "/asset." + ext, nil
		},
	)

	if !ok {
		t.Fatal("processExtensions returned false, want true")
	}

	got := drainResults(results)
	want := []Result{{
		Prefix: "FILE",
		URL:    server.URL + "/asset.txt",
		Size:   len("found"),
		Status: http.StatusOK,
	}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("results = %#v, want %#v", got, want)
	}

	if wantPaths := []string{"/asset.txt", "/asset.skip"}; !reflect.DeepEqual(requestedPaths, wantPaths) {
		t.Fatalf("requested paths = %#v, want %#v", requestedPaths, wantPaths)
	}
}

func TestProcessExtensionsContinuesOnCalibrationMatch(t *testing.T) {
	var requestedPaths []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestedPaths = append(requestedPaths, r.URL.Path)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "baseline")
	}))
	defer server.Close()

	c := &Core{
		Client:     transport.New(&fasthttp.Client{}),
		Ctx:        context.Background(),
		Method:     transport.MethodGET,
		MethodMode: transport.MethodModeFixed,
		UserAgents: []string{"godirb-test"},
		Exts:       []string{"one", "two"},
		Calibration: &calibration.Calibration{
			Status:    http.StatusOK,
			Length:    len("baseline"),
			Tolerance: 0,
			Stable:    true,
		},
	}
	results := make(chan Result, len(c.Exts))
	request := &transport.RequestOptions{MethodMode: transport.MethodModeFixed}

	ok := c.processExtensions(
		request,
		results,
		"FILE",
		"test-ext",
		func(ext string) (string, error) {
			return server.URL + "/asset." + ext, nil
		},
	)

	if !ok {
		t.Fatal("processExtensions returned false, want true")
	}

	if got := drainResults(results); len(got) != 0 {
		t.Fatalf("results = %#v, want none", got)
	}

	if wantPaths := []string{"/asset.one", "/asset.two"}; !reflect.DeepEqual(requestedPaths, wantPaths) {
		t.Fatalf("requested paths = %#v, want %#v", requestedPaths, wantPaths)
	}
}

func TestProcessExtensionsDelaysIgnoredResponses(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "missing")
	}))
	defer server.Close()

	c := &Core{
		Client:      transport.New(&fasthttp.Client{}),
		Ctx:         context.Background(),
		Method:      transport.MethodGET,
		MethodMode:  transport.MethodModeFixed,
		UserAgents:  []string{"godirb-test"},
		IgnoreCodes: []int{http.StatusNotFound},
		Exts:        []string{"one", "two"},
		Calibration: &calibration.Calibration{},
		Delay:       20 * time.Millisecond,
	}
	results := make(chan Result, len(c.Exts))
	request := &transport.RequestOptions{MethodMode: transport.MethodModeFixed}

	start := time.Now()
	ok := c.processExtensions(
		request,
		results,
		"FILE",
		"test-ext",
		func(ext string) (string, error) {
			return server.URL + "/asset." + ext, nil
		},
	)
	elapsed := time.Since(start)

	if !ok {
		t.Fatal("processExtensions returned false, want true")
	}

	if got := drainResults(results); len(got) != 0 {
		t.Fatalf("results = %#v, want none", got)
	}

	if elapsed < 40*time.Millisecond {
		t.Fatalf("elapsed = %s, want at least %s", elapsed, 40*time.Millisecond)
	}
}

func drainResults(results <-chan Result) []Result {
	var drained []Result
	for {
		select {
		case result := <-results:
			drained = append(drained, result)
		default:
			return drained
		}
	}
}
