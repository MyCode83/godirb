package transport

type Request struct {
	URL    string
	Method Method

	Headers []byte
}
