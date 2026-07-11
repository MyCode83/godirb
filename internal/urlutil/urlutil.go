package urlutil

import "strings"

func JoinPath(fullURL, newPath string) (string, error) {
	base := DropQueryAndFragment(fullURL)
	if base == "" {
		path := strings.TrimLeft(newPath, "/")
		if path == "" {
			return "/", nil
		}
		return DropQueryAndFragment("/" + path), nil
	}

	if newPath == "" {
		return base, nil
	}

	path := strings.TrimLeft(newPath, "/")
	if strings.HasSuffix(base, "/") {
		return DropQueryAndFragment(base + path), nil
	}

	return DropQueryAndFragment(base + "/" + path), nil
}

func AddExtension(fullURL, ext string) (string, error) {
	base := DropQueryAndFragment(fullURL)
	cleanExt := strings.TrimSpace(ext)
	if cleanExt == "" {
		return base, nil
	}

	if !strings.HasPrefix(cleanExt, ".") {
		cleanExt = "." + cleanExt
	}

	base = strings.TrimSuffix(base, "/")
	return DropQueryAndFragment(base + cleanExt), nil
}

func DropQueryAndFragment(raw string) string {
	for i, char := range raw {
		if char == '?' || char == '#' {
			return raw[:i]
		}
	}

	return raw
}
