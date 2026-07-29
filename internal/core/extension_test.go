package core

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/MyCode83/godirb/internal/calibration"
	"github.com/MyCode83/godirb/internal/transport"
	"github.com/MyCode83/godirb/internal/urlutil"
	"github.com/valyala/fasthttp"
)

func TestProcessExtensionsEmitsOnlyUnfilteredResults(t *testing.T) {
	var (
		mu             sync.Mutex
		requestedPaths []string
	)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/asset.txt":
			mu.Lock()
			requestedPaths = append(requestedPaths, r.URL.Path)
			mu.Unlock()
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "found")
		case "/asset.skip":
			mu.Lock()
			requestedPaths = append(requestedPaths, r.URL.Path)
			mu.Unlock()
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "missing")
		default:
			w.WriteHeader(http.StatusTeapot)
			fmt.Fprint(w, "calibration")
		}
	}))
	defer server.Close()

	c := newExtensionTestCore([]string{"txt", "skip"})
	c.IgnoreCodes = []int{http.StatusNotFound}

	buildCalibrationURL := extensionCalibrationURLBuilder(server.URL)
	buildExtensionCalibrations(t, c, buildCalibrationURL)

	results := make(chan Result, len(c.Exts))
	request := &transport.RequestOptions{
		MethodMode: transport.MethodModeFixed,
	}

	ok := c.processExtensions(
		request,
		results,
		"FILE",
		"test-ext",
		func(ext string) string {
			return urlutil.AddExtension(server.URL+"/asset", ext)
		},
		buildCalibrationURL,
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

	mu.Lock()
	paths := append([]string(nil), requestedPaths...)
	mu.Unlock()

	wantPaths := []string{"/asset.txt", "/asset.skip"}
	if !reflect.DeepEqual(paths, wantPaths) {
		t.Fatalf("requested paths = %#v, want %#v", paths, wantPaths)
	}
}

func TestProcessExtensionsContinuesOnCalibrationMatch(t *testing.T) {
	var (
		mu             sync.Mutex
		requestedPaths []string
	)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/asset.") {
			mu.Lock()
			requestedPaths = append(requestedPaths, r.URL.Path)
			mu.Unlock()
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "baseline")
	}))
	defer server.Close()

	c := newExtensionTestCore([]string{"one", "two"})

	buildCalibrationURL := extensionCalibrationURLBuilder(server.URL)
	buildExtensionCalibrations(t, c, buildCalibrationURL)

	results := make(chan Result, len(c.Exts))
	request := &transport.RequestOptions{
		MethodMode: transport.MethodModeFixed,
	}

	ok := c.processExtensions(
		request,
		results,
		"FILE",
		"test-ext",
		func(ext string) string {
			return urlutil.AddExtension(server.URL+"/asset", ext)
		},
		buildCalibrationURL,
	)

	if !ok {
		t.Fatal("processExtensions returned false, want true")
	}

	if got := drainResults(results); len(got) != 0 {
		t.Fatalf("results = %#v, want none", got)
	}

	mu.Lock()
	paths := append([]string(nil), requestedPaths...)
	mu.Unlock()

	wantPaths := []string{"/asset.one", "/asset.two"}
	if !reflect.DeepEqual(paths, wantPaths) {
		t.Fatalf("requested paths = %#v, want %#v", paths, wantPaths)
	}
}

func TestProcessExtensionsSkipsMissingCalibration(t *testing.T) {
	var (
		mu             sync.Mutex
		requestedPaths []string
	)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/asset.") {
			mu.Lock()
			requestedPaths = append(requestedPaths, r.URL.Path)
			mu.Unlock()
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "found")
	}))
	defer server.Close()

	c := newExtensionTestCore([]string{"known", "missing"})
	buildCalibrationURL := extensionCalibrationURLBuilder(server.URL)
	buildSingleExtensionCalibration(t, c, buildCalibrationURL, "known")

	results := make(chan Result, len(c.Exts))
	ok := c.processExtensions(
		&transport.RequestOptions{MethodMode: transport.MethodModeFixed},
		results,
		"FILE",
		"test-ext",
		func(ext string) string {
			return urlutil.AddExtension(server.URL+"/asset", ext)
		},
		buildCalibrationURL,
	)

	if !ok {
		t.Fatal("processExtensions returned false, want true")
	}

	if got := drainResults(results); len(got) != 0 {
		t.Fatalf("results = %#v, want none", got)
	}

	mu.Lock()
	paths := append([]string(nil), requestedPaths...)
	mu.Unlock()

	wantPaths := []string{"/asset.known"}
	if !reflect.DeepEqual(paths, wantPaths) {
		t.Fatalf("requested paths = %#v, want %#v", paths, wantPaths)
	}
}

func TestProcessExtensionsUsesIndependentCalibrations(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path != "/asset.one" && strings.HasSuffix(r.URL.Path, ".one"):
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "one-calibration")
		case r.URL.Path != "/asset.two" && strings.HasSuffix(r.URL.Path, ".two"):
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "two-calibration")
		case r.URL.Path == "/asset.one":
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "one-calibration")
		case r.URL.Path == "/asset.two":
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "real-two")
		default:
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "missing")
		}
	}))
	defer server.Close()

	c := newExtensionTestCore([]string{"one", "two"})
	buildCalibrationURL := extensionCalibrationURLBuilder(server.URL)
	buildExtensionCalibrations(t, c, buildCalibrationURL)

	results := make(chan Result, len(c.Exts))
	ok := c.processExtensions(
		&transport.RequestOptions{MethodMode: transport.MethodModeFixed},
		results,
		"FILE",
		"test-ext",
		func(ext string) string {
			return urlutil.AddExtension(server.URL+"/asset", ext)
		},
		buildCalibrationURL,
	)

	if !ok {
		t.Fatal("processExtensions returned false, want true")
	}

	want := []Result{{
		Prefix: "FILE",
		URL:    server.URL + "/asset.two",
		Size:   len("real-two"),
		Status: http.StatusOK,
	}}
	if got := drainResults(results); !reflect.DeepEqual(got, want) {
		t.Fatalf("results = %#v, want %#v", got, want)
	}
}

