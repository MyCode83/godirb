package core

type Result struct {
	URL    string `json:"url" csv:"url"`
	Size   int    `json:"size" csv:"size"`
	Status int    `json:"status_code" csv:"status_code"`

	Method        string `json:"method" csv:"method"`
	ContentType   string `json:"content_type" csv:"content_type"`
	ContentLength int    `json:"content_length" csv:"content_length"`
	Location      string `json:"location" csv:"location"`
	Duration      string `json:"duration" csv:"duration"`
	Title		  string `json:"title" csv:"title"`

	Kind  string `json:"kind" csv:"kind"`
	Error string `json:"error,omitempty" csv:"error"`
}
