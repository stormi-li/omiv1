package omi

import (
	"embed"

	cert "github.com/stormi-li/omiv1/omiregister/omicert"
	web "github.com/stormi-li/omiv1/omiweb"
)

func NewWeb(embeddedSource *embed.FS) *web.Web {
	return web.NewWeb(embeddedSource)
}

func WriteDefaultCertAndKey() {
	cert.WriteDefaultCertAndKey()
}
