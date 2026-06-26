package core

import (
	"testing"

	"github.com/MyCode83/godirb/internal/transport"
)

func TestNextRequestMethodSwitchRotates(t *testing.T) {
	c := &Core{
		Method:     transport.MethodHEAD,
		MethodMode: transport.MethodModeSwitch,
	}

	want := []transport.Method{
		transport.MethodGET,
		transport.MethodHEAD,
		transport.MethodGET,
	}
	for i, expected := range want {
		if got := c.nextRequestMethod(); got != expected {
			t.Fatalf("method %d = %q, want %q", i, got, expected)
		}
	}
}

func TestNextRequestMethodFixedDoesNotRotate(t *testing.T) {
	c := &Core{
		Method:     transport.MethodHEAD,
		MethodMode: transport.MethodModeFixed,
	}

	for i := range 3 {
		if got := c.nextRequestMethod(); got != transport.MethodHEAD {
			t.Fatalf("method %d = %q, want %q", i, got, transport.MethodHEAD)
		}
	}
}
