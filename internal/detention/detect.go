package detention

import (
	"github.com/MyCode83/godirb/internal/debug"
	"github.com/MyCode83/godirb/internal/transport"
)

func Detect(client *transport.Client, opts transport.RequestOptions) (DetentionResult, error) {
	var result DetentionResult

	// pathURL = https://example.com/admin
	// slashURL = https://example.com/admin/
	pathURL := removeFinalSlash(opts.URL)
	slashURL := addFinalSlash(opts.URL)

	debug.Printf(
		"detention start url=%q slash_url=%q method=%q mode=%q",
		pathURL,
		slashURL,
		opts.Method.String(),
		opts.MethodMode,
	)

	// Request 1: without "/"
	pathOpts := opts
	pathOpts.URL = pathURL

	pathRes, err := client.Do(&pathOpts)
	if err != nil {
		debug.Error("detention-path", err)
		return result, err
	}

	// Request 2: with "/"
	slashOpts := opts
	slashOpts.URL = slashURL

	slashRes, err := client.Do(&slashOpts)
	if err != nil {
		debug.Error("detention-slash", err)
		return result, err
	}

	result = classify(pathURL, pathRes, slashRes)

		debug.Printf(
		"detention result url=%q slash_url=%q status=%d slash_status=%d is_dir=%t is_file=%t unknown=%t",
		pathURL,
		slashURL,
		pathRes.StatusCode,
		slashRes.StatusCode,
		result.IsDir,
		result.IsFile,
		result.Unknown,
	)

	return result, err
}
