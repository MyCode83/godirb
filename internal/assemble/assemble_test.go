package assemble

import (
	"context"
	"sync"
	"testing"

	"github.com/MyCode83/godirb/internal/cli"
	"github.com/MyCode83/godirb/internal/core"
	"github.com/MyCode83/godirb/internal/output"
	"github.com/MyCode83/godirb/internal/transport"
)

func TestBuildAuth(t *testing.T) {
	got := BuildAuth(cli.Config{Username: "user", Password: "pass"})
	if got != "Basic dXNlcjpwYXNz" {
		t.Fatalf("BuildAuth() = %q", got)
	}
}

func TestBuildMethod(t *testing.T) {
	method, mode, err := BuildMethod(cli.Config{Method: "switch"})
	if err != nil {
		t.Fatalf("BuildMethod() returned error: %v", err)
	}
	if method != transport.MethodHEAD || mode != transport.MethodModeSwitch {
		t.Fatalf("BuildMethod() = %q, %q", method, mode)
	}
}

func TestBuildOutputConfig(t *testing.T) {
	format, quiet := BuildOutputConfig(cli.Config{JSON: true})
	if format != output.FormatJSON || !quiet {
		t.Fatalf("BuildOutputConfig() = %d, %t", format, quiet)
	}
}

func TestBuildEngine(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var wg sync.WaitGroup

	engine := BuildEngine(ctx, cancel, cli.Config{Threads: 2}, core.ModeDir, nil,
		transport.MethodGET, transport.MethodModeFixed, "auth", []string{"admin"}, &wg)

	if engine.Mode != core.ModeDir || engine.AuthHeader != "auth" {
		t.Fatalf("BuildEngine() did not copy core configuration")
	}
	if cap(engine.Limiter) != 2 || cap(engine.DirsChan) != 100 {
		t.Fatalf("BuildEngine() created unexpected channels")
	}
}
