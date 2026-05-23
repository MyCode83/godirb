package core

type Result struct {
	Prefix string `json:"prefix"`
	URL    string `json:"url"`
	Size   int    `json:"size"`
	Status int    `json:"status"`
	Extra  string `json:"extra,omitempty"`
}
