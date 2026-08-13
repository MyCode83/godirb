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
	"github.com/valyala/fasthttp"
)

func TestRunFuzzProcessesExtensionsBeforeFilteringBaseCalibration(t *testing.T) {
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
		case "/asset.php":
			mu.Lock()
			requestedPaths = append(requestedPaths, r.URL.Path)
			mu.Unlock()
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "found")
		case "/asset.js":
			mu.Lock()
			requestedPaths = append(requestedPaths, r.URL.Path)
			mu.Unlock()
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "missing")
		default:
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "missing")
		}
	}))
	defer server.Close()

	c := newFuzzTestCore([]string{"asset"})
	c.Exts = []string{"php", "js"}
	c.Calibration = &calibration.Calibration{
		Status:    http.StatusNotFound,
		Length:    len("missing"),
		Tolerance: 0,
		Stable:    true,
	}
	buildFuzzExtensionCalibrations(t, c, server.URL+"/"+c.Placeholder)

	got := resultSummaries(collectResults(c.RunFuzz(server.URL + "/" + c.Placeholder)))
	want := []resultSummary{{
		Kind:   prefix,
		URL:    server.URL + "/asset.php",
		Size:   len("found"),
		Status: http.StatusOK,
	}}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("results = %#v, want %#v", got, want)
	}

	mu.Lock()
	paths := append([]string(nil), requestedPaths...)
	mu.Unlock()

	if wantPaths := []string{"/asset", "/asset.php", "/asset.js"}; !reflect.DeepEqual(paths, wantPaths) {
		t.Fatalf("requested paths = %#v, want %#v", paths, wantPaths)
	}
}

func TestRunFuzzDelaysCalibrationFilteredResponses(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "missing")
	}))
	defer server.Close()

	c := newFuzzTestCore([]string{"one", "two"})
	c.Calibration = &calibration.Calibration{
		Status:    http.StatusNotFound,
		Length:    len("missing"),
		Tolerance: 0,
		Stable:    true,
	}
	c.Delay = 20 * time.Millisecond

	start := time.Now()
	got := collectResults(c.RunFuzz(server.URL + "/" + c.Placeholder))
	elapsed := time.Since(start)

	if len(got) != 0 {
		t.Fatalf("results = %#v, want none", got)
	}

	if elapsed < 40*time.Millisecond {
		t.Fatalf("elapsed = %s, want at least %s", elapsed, 40*time.Millisecond)
	}
}

func newFuzzTestCore(words []string) *Core {
	return &Core{
		Client:      transport.New(&fasthttp.Client{}),
		Ctx:         context.Background(),
		Method:      transport.MethodGET,
		MethodMode:  transport.MethodModeFixed,
		Placeholder: "FUZZ",
		UserAgents:  []string{"godirb-test"},
		Limiter:     make(chan struct{}, 1),
		WG:          &sync.WaitGroup{},
		WL:          words,
	}
}

func buildFuzzExtensionCalibrations(t *testing.T, c *Core, baseURL string) {
	t.Helper()

	for _, ext := range c.Exts {
		baseURL := buildFuzzExtensionCalibrationURL(baseURL, c.Placeholder, ext)
		err := calibration.Build(c.Client, calibration.Options{
			BaseURL:     baseURL,
			Placeholder: ExtPlaceholder,
			Tries:       3,
			UserAgents:  c.UserAgents,
		})
		if err != nil {
			t.Fatalf("build fuzz extension calibration for %q: %v", baseURL, err)
		}

		if _, ok := calibration.Get(baseURL, ExtPlaceholder); !ok {
			t.Fatalf("fuzz extension calibration missing after build for %q", baseURL)
		}
	}
}

func buildFuzzExtensionCalibrationURL(baseURL, placeholder, ext string) string {
	if len(ext) > 0 && ext[0] != '.' {
		ext = "." + ext
	}

	urlParts := strings.Split(baseURL, placeholder)
	return urlParts[0] + ExtPlaceholder + ext + urlParts[1]
}
