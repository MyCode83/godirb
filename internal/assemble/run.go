package assemble

import (
	"context"

	"github.com/MyCode83/godirb/internal/core"
	"github.com/MyCode83/godirb/internal/debug"
	"github.com/MyCode83/godirb/internal/output"
)

func Run(engine *core.Core, stream *output.Stream, baseURL string, cancel context.CancelFunc) error {
	var outputErr error
	for result := range engine.Run(baseURL) {
		debug.Printf("result prefix=%s status=%d size=%d url=%s extra=%q", result.Kind, result.Status, result.Size, result.URL, result.Error)
		if outputErr != nil {
			continue
		}
		if err := stream.Write(result); err != nil {
			debug.Error("output stream write", err)
			outputErr = err
			cancel()
		}
		continue

	}
	if stream != nil {
		if err := stream.Close(); outputErr == nil {
			outputErr = err
		}
	}
	if outputErr != nil {
		return outputErr
	}

	debug.Printf("scan finished")
	return nil
}
