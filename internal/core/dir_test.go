package core

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"slices"
	"strings"
	"sync"
	"testing"

	"github.com/MyCode83/godirb/internal/calibration"
	"github.com/MyCode83/godirb/internal/transport"
	"github.com/valyala/fasthttp"
)

func TestRunDirProcessesExtensionsBeforeFilteringBaseCalibration(t *testing.T) {
	var (
		mu             sync.Mutex
		requestedPaths []string
	)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/asset":
			mu.Lock()
			requestedPaths = append(requestedPaths, r.URL.Path)
			mu.Unlock()
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "missing")
		case "/asset.txt":
			mu.Lock()
			requestedPaths = append(requestedPaths, r.URL.Path)
			mu.Unlock()
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "found")
		default:
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "missing")
		}
	}))
	defer server.Close()

	c := newDirTestCore([]string{"asset"})
	c.Exts = []string{"txt"}
	buildDirCalibration(t, c, server.URL)
	buildExtensionCalibrations(t, c, extensionCalibrationURLBuilder(server.URL))

	got := resultSummaries(collectResults(c.RunDir(server.URL)))
	want := []resultSummary{{
		Kind:   "FILE",
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

	if wantPaths := []string{"/asset", "/asset.txt"}; !reflect.DeepEqual(paths, wantPaths) {
		t.Fatalf("requested paths = %#v, want %#v", paths, wantPaths)
	}
}

func TestRunDirReusesRootDirectoryCalibration(t *testing.T) {
	var (
		mu              sync.Mutex
		runStarted      bool
		runUnknownPaths []string
	)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		if runStarted && r.URL.Path != "/asset" {
			runUnknownPaths = append(runUnknownPaths, r.URL.Path)
		}
		mu.Unlock()

		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "missing")
	}))
	defer server.Close()

	c := newDirTestCore([]string{"asset"})
	buildDirCalibration(t, c, server.URL)

	mu.Lock()
	runStarted = true
	mu.Unlock()

	if got := collectResults(c.RunDir(server.URL)); len(got) != 0 {
		t.Fatalf("results = %#v, want none", got)
	}

	mu.Lock()
	paths := append([]string(nil), runUnknownPaths...)
	mu.Unlock()

	if len(paths) != 0 {
		t.Fatalf("RunDir made unexpected calibration requests after root calibration was built: %#v", paths)
	}
}

func TestRunDirRecursiveUsesDirectoryCalibration(t *testing.T) {
	var (
		mu                        sync.Mutex
		requestedPaths            []string
		recursiveCalibrationPaths []string
	)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		requestedPaths = append(requestedPaths, r.URL.Path)
		if strings.HasPrefix(r.URL.Path, "/admin/") &&
			r.URL.Path != "/admin/child" &&
			r.URL.Path != "/admin/child/" &&
			r.URL.Path != "/admin/admin/" {
			recursiveCalibrationPaths = append(recursiveCalibrationPaths, r.URL.Path)
		}
		mu.Unlock()

		switch r.URL.Path {
		case "/admin":
			http.Redirect(w, r, "/admin/", http.StatusMovedPermanently)
		case "/admin/":
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "dir")
		case "/admin/child":
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "found")
		default:
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "missing")
		}
	}))
	defer server.Close()

	c := newDirTestCore([]string{"admin/", "child"})
	c.Recursive = true
	c.Depth = -1
	buildDirCalibration(t, c, server.URL)

	got := resultSummaries(collectResults(c.RunDir(server.URL)))
	want := resultSummary{
		Kind:   "FILE",
		URL:    server.URL + "/admin/child",
		Size:   len("found"),
		Status: http.StatusOK,
	}
	if !slices.Contains(got, want) {
		t.Fatalf("results = %#v, want to contain %#v", got, want)
	}

	mu.Lock()
	paths := append([]string(nil), requestedPaths...)
	recursivePaths := append([]string(nil), recursiveCalibrationPaths...)
	mu.Unlock()

	if !slices.Contains(paths, "/admin/child") {
		t.Fatalf("requested paths = %#v, want recursive child path /admin/child", paths)
	}
	if len(recursivePaths) == 0 {
		t.Fatalf("requested paths = %#v, want recursive calibration under /admin/", paths)
	}
}

