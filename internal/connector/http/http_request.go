package http

import (
	"context"
	"io"
	"net"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/alikWu/go-tomcat/internal"
	"github.com/alikWu/go-tomcat/internal/cookie"
	"github.com/alikWu/go-tomcat/internal/core"
	"github.com/alikWu/go-tomcat/internal/logger"
	"github.com/alikWu/go-tomcat/internal/util"
	"github.com/alikWu/go-tomcat/servlet"
)

const (
	queryDelimiter     = byte('?')
	parameterDelimiter = byte('&')
	semicolonDelimiter = byte(';')
	equalDelimiter     = byte('=')

	sessionExpireTime = 600 // seconds
)

//目前我们处理方法比较简单，只考虑文本类型。其实可以支持文本、二进制、压缩包，都是通过 Content-Type 指定。常见的有 application/json、application/xml 等。
//POST 可以混合，也就是 multipart/form-data 多部分，有的是文本，有的是二进制，比如图片之类的
type HttpRequest struct {
	uri     string
	headers map[string]string
	charset string

	queryStr    string
	paramParsed bool
	//value 是字符串数组，因为部分参数存在多个值与之对应，例如 options、checkbox 等。
	parameters map[string][]string

	sis         *SocketInputStream
	requestLine *HttpRequestLine
	ctx         context.Context

	cookies   []*cookie.Cookie
	session   internal.HttpSession
	sessionID string
	sk        *SessionKeeper

	response *HttpResponse
}

func NewHttpRequest(ctx context.Context, conn net.Conn, sk *SessionKeeper) *HttpRequest {
	return &HttpRequest{
		ctx:         ctx,
		headers:     make(map[string]string),
		parameters:  make(map[string][]string),
		sis:         NewSocketInputStream(ctx, conn, 2048),
		requestLine: NewHttpRequestLine(),
		sk:          sk,
	}
}

func (h *HttpRequest) SetResponse(response *HttpResponse) {
	h.response = response
}

func (h *HttpRequest) Parse() error {
	//parse request line
	if err := h.sis.ReadRequestLine(h.requestLine); err != nil {
		logger.LogError(h.ctx, "HttpRequest#Parse read request line err", err)
		return err
	}
	h.parseRequestLine()

	//parse headers
	if err := h.parseHeaders(); err != nil {
		return err
	}

	//parse session
	if len(h.sessionID) == 0 {
		_ = h.GetSession(true)
	}
	return nil
}

func (h *HttpRequest) parseRequestLine() {
	i := int64(0)
	for ; i < h.requestLine.UriEnd; i++ {
		if h.requestLine.Uri[i] == queryDelimiter {
			break
		}
	}
	if i == h.requestLine.UriEnd {
		h.uri = string(h.requestLine.Uri[:i])
	} else {
		h.uri = string(h.requestLine.Uri[:i])
		h.requestLine.UriEnd = i
		h.queryStr = string(h.requestLine.Uri[i+1:])
	}

	// /test/TestServlet;jsessionid=5AC6268DD8D4D5D1FDF5D41E9F2FD960?curAlbumID=9。
	//浏览器是在 URL 之后加上 ;jsessionid= 这个固定搭配来传递 Session，不是普通的参数格式。
	tmp := ";" + string(JSESSIONID_NAME) + "="
	if semicolonIndex := strings.Index(h.uri, tmp); semicolonIndex >= 0 {
		h.sessionID = h.uri[semicolonIndex+len(tmp):]
		h.uri = h.uri[:semicolonDelimiter]
	}
}

func (h *HttpRequest) parseParameters() error {
	if h.paramParsed {
		return nil
	}
	defer func() {
		h.paramParsed = true
	}()

	encoding := h.GetCharacterEncoding()
	if len(encoding) == 0 {
		encoding = "ISO-8859-1"
	}

	if len(h.queryStr) > 0 {
		h.doParseParameter([]byte(h.queryStr), encoding)
	}

	contentType := h.GetContentType()
	delimiterIndex := strings.Index(contentType, ";")
	if delimiterIndex > 0 {
		contentType = strings.TrimSpace(contentType[:delimiterIndex])
	} else {
		contentType = strings.TrimSpace(contentType)
	}

	if h.GetMethod() == "POST" && h.GetContentLength() > 0 && "application/x-www-form-urlencoded" == contentType {
		maxL := h.GetContentLength()
		buf := make([]byte, maxL)
		curL := int64(0)
		for curL < maxL {
			n, err := h.sis.Read(buf[curL : maxL-curL])
			if err != nil {
				logger.LogError(h.ctx, "HttpRequest#parseParameters readBatch err", err)
				return err
			}
			if n == 0 {
				logger.LogWarnf(h.ctx, "HttpRequest#parseParameters readBatch %d bytes", n)
				break
			}
			curL += int64(n)
		}
		if curL < maxL {
			logger.LogWarnf(h.ctx, "HttpRequest#parseParameters contentLength=%d, but only read %d bytes", maxL, curL)
			return errors.New("not read contentLength bytes")
		}
		h.doParseParameter(buf, encoding)
	}
	return nil
}

