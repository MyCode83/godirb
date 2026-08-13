package transport

import "time"

type Response struct {
	URL    string
	Method Method

	StatusCode    int
	ContentLenght int
	Lenght        int

	ContentType string
	Location    string

	Body []byte

	Duration time.Duration
}
