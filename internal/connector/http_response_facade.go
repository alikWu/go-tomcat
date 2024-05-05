package connector

import (
	"io"

	"github.com/alikWu/go-tomcat/internal"
	"github.com/alikWu/go-tomcat/internal/cookie"
)

type HttpResponseFacade struct {
	response internal.HttpServletResponse
}

func NewHttpResponseFacade(response internal.HttpServletResponse) internal.HttpServletResponse {
	return &HttpResponseFacade{response: response}
}

func (h HttpResponseFacade) SetStatus(status int32) {
	h.response.SetStatus(status)
}

func (h HttpResponseFacade) GetStatus() int32 {
	return h.response.GetStatus()
}

func (h HttpResponseFacade) GetWriter() io.Writer {
	return h.response.GetWriter()
}

func (h HttpResponseFacade) GetContentType() string {
	return h.response.GetContentType()
}

func (h HttpResponseFacade) GetContentLength() int64 {
	return h.response.GetContentLength()
}

func (h HttpResponseFacade) SetContentLength(l int64) {
	h.response.SetContentLength(l)
}

func (h HttpResponseFacade) SetContentType(s string) {
	h.response.SetContentType(s)
}

func (h HttpResponseFacade) AddCookie(c *cookie.Cookie) {
	h.response.AddCookie(c)
}

func (h HttpResponseFacade) SetHeader(name, value string) {
	h.response.SetHeader(name, value)
}

func (h HttpResponseFacade) GetHeader(name string) string {
	return h.response.GetHeader(name)
}

func (h HttpResponseFacade) GetHeaderNames() []string {
	return h.response.GetHeaderNames()
}
