package core

import (
	"io"

	"github.com/alikWu/go-tomcat/servlet"
)

type servletResponseFacade struct {
	response servlet.ServletResponse
}

func NewServletResponseFacade(response servlet.ServletResponse) servlet.ServletResponse {
	return &servletResponseFacade{response: response}
}

func (srf servletResponseFacade) GetWriter() io.Writer {
	return srf.response.GetWriter()
}

func (srf servletResponseFacade) GetContentType() string {
	return srf.response.GetContentType()
}

func (srf servletResponseFacade) SetContentType(s string) {
	srf.response.SetContentType(s)
}

func (srf servletResponseFacade) GetContentLength() int64 {
	return srf.response.GetContentLength()
}

func (srf servletResponseFacade) SetContentLength(l int64) {
	srf.response.SetContentLength(l)
}
