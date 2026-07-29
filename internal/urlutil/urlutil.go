package urlutil

import "strings"

func JoinPath(fullURL, newPath string) string {
	base := DropQueryAndFragment(fullURL)
	if base == "" {
		path := strings.TrimLeft(newPath, "/")
		if path == "" {
			return "/"
		}
		return DropQueryAndFragment("/" + path)
	}

	if newPath == "" {
		return base
	}

	path := strings.TrimLeft(newPath, "/")
	if strings.HasSuffix(base, "/") {
		return DropQueryAndFragment(base + path)
	}

	return DropQueryAndFragment(base + "/" + path)
}

func AddExtension(fullURL, ext string) string {
	base := DropQueryAndFragment(fullURL)
	cleanExt := strings.TrimSpace(ext)
	if cleanExt == "" {
		return base
	}

	if !strings.HasPrefix(cleanExt, ".") {
		cleanExt = "." + cleanExt
	}

	base = strings.TrimSuffix(base, "/")
	return DropQueryAndFragment(base + cleanExt)
}

func DropQueryAndFragment(raw string) string {
	for i, char := range raw {
		if char == '?' || char == '#' {
			return raw[:i]
		}
	}

	return raw
}
