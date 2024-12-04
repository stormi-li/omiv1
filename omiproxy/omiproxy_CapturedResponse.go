package proxy

import (
	"bytes"
	"net/url"
)

type CapturedResponse struct {
	StatusCode int
	Body       bytes.Buffer
	Error      error
	TargetURL  *url.URL
}
