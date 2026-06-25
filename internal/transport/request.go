package transport

type RequestOptions struct {
	URL    string
	Method Method
	MethodMode MethodMode

	Headers []string
}
