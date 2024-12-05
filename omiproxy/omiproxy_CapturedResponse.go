package proxy

import (
	"net/url"
)

type CapturedResponse struct {
	StatusCode int
	Body       []byte
	Error      error
	TargetURL  *url.URL
}