func TestRunDirRecursiveJoinDoesNotCreateDoubleSlashPaths(t *testing.T) {
	var (
		mu             sync.Mutex
		requestedPaths []string
	)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		requestedPaths = append(requestedPaths, r.URL.Path)
		mu.Unlock()

		switch r.URL.Path {
		case "/admin":
			http.Redirect(w, r, "/admin/", http.StatusMovedPermanently)
		case "/admin/":
			w.WriteHeader(http.StatusOK)
		case "/admin/child":
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "found")
		default:
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "missing")
		}
	}))
	defer server.Close()

	c := newDirTestCore([]string{"admin/", "child"})
	c.Recursive = true
	c.Depth = -1
	buildDirCalibration(t, c, server.URL)

	collectResults(c.RunDir(server.URL))

	mu.Lock()
	paths := append([]string(nil), requestedPaths...)
	mu.Unlock()

	for _, path := range paths {
		if strings.Contains(path, "//") {
			t.Fatalf("requested path %q contains a double slash; all paths: %#v", path, paths)
		}
	}

	if !slices.Contains(paths, "/admin/child") {
		t.Fatalf("requested paths = %#v, want recursive child path /admin/child", paths)
	}
}

func TestRunDirRecursiveHonorsDepthLimit(t *testing.T) {
	var (
		mu             sync.Mutex
		requestedPaths []string
	)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		requestedPaths = append(requestedPaths, r.URL.Path)
		mu.Unlock()

		switch r.URL.Path {
		case "/admin/", "/admin/child", "/admin/child/grandchild":
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "found")
		default:
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "missing")
		}
	}))
	defer server.Close()

	c := newDirTestCore([]string{"admin/", "child", "grandchild"})
	c.Recursive = true
	c.Depth = 1
	buildDirCalibration(t, c, server.URL)

	collectResults(c.RunDir(server.URL))

	mu.Lock()
	paths := append([]string(nil), requestedPaths...)
	mu.Unlock()

	if !slices.Contains(paths, "/admin/child") {
		t.Fatalf("requested paths = %#v, want recursive child path /admin/child", paths)
	}
	if slices.Contains(paths, "/admin/child/grandchild") {
		t.Fatalf("requested paths = %#v, did not want grandchild path beyond depth limit", paths)
	}
}

func newDirTestCore(words []string) *Core {
	return &Core{
		Client:      transport.New(&fasthttp.Client{}),
		Ctx:         context.Background(),
		Method:      transport.MethodGET,
		MethodMode:  transport.MethodModeFixed,
		UserAgents:  []string{"godirb-test"},
		Limiter:     make(chan struct{}, 1),
		DirsChan:    make(chan DirTask, 8),
		WG:          &sync.WaitGroup{},
		WL:          words,
		VisitedDirs: map[string]bool{},
	}
}

func buildDirCalibration(t *testing.T, c *Core, baseURL string) {
	t.Helper()

	err := calibration.Build(c.Client, calibration.Options{
		BaseURL:     baseURL,
		Placeholder: "",
		Tries:       3,
		UserAgents:  c.UserAgents,
	})
	if err != nil {
		t.Fatalf("build directory calibration for %q: %v", baseURL, err)
	}

	if _, ok := calibration.Get(baseURL, ""); !ok {
		t.Fatalf("directory calibration missing after build for %q", baseURL)
	}
}

func collectResults(results <-chan Result) []Result {
	var collected []Result
	for result := range results {
		collected = append(collected, result)
	}
	return collected
}

type resultSummary struct {
	URL    string
	Size   int
	Status int
	Kind   string
}

func resultSummaries(results []Result) []resultSummary {
	summaries := make([]resultSummary, 0, len(results))
	for _, result := range results {
		summaries = append(summaries, resultSummary{
			URL:    result.URL,
			Size:   result.Size,
			Status: result.Status,
			Kind:   result.Kind,
		})
	}
	return summaries
}
