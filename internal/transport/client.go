package transport

import "github.com/valyala/fasthttp"

type Client struct {
	raw *fasthttp.Client
}

func New(raw *fasthttp.Client) *Client {
	return &Client{
		raw: raw,
	}
}

func (c *Client) Do(opts RequestOptions) (Response, error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()

	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(opts.URL)

	if opts.MethodMode == MethodModeSwitch {
		opts.Method.Toggle()
	}
	req.Header.SetMethod(opts.Method.String())

	if opts.Headers != nil {
		err := applyHeaders(req, opts.Headers)
		if err != nil {
			return Response{}, err
		}
	}

	err := c.raw.Do(req, resp)
	if err != nil {
		return Response{}, err
	}

	body := resp.Body()
	lenght := len(body)

	return Response{
		URL:           opts.URL,
		Method:        opts.Method,
		StatusCode:    resp.StatusCode(),
		ContentLenght: resp.Header.ContentLength(),
		Lenght:        lenght,
		ContentType:   string(resp.Header.ContentType()),
		Location:      string(resp.Header.Peek("Location")),

		Body: body,
	}, err
}