func (h *HttpRequest) doParseParameter(data []byte, encoding string) {
	if len(data) <= 0 {
		return
	}

	parameters := h.parameters
	var err error
	i := 0
	overwriteIndex := 0
	key := ""
	for ; i < len(data); i++ {
		c := data[i]
		switch c {
		case parameterDelimiter: //两个k/v之间的分隔符
			var v string
			if v, err = util.DecodeToUtf8(string(data[:overwriteIndex]), encoding); err != nil {
				logger.LogError(h.ctx, "HttpRequest#doParseParameter DecodeToUtf8 err", err)
			}
			parameters[key] = append(parameters[key], v)
			overwriteIndex = 0
		case equalDelimiter: //参数key=value的分隔符
			if key, err = util.DecodeToUtf8(string(data[:overwriteIndex]), encoding); err != nil {
				logger.LogError(h.ctx, "HttpRequest#doParseParameter DecodeToUtf8 err", err)
			}
			overwriteIndex = 0
		case '+': //特殊字符，空格， 历史原因：之前某一个版本的JDK，在使用URI encoding的时候会将空格编码成+
			data[overwriteIndex] = ' '
			overwriteIndex++
		case '%': //处理%NN表示的ASCII字符
			i++
			b := convertHexDigit(data[i]) << 4
			i++
			b += convertHexDigit(data[i])
			data[overwriteIndex] = b
			overwriteIndex++
		default:
			data[overwriteIndex] = data[i]
			overwriteIndex++
		}
	}
	var v string
	if v, err = util.DecodeToUtf8(string(data[:overwriteIndex]), encoding); err != nil {
		logger.LogError(h.ctx, "HttpRequest#doParseParameter DecodeToUtf8 err", err)
	}
	parameters[key] = append(parameters[key], v)
}

func convertHexDigit(b byte) byte {
	if b >= '0' && b <= '9' {
		return b - '0'
	}
	if b >= 'a' && b <= 'f' {
		return b - 'a' + 10
	}
	if b >= 'A' && b <= 'F' {
		return b - 'A' + 10
	}
	return 0
}

func (h *HttpRequest) parseHeaders() error {
	hh := NewHttpHeader()

	for true {
		if err := h.sis.ReadHeader(hh); err != nil {
			return err
		}
		if hh.NameEnd == 0 {
			return nil
		}
		name := strings.ToLower(string(hh.Name[:hh.NameEnd]))
		value := string(hh.Value[:hh.ValueEnd])
		h.headers[name] = value
		switch name {
		case string(ACCEPT_LANGUAGE_NAME):
		case string(CONTENT_TYPE_NAME):
			h.charset = parseCharset(value)
		case string(CONTENT_LENGTH_NAME):
		case string(HOST_NAME):
		case string(CONNECTION_NAME):
			if strings.ToLower(value) == CONNECTION_CLOSE {
				h.response.SetHeader(name, value)
			}
		case string(TRANSFER_ENCODING_NAME):
		case string(COOKIE_NAME):
			cookies := h.parseCookies(value)
			h.cookies = cookies
			for _, ck := range cookies {
				if ck.GetName() == string(JSESSIONID_NAME) {
					h.sessionID = ck.GetValue()
				}
			}
		}

		hh.Recycle()
	}
	return nil
}

func parseCharset(contentType string) string {
	contentType = strings.ToLower(contentType)
	charset := "charset="
	index := strings.Index(contentType, charset)
	if index < 0 {
		return ""
	}

	s := contentType[index+len(charset):]
	index = strings.Index(s, ";")
	if index < 0 {
		return strings.TrimSpace(s)
	}
	return strings.TrimSpace(s[:index])
}

//格式为: key1=value1;key2=value2
func (h *HttpRequest) parseCookies(header string) []*cookie.Cookie {
	var res []*cookie.Cookie
	for len(header) > 0 {
		indexSemicolon := strings.IndexByte(header, semicolonDelimiter)
		if indexSemicolon < 0 {
			indexSemicolon = len(header)
		} else if indexSemicolon == 0 {
			break
		}
		token := header[:indexSemicolon]
		if indexSemicolon < len(header) {
			header = header[indexSemicolon+1:]
		} else {
			header = ""
		}

		indexEqual := strings.IndexByte(token, equalDelimiter)
		if indexEqual > 0 {
			name := token[:indexEqual]
			value := token[indexEqual+1:]
			res = append(res, cookie.NewCookie(name, value))
		} else {
			logger.LogWarnf(h.ctx, "HttpRequest#parseCookies get a cookie which don't have =")
		}
	}
	return res
}

