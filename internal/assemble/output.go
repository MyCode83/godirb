package assemble

import (
	"github.com/MyCode83/godirb/internal/cli"
	"github.com/MyCode83/godirb/internal/debug"
	"github.com/MyCode83/godirb/internal/output"
)

func BuildOutputConfig(cfg cli.Config) (output.Format, bool) {
	outputFormat := output.FromFlags(cfg.JSON, cfg.CSV, cfg.Quiet)
	quietOutput := cfg.Output != "" || outputFormat != output.FormatHuman
	debug.Printf("output format=%d quiet_output=%t", outputFormat, quietOutput)

	return outputFormat, quietOutput
}

func BuildOutputStream(outputFormat output.Format, cfg cli.Config) (*output.Stream, error) {
	stream, err := output.NewStream(outputFormat, cfg.Output, cfg.NoColor)
	if err != nil {
		return nil, err
	}

	return stream, nil
}
