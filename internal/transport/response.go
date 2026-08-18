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
	Title 		string

	Body []byte
	Headers map[string]string

	Duration time.Duration
}
