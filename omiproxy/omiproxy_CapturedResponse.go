package proxy

import (
	"bytes"
	"net/url"
)

type CapturedResponse struct {
	StatusCode  int
	Body        bytes.Buffer
	Error       error
	OriginalURL *url.URL
	TargetURL   *url.URL
}
