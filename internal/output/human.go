package output

import (
	"fmt"
	"io"

	"github.com/MyCode83/godirb/internal/core"
	"github.com/MyCode83/godirb/internal/ui"
	"github.com/charmbracelet/x/ansi"
)

type ansiStripWriter struct {
	w io.Writer
}

func (s ansiStripWriter) Write(p []byte) (int, error) {
	clean := ansi.Strip(string(p))

	if _, err := io.WriteString(s.w, clean); err != nil {
		return 0, err
	}

	return len(p), nil
}

func WriteHuman(w io.Writer, result core.Result) error {
	_, err := fmt.Fprintln(w, ui.RenderResult(result))

	return err
}
