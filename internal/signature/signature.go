package signature

import (
	"bytes"
	_ "embed"
	"fmt"
	"reflect"
	"sort"
	"unsafe"

	"github.com/MyCode83/godirb/internal/transport"
	"github.com/projectdiscovery/cleanhttp"
)

//go:embed rules.json
var customRules []byte

// Matcher detects known default-error HTTP response signatures.
type Matcher struct {
	matcher *cleanhttp.Matcher
}

func New() (*Matcher, error) {
	matcher, err := newEmptyCleanHTTPMatcher()
	if err != nil {
		return nil, err
	}

	if len(bytes.TrimSpace(customRules)) > 0 {
		err := matcher.AddRules(customRules)
		if err != nil {
			return nil, fmt.Errorf("load custom signatures: %w", err)
		}
	}

	return &Matcher{matcher: matcher}, nil
}

func (m *Matcher) Match(resp *transport.Response) []string {
	if m == nil || m.matcher == nil {
		return nil
	}

	matches := m.matcher.Match(cleanhttp.Response{
		StatusCode: resp.StatusCode,
		Body:       string(resp.Body),
		Headers:    resp.Headers,
		Title:      resp.Title,
		RequestURL: resp.URL,
	})

	sort.Strings(matches)
	return matches
}

func (m *Matcher) MatchDefaultError(resp *transport.Response) []string {
	return m.Match(resp)
}

func newEmptyCleanHTTPMatcher() (*cleanhttp.Matcher, error) {
	matcher := &cleanhttp.Matcher{}
	value := reflect.ValueOf(matcher).Elem()
	rulesField := value.FieldByName("rules")
	if !rulesField.IsValid() {
		return nil, fmt.Errorf("cleanhttp matcher rules field not found")
	}

	writableRules := reflect.NewAt(rulesField.Type(), unsafe.Pointer(rulesField.UnsafeAddr())).Elem()
	writableRules.Set(reflect.MakeMap(rulesField.Type()))

	return matcher, nil
}
