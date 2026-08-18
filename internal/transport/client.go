package transport

import (
	"time"

	"github.com/valyala/fasthttp"
)

type Client struct {
	raw *fasthttp.Client
}

func New(raw *fasthttp.Client) *Client {
	return &Client{
		raw: raw,
	}
}

func (c *Client) Do(opts *RequestOptions) (Response, error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()

	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(opts.URL)
	req.Header.SetUserAgent(opts.UserAgent)

	method := opts.methodForRequest()
	req.Header.SetMethod(method.String())

	if opts.Headers != nil {
		err := applyHeaders(req, opts.Headers)
		if err != nil {
			return Response{}, err
		}
	}

	start := time.Now()

	err := c.raw.Do(req, resp)
	finished := time.Since(start)

	if err != nil {
		return Response{}, err
	}

	body := append([]byte(nil), resp.Body()...)
	lenght := len(body)

	headers := make(map[string]string, resp.Header.Len())

	for key, value := range resp.Header.All() {
		headers[string(key)] = string(value)
	}

	title := ""
	if CanHaveTitleTag(string(resp.Header.ContentType())) {
		title = ExtractTitle(body)
	}

	return Response{
		URL:           opts.URL,
		Method:        method,
		StatusCode:    resp.StatusCode(),
		ContentLenght: resp.Header.ContentLength(),
		Lenght:        lenght,
		ContentType:   string(resp.Header.ContentType()),
		Location:      string(resp.Header.Peek("Location")),

		Body: body,
		Headers: headers,
		Title: title,

		Duration: finished,
	}, err
}
