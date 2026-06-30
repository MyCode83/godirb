package detection

import (
	"path/filepath"
	"strings"

	"github.com/MyCode83/godirb/internal/transport"
)

func classify(rawURL string, pathRes, slashRes transport.Response) DetectionResult {
	var result DetectionResult

	path := urlPath(rawURL)
	hasExt := filepath.Ext(strings.TrimRight(path, "/")) != ""

	status := pathRes.StatusCode
	slashStatus := slashRes.StatusCode

	pathFileCT := contentTypeLooksLikeFile(pathRes.ContentType)

	switch {
	// /admin -> /admin/ (DIR)
	case isRedirect(status) && redirectsToSlash(pathRes.Location):
		result.IsDir = true
	// /admin -> OK, /admin/ -> 404 (FILE)
	case isOK(status) && isNotFound(slashStatus):
		result.IsFile = true
	// /admin -> 404, /admin/ -> OK or forbidden (DIR)
	case isNotFound(status) && (isOK(slashStatus) || isForbidden(slashStatus)):
		result.IsDir = true
	// /style.css -> OK, /style.css/ NOT OK (FILE)
	case isOK(status) && hasExt && !isOK(slashStatus):
		result.IsFile = true
	case isOK(status) && pathFileCT:
		result.IsFile = true
	default:
		result.Unknown = true
	}

	return result
}

func isOK(status int) bool {
	return status >= 200 && status < 300
}

func isRedirect(status int) bool {
	return status == 301 || status == 302 || status == 307 || status == 308
}

func isNotFound(status int) bool {
	return status == 404
}

func isForbidden(status int) bool {
	return status == 401 || status == 403
}

func redirectsToSlash(location string) bool {
	return strings.HasSuffix(location, "/")
}

func contentTypeLooksLikeFile(contentType string) bool {
	ct := strings.ToLower(strings.TrimSpace(contentType))

	if idx := strings.Index(ct, ";"); idx != -1 {
		ct = strings.TrimSpace(ct[:idx])
	}

	switch {
	case strings.HasPrefix(ct, "image/"):
		return true
	case strings.HasPrefix(ct, "font/"):
		return true
	case strings.HasPrefix(ct, "audio/"):
		return true
	case strings.HasPrefix(ct, "video/"):
		return true
	case ct == "text/css":
		return true
	case ct == "application/javascript":
		return true
	case ct == "text/javascript":
		return true
	case ct == "application/json":
		return true
	case ct == "application/pdf":
		return true
	case ct == "application/octet-stream":
		return true
	default:
		return false
	}
}
