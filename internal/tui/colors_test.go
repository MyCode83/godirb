package tui

import (
	"os"
	"strings"
	"testing"
)

func TestColorDisabledSources(t *testing.T) {
	if !colorDisabled(true, false) {
		t.Fatal("colorDisabled noColor = false, want true")
	}

	t.Setenv("NO_COLOR", "1")
	if !colorDisabled(false, false) {
		t.Fatal("colorDisabled NO_COLOR = false, want true")
	}

	if err := os.Unsetenv("NO_COLOR"); err != nil {
		t.Fatalf("Unsetenv(NO_COLOR): %v", err)
	}
	t.Setenv("GODIRB_NO_COLOR", "1")
	if !colorDisabled(false, false) {
		t.Fatal("colorDisabled GODIRB_NO_COLOR = false, want true")
	}
}

func TestConfigureColorDisablesAnsiRendering(t *testing.T) {
	ConfigureColor(true, false)

	line := successStyle.Render("200")
	if strings.Contains(line, "\x1b[") {
		t.Fatalf("rendered line contains ANSI escape sequence: %q", line)
	}
}
