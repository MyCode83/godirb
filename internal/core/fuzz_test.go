package core

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync"
	"testing"

	"github.com/MyCode83/godirb/internal/calibration"
	"github.com/MyCode83/godirb/internal/transport"
	"github.com/valyala/fasthttp"
)

func TestRunFuzzProcessesExtensionsBeforeFilteringBaseCalibration(t *testing.T) {
	var requestedPaths []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestedPaths = append(requestedPaths, r.URL.Path)

		switch r.URL.Path {
		case "/asset":
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "missing")
		case "/asset.php":
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "found")
		case "/asset.js":
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
		Placeholder: "FUZZ",
		UserAgents:  []string{"godirb-test"},
		Exts:        []string{"php", "js"},
		Calibration: &calibration.Calibration{
			Status:    http.StatusNotFound,
			Length:    len("missing"),
			Tolerance: 0,
			Stable:    true,
		},
		Limiter: make(chan struct{}, 1),
		WG:      &sync.WaitGroup{},
		WL:      []string{"asset"},
	}

	got := collectResults(c.RunFuzz(server.URL + "/FUZZ"))
	want := []Result{{
		Prefix: prefix,
		URL:    server.URL + "/asset.php",
		Size:   len("found"),
		Status: http.StatusOK,
	}}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("results = %#v, want %#v", got, want)
	}

	if wantPaths := []string{"/asset", "/asset.php", "/asset.js"}; !reflect.DeepEqual(requestedPaths, wantPaths) {
		t.Fatalf("requested paths = %#v, want %#v", requestedPaths, wantPaths)
	}
}
