package http

import (
	"context"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"

	"go-tomcat/internal"
	"go-tomcat/internal/cookie"
	"go-tomcat/internal/logger"
)

const lineBreak = "\r\n"

type HttpResponse struct {
	request       internal.HttpServletRequest
	writer        *httpResponseWriter
	contentType   string
	contentLength int64
	protocol      string
	headers       map[string]string
	ctx           context.Context

	status int32
	//charset           string
	//characterEncoding string

	cookies []*cookie.Cookie
}

func NewHttpResponse(ctx context.Context, conn net.Conn) *HttpResponse {
	h := &HttpResponse{
		contentLength: -1,
		protocol:      "HTTP/1.1",
		headers:       make(map[string]string),
		ctx:           ctx,
	}
	writer := NewHttpResponseWriter(ctx, h, conn)
	h.writer = writer
	return h
}

func (h *HttpResponse) SetRequest(request internal.HttpServletRequest) {
	h.request = request
}

func (h *HttpResponse) GetResponseWriter() *httpResponseWriter {
	return h.writer
}

func (h *HttpResponse) GetWriter() io.Writer {
	return h.writer
}

func (h *HttpResponse) GetContentType() string {
	return h.contentType
}

func (h *HttpResponse) SetContentLength(l int64) {
	h.contentLength = l
}

func (h *HttpResponse) GetContentLength() int64 {
	return h.contentLength
}

func (h *HttpResponse) SetContentType(s string) {
	h.contentType = s
}

func (h *HttpResponse) AddCookie(c *cookie.Cookie) {
	h.cookies = append(h.cookies, c)
}

func (h *HttpResponse) SetHeader(name, value string) {
	h.headers[name] = value
	if strings.ToLower(name) == string(CONTENT_LENGTH_NAME) {
		l, err := strconv.Atoi(value)
		if err != nil {
			logger.LogError(h.ctx, fmt.Sprintf("HttpResponse#SetHeader atoi err, value=%s", value), err)
			return
		}
		h.SetContentLength(int64(l))
	}

	if strings.ToLower(name) == string(CONTENT_TYPE_NAME) {
		h.SetContentType(value)
	}
}

func (h *HttpResponse) GetHeader(name string) string {
	return h.headers[name]
}

func (h *HttpResponse) GetHeaderNames() []string {
	var res []string
	for name, _ := range h.headers {
		res = append(res, name)
	}
	return res
}

func (h *HttpResponse) SendHeaders() error {
	var sb strings.Builder
	if h.status == 0 || h.status == int32(internal.SC_OK) {
		sb = h.buildOKHeader()
	} else {
		sb = h.buildBriefHeader()
	}

	if _, err := h.GetResponseWriter().WriteHeaders([]byte(sb.String())); err != nil {
		logger.LogError(h.ctx, "HttpResponse#SendHeaders WriteHeaders err", err)
		return err
	}
	return nil
}

func (h *HttpResponse) buildBriefHeader() strings.Builder {
	var sb strings.Builder
	//状态行
	sb.WriteString(h.protocol + " ")
	sb.WriteString(strconv.Itoa(int(h.status)) + " ")
	sb.WriteString(internal.StatusMessageMap[internal.StatusCode(h.status)] + lineBreak)

	sb.WriteString(string(CONTENT_LENGTH_NAME) + ": 0" + lineBreak)
	sb.WriteString(lineBreak)
	return sb
}

func (h *HttpResponse) buildOKHeader() strings.Builder {
	var sb strings.Builder
	sb.Grow(1024)

	//状态行
	sb.WriteString(h.protocol + " ")
	status := int32(internal.SC_OK)
	sb.WriteString(strconv.Itoa(int(status)) + " ")
	sb.WriteString(internal.StatusMessageMap[internal.StatusCode(status)] + lineBreak)

	//头信息
	if len(h.GetContentType()) > 0 {
		sb.WriteString(string(CONTENT_TYPE_NAME) + ": " + h.GetContentType() + lineBreak)
	}
	if h.GetContentLength() >= 0 {
		sb.WriteString(fmt.Sprintf("%s: %d%s", CONTENT_LENGTH_NAME, h.GetContentLength(), lineBreak))
	} else {
		sb.WriteString(string(TRANSFER_ENCODING_NAME) + ": chunked" + lineBreak)
	}

	for name, value := range h.headers {
		sb.WriteString(name + ": " + value + lineBreak)
	}

	//cookies
	session := h.request.GetSession(true)
	if session != nil {
		ck := cookie.NewCookie(string(JSESSIONID_NAME), session.GetId())
		ck.SetMaxAge(-1)
		h.cookies = append(h.cookies, ck)
	}
	for _, ck := range h.cookies {
		sb.WriteString(ck.ToCookieHeaderName() + ": ")
		sb.WriteString(ck.ToCookieHeaderValue() + lineBreak)
	}

	//输出空行
	sb.WriteString(lineBreak)
	return sb
}

func (h *HttpResponse) FinishResponse() error {
	_, err := h.writer.Write([]byte{})
	return err
}

func (h *HttpResponse) SetStatus(status int32) {
	h.status = status
}

func (h *HttpResponse) GetStatus() int32 {
	return h.status
}
