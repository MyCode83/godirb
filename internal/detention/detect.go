package detention

import (
	"github.com/MyCode83/godirb/internal/debug"
	"github.com/valyala/fasthttp"
	"strings"
)

func Detect(client *fasthttp.Client, baseURL string, path string, method string) (DetentionResult, error) {
	var result DetentionResult
	var err error
	debug.Printf("detention start base_url=%q path=%q method=%q", baseURL, path, method)
	url := strings.TrimSpace(baseURL)
	url = strings.TrimRight(url, "/") + "/" + strings.TrimLeft(url, "/")

	request := fasthttp.AcquireRequest()
	response := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(request)
	defer fasthttp.ReleaseResponse(response)

	// /path
	request.Reset()
	response.Reset()
	request.SetRequestURI(url)
	request.Header.SetMethod(method)
	debug.Request("detention-path", request)
	err = client.Do(request, response)
	if err != nil {
		debug.Error("detention-path", err)
		return result, err
	}
	debug.Response("detention-path", response)
	status := response.StatusCode()

	// /path/

	request.Reset()
	response.Reset()
	request.SetRequestURI(url)
	request.Header.SetMethod(method)
	debug.Request("detention-slash", request)
	err = client.Do(request, response)
	if err != nil {
		debug.Error("detention-slash", err)
		return result, err
	}
	debug.Response("detention-slash", response)
	status2 := response.StatusCode()
	if status == 200 && (status2 == 200 || status2 == 301 || status2 == 302) {
		result.IsDir = true
	} else if status == 200 && status2 >= 400 {
		result.IsFile = true
	} else {
		result.Unknown = true
	}
	debug.Printf("detention result status=%d slash_status=%d is_dir=%t is_file=%t unknown=%t",
		status, status2, result.IsDir, result.IsFile, result.Unknown)
	return result, nil

}