func TestProcessExtensionsDelaysIgnoredResponses(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/asset.") {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "missing")
			return
		}

		w.WriteHeader(http.StatusTeapot)
		fmt.Fprint(w, "calibration")
	}))
	defer server.Close()

	c := newExtensionTestCore([]string{"one", "two"})
	c.IgnoreCodes = []int{http.StatusNotFound}
	c.Delay = 20 * time.Millisecond

	buildCalibrationURL := extensionCalibrationURLBuilder(server.URL)
	buildExtensionCalibrations(t, c, buildCalibrationURL)

	results := make(chan Result, len(c.Exts))
	request := &transport.RequestOptions{
		MethodMode: transport.MethodModeFixed,
	}

	start := time.Now()

	ok := c.processExtensions(
		request,
		results,
		"FILE",
		"test-ext",
		func(ext string) string {
			return urlutil.AddExtension(server.URL+"/asset", ext)
		},
		buildCalibrationURL,
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

func TestProcessExtensionsCancellationReturnsCleanly(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "missing")
	}))
	defer server.Close()

	ctx, cancel := context.WithCancel(context.Background())
	c := newExtensionTestCore([]string{"one"})
	c.Ctx = ctx
	c.IgnoreCodes = []int{http.StatusNotFound}
	c.Delay = time.Hour

	buildCalibrationURL := extensionCalibrationURLBuilder(server.URL)
	buildExtensionCalibrations(t, c, buildCalibrationURL)

	done := make(chan bool, 1)
	go func() {
		done <- c.processExtensions(
			&transport.RequestOptions{MethodMode: transport.MethodModeFixed},
			make(chan Result, 1),
			"FILE",
			"test-ext",
			func(ext string) string {
				return urlutil.AddExtension(server.URL+"/asset", ext)
			},
			buildCalibrationURL,
		)
	}()

	cancel()

	select {
	case ok := <-done:
		if ok {
			t.Fatal("processExtensions returned true after cancellation, want false")
		}
	case <-time.After(time.Second):
		t.Fatal("processExtensions did not return after context cancellation")
	}
}

func TestProcessExtensionsNormalizesExtensionURLs(t *testing.T) {
	var (
		mu             sync.Mutex
		requestedPaths []string
	)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/asset.txt" {
			mu.Lock()
			requestedPaths = append(requestedPaths, r.URL.Path)
			mu.Unlock()
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "baseline")
	}))
	defer server.Close()

	c := newExtensionTestCore([]string{"txt", ".txt"})
	buildCalibrationURL := extensionCalibrationURLBuilder(server.URL)
	buildSingleExtensionCalibration(t, c, buildCalibrationURL, "txt")

	results := make(chan Result, len(c.Exts))
	ok := c.processExtensions(
		&transport.RequestOptions{MethodMode: transport.MethodModeFixed},
		results,
		"FILE",
		"test-ext",
		func(ext string) string {
			return urlutil.AddExtension(server.URL+"/asset?download=1#top", ext)
		},
		buildCalibrationURL,
	)

	if !ok {
		t.Fatal("processExtensions returned false, want true")
	}

	if got := drainResults(results); len(got) != 0 {
		t.Fatalf("results = %#v, want none", got)
	}

	mu.Lock()
	paths := append([]string(nil), requestedPaths...)
	mu.Unlock()

	wantPaths := []string{"/asset.txt", "/asset.txt"}
	if !reflect.DeepEqual(paths, wantPaths) {
		t.Fatalf("requested paths = %#v, want %#v", paths, wantPaths)
	}
}

func newExtensionTestCore(exts []string) *Core {
	return &Core{
		Client:     transport.New(&fasthttp.Client{}),
		Ctx:        context.Background(),
		Method:     transport.MethodGET,
		MethodMode: transport.MethodModeFixed,
		UserAgents: []string{"godirb-test"},
		Exts:       exts,
	}
}

func extensionCalibrationURLBuilder(baseURL string) func(ext string) string {
	return func(ext string) string {
		return urlutil.AddExtension(
			urlutil.JoinPath(baseURL, ExtPlaceholder),
			ext,
		)
	}
}

func buildExtensionCalibrations(
	t *testing.T,
	c *Core,
	buildCalibrationURL func(ext string) string,
) {
	t.Helper()

	for _, ext := range c.Exts {
		buildSingleExtensionCalibration(t, c, buildCalibrationURL, ext)
	}
}

func buildSingleExtensionCalibration(
	t *testing.T,
	c *Core,
	buildCalibrationURL func(ext string) string,
	ext string,
) {
	t.Helper()

	baseURL := buildCalibrationURL(ext)

	err := calibration.Build(c.Client, calibration.Options{
		BaseURL:     baseURL,
		Placeholder: ExtPlaceholder,
		Tries:       3,
		UserAgents:  c.UserAgents,
	})
	if err != nil {
		t.Fatalf("build calibration for extension %q: %v", ext, err)
	}

	if _, ok := calibration.Get(baseURL, ExtPlaceholder); !ok {
		t.Fatalf("calibration missing after build for extension %q", ext)
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
