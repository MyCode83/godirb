package detention

import (
	"github.com/valyala/fasthttp"
	"strings"
)
func Detect(client *fasthttp.Client, baseURL string, path string, method string) (DetentionResult, error) {
	var result DetentionResult
	var err error
	url := strings.TrimSpace(baseURL)
	url = strings.TrimRight(url, "/") + "/" + strings.TrimLeft(url,"/")

	request := fasthttp.AcquireRequest()
	response := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(request)
	defer fasthttp.ReleaseResponse(response)

	// /path
	request.Reset()
	response.Reset()
	request.SetRequestURI(url)
	request.Header.SetMethod(method)
	err = client.Do(request, response)
	if err != nil {
		return result, err
	}
	status := response.StatusCode()

	// /path/

	request.Reset()
	response.Reset()
	request.SetRequestURI(url)
	request.Header.SetMethod(method)
	err = client.Do(request, response)
	if err != nil {
		return result, err
	}
	status2 := response.StatusCode()
	if status == 200 && (status2 == 200 || status2 == 301 || status2 == 302) {
		result.IsDir = true
	} else if status == 200 && status2 >= 400 {
		result.IsFile = true
	} else {
		result.Unknown = true
	}
	return result, nil

}