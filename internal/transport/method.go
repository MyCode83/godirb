package transport

import (
	"errors"
	"strings"
)

type Method string

type MethodMode string

const (
	MethodModeFixed  MethodMode = "FIXED"
	MethodModeSwitch MethodMode = "SWITCH"
)
const (
	MethodGET    Method = "GET"
	MethodHEAD   Method = "HEAD"
)

var ErrInvalidMethod = errors.New("invalid method")

func (m Method) String() string {
	return string(m)
}

func (m Method) Valid() bool {
	switch m {
	case MethodGET, MethodHEAD:
		return true
	default:
		return false
	}
}

func ParseMethod(raw string) (Method, MethodMode, error) {
	switch strings.ToUpper(strings.TrimSpace(raw)) {
	case "GET":
		return MethodGET, MethodModeFixed, nil
	case "HEAD":
		return MethodHEAD, MethodModeFixed, nil
	case "SWITCH":
		return MethodHEAD, MethodModeSwitch, nil
	default:
		return "", "", ErrInvalidMethod
	}
}

func (m *Method) Toggle() {
	switch *m {
	case MethodGET:
		*m = MethodHEAD
	case MethodHEAD:
		*m = MethodGET
	}
}