func (h *HttpRequest) GetCharacterEncoding() string {
	return h.charset
}

func (h *HttpRequest) GetContentLength() int64 {
	s := h.headers[string(CONTENT_LENGTH_NAME)]
	if len(s) == 0 {
		return 0
	}

	l, err := strconv.Atoi(s)
	if err != nil {
		logger.LogError(h.GetContext(), "HttpRequest#GetContentLength err", err)
	}
	return int64(l)
}

func (h *HttpRequest) GetContentType() string {
	return h.headers[string(CONTENT_TYPE_NAME)]
}

func (h *HttpRequest) GetInputStream() io.Reader {
	return h.sis
}

func (h *HttpRequest) GetParameter(s string) string {
	if err := h.parseParameters(); err != nil {
		return ""
	}
	values := h.parameters[s]
	if values == nil {
		return ""
	}
	return values[0]
}

func (h *HttpRequest) GetParameterMap() map[string][]string {
	if err := h.parseParameters(); err != nil {
		return make(map[string][]string)
	}
	return h.parameters
}

func (h *HttpRequest) GetParameterNames() []string {
	if err := h.parseParameters(); err != nil {
		return []string{}
	}

	res := make([]string, 0, len(h.parameters))
	for name, _ := range h.parameters {
		res = append(res, name)
	}
	return res
}

func (h *HttpRequest) GetParameterValues(s string) []string {
	if err := h.parseParameters(); err != nil {
		return []string{}
	}

	return h.parameters[s]
}

func (h *HttpRequest) GetHeader(arg string) string {
	return h.headers[arg]
}

func (h *HttpRequest) GetHeaderNames() []string {
	names := make([]string, 0, len(h.headers))
	for name, _ := range h.headers {
		names = append(names, name)
	}
	return names
}

func (h *HttpRequest) GetServletContext() servlet.ServletContext {
	return core.NewServletContext(h.ctx)
}

func (h *HttpRequest) GetCookies() []*cookie.Cookie {
	return h.cookies
}

func (h *HttpRequest) GetQueryString() string {
	return h.queryStr
}

func (h *HttpRequest) GetMethod() string {
	return string(h.requestLine.Method[:h.requestLine.MethodEnd])
}

func (h *HttpRequest) GetProtocol() string {
	return string(h.requestLine.Protocol[:h.requestLine.ProtocolEnd])
}

func (h *HttpRequest) GetServerPort() int64 {
	//TODO implement me
	panic("implement me")
}

func (h *HttpRequest) GetRemotePort() int64 {
	//TODO implement me
	panic("implement me")
}

func (h *HttpRequest) GetRemoteHost() string {
	//TODO implement me
	panic("implement me")
}

func (h *HttpRequest) GetPathInfo() string {
	return h.uri
}

func (h *HttpRequest) GetSession(create bool) internal.HttpSession {
	if h.session != nil {
		return h.session
	}

	if len(h.sessionID) > 0 {
		session := h.sk.GetSession(h.sessionID)
		if session != nil {
			return session
		}
	}

	if !create {
		return nil
	}
	h.session = h.sk.CreateSession(sessionExpireTime)
	h.sessionID = h.session.GetId()
	return h.session
}

func (h *HttpRequest) GetContext() context.Context {
	return h.ctx
}

func (h *HttpRequest) GetAttribute(name string) interface{} {
	session := h.GetSession(false)
	if session == nil {
		return nil
	}

	v := session.GetAttribute(name)
	if v != nil {
		return v
	}

	for _, c := range h.cookies {
		if strings.ToLower(c.GetName()) == strings.ToLower(name) {
			return c.GetValue()
		}
	}
	return nil
}

func (h *HttpRequest) GetAttributeNames() []string {
	session := h.GetSession(false)
	if session == nil {
		return []string{}
	}
	names := session.GetAttributeNames()
	for _, c := range h.cookies {
		names = append(names, c.GetName())
	}
	return names
}

func (h *HttpRequest) SetAttribute(name string, value interface{}) {
	ses := h.GetSession(true)
	ses.SetAttribute(name, value)
}

func (h *HttpRequest) RemoveAttribute(name string) {
	ses := h.GetSession(false)
	if ses == nil {
		return
	}
	ses.RemoveAttribute(name)
}
