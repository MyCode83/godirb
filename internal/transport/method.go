package transport

import "errors"

type Method string

const (
	MethodGET    Method = "GET"
	MethodHEAD   Method = "HEAD"
	MethodSwitch Method = "SWITCH"
)

var ErrInvalidMethod = errors.New("invalid method")

func (m Method) String() string {
	return string(m)
}

func (m Method) Valid() bool {
	switch m {
	case MethodGET, MethodHEAD, MethodSwitch:
		return true
	default:
		return false
	}
}

func ParseMethod(raw string) (Method, error) {
	m := Method(raw)

	if !m.Valid() {
		return "", ErrInvalidMethod
	}
	return m, nil
}

func (m *Method) Toggle() {
	switch *m {
	case MethodGET:
		*m = MethodHEAD
	case MethodHEAD:
		*m = MethodGET
	}
}
