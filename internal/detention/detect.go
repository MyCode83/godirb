package detention

import (
	"github.com/MyCode83/godirb/internal/debug"
	"github.com/MyCode83/godirb/internal/transport"
	"strings"
)

func Detect(client *transport.Client, baseURL string, path string, method transport.Method, methodMode transport.MethodMode) (DetentionResult, error) {
	var result DetentionResult
	var err error
	debug.Printf("detention start base_url=%q path=%q method=%q", baseURL, path, method)
	url := strings.TrimSpace(baseURL)
	url = strings.TrimRight(url, "/") + "/" + strings.TrimLeft(url, "/")

	request := transport.RequestOptions{
		URL:        url,
		Method:     method,
		MethodMode: methodMode,
	}

	// /path
	response, err := client.Do(&request)
	if err != nil {
		debug.Error("detention-path", err)
		return result, err
	}
	debug.Printf("detention-path response status=%d body=%d", response.StatusCode, response.Lenght)
	status := response.StatusCode

	// /path/

	response, err = client.Do(&request)
	if err != nil {
		debug.Error("detention-slash", err)
		return result, err
	}
	debug.Printf("detention-slash response status=%d body=%d", response.StatusCode, response.Lenght)
	status2 := response.StatusCode
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
