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
	var requestedPaths []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestedPaths = append(requestedPaths, r.URL.Path)

		switch r.URL.Path {
		case "/asset":
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "missing")
		case "/asset.txt":
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "found")
		default:
			t.Fatalf("unexpected request path %q", r.URL.Path)
		}
	}))
	defer server.Close()

	c := &Core{
		Client:     transport.New(&fasthttp.Client{}),
		Ctx:        context.Background(),
		Method:     transport.MethodGET,
		MethodMode: transport.MethodModeFixed,
		UserAgents: []string{"godirb-test"},
		Exts:       []string{"txt"},
		Calibration: &calibration.Calibration{
			Status:    http.StatusNotFound,
			Length:    len("missing"),
			Tolerance: 0,
			Stable:    true,
		},
		Limiter:     make(chan struct{}, 1),
		DirsChan:    make(chan string, 1),
		WG:          &sync.WaitGroup{},
		WL:          []string{"asset"},
		VisitedDirs: map[string]bool{},
	}

	got := collectResults(c.RunDir(server.URL))
	want := []Result{{
		Prefix: "FILE",
		URL:    server.URL + "/asset.txt",
		Size:   len("found"),
		Status: http.StatusOK,
	}}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("results = %#v, want %#v", got, want)
	}

	if wantPaths := []string{"/asset", "/asset.txt"}; !reflect.DeepEqual(requestedPaths, wantPaths) {
		t.Fatalf("requested paths = %#v, want %#v", requestedPaths, wantPaths)
	}
}

func TestRunDirRecursiveJoinDoesNotCreateDoubleSlashPaths(t *testing.T) {
	var requestedPaths []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestedPaths = append(requestedPaths, r.URL.Path)

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

	c := &Core{
		Client:     transport.New(&fasthttp.Client{}),
		Ctx:        context.Background(),
		Method:     transport.MethodGET,
		MethodMode: transport.MethodModeFixed,
		Recursive:  true,
		UserAgents: []string{"godirb-test"},
		Calibration: &calibration.Calibration{
			Status:    http.StatusNotFound,
			Length:    len("missing"),
			Tolerance: 0,
			Stable:    true,
		},
		Limiter:     make(chan struct{}, 1),
		DirsChan:    make(chan string, 4),
		WG:          &sync.WaitGroup{},
		WL:          []string{"admin/", "child"},
		VisitedDirs: map[string]bool{},
	}

	collectResults(c.RunDir(server.URL))

	for _, path := range requestedPaths {
		if strings.Contains(path, "//") {
			t.Fatalf("requested path %q contains a double slash; all paths: %#v", path, requestedPaths)
		}
	}

	if !slices.Contains(requestedPaths, "/admin/child") {
		t.Fatalf("requested paths = %#v, want recursive child path /admin/child", requestedPaths)
	}
}

func collectResults(results <-chan Result) []Result {
	var collected []Result
	for result := range results {
		collected = append(collected, result)
	}
	return collected
}
