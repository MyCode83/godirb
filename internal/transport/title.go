// Portions of this file are adapted from ProjectDiscovery httpx:
// https://github.com/projectdiscovery/httpx/blob/dev/common/httpx/title.go
//
// Copyright (c) 2021 ProjectDiscovery, Inc.
// Licensed under the MIT License.
// See THIRD_PARTY_NOTICES.md for the complete license text.
package transport

import (
	"bytes"
	"fmt"
	"io"
	"mime"
	"regexp"
	"slices"
	"strings"

	"golang.org/x/net/html"
)

var (
	titleCutset   = "\n\t\v\f\r"
	titleRegexp   = regexp.MustCompile(`(?im)<\s*title.*>(.*?)<\s*/\s*title>`)
	titleReplacer = strings.NewReplacer(
		"\n", "",
		"\t", "",
		"\v", "",
		"\f", "",
		"\r", "",
	)

	supportedTitleMimeTypes = []string{
		"text/html",
		"application/xhtml+xml",
		"application/xml",
		"application/rss+xml",
		"application/atom+xml",
		"application/vnd.wap.xhtml+xml",
	}
)

// ExtractTitle extracts the title from an HTML response body.
func ExtractTitle(body []byte) (title string) {
	titleDOM, err := getTitleWithDOM(body)
	if err != nil {
		title = titleRegexp.FindString(string(body))
	} else {
		title = renderNode(titleDOM)
	}

	title = html.UnescapeString(trimTitleTags(title))
	title = strings.TrimSpace(strings.Trim(title, titleCutset))
	title = titleReplacer.Replace(title)

	return title
}

// CanHaveTitleTag reports whether the content type may contain a title tag.
func CanHaveTitleTag(contentType string) bool {
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		mediaType, _, _ = strings.Cut(contentType, ";")
	}

	mediaType = strings.ToLower(strings.TrimSpace(mediaType))

	return slices.Contains(supportedTitleMimeTypes, mediaType)
}

func getTitleWithDOM(body []byte) (*html.Node, error) {
	var title *html.Node

	var crawler func(*html.Node)
	crawler = func(node *html.Node) {
		if node.Type == html.ElementNode &&
			strings.EqualFold(node.Data, "title") {
			title = node
			return
		}

		for child := node.FirstChild; child != nil && title == nil; child = child.NextSibling {
			crawler(child)
		}
	}

	document, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	crawler(document)

	if title == nil {
		return nil, fmt.Errorf("title not found")
	}

	return title, nil
}

func renderNode(node *html.Node) string {
	var buffer bytes.Buffer
	_ = html.Render(io.Writer(&buffer), node)

	return buffer.String()
}

func trimTitleTags(title string) string {
	titleBegin := strings.Index(title, ">")
	titleEnd := strings.Index(title, "</")

	if titleEnd < 0 || titleBegin < 0 || titleEnd <= titleBegin {
		return title
	}

	return title[titleBegin+1 : titleEnd]
}