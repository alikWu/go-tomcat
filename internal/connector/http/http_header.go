package http

type HeaderType string

const (
	HOST_NAME              HeaderType = "host"
	CONNECTION_NAME        HeaderType = "connection"
	ACCEPT_LANGUAGE_NAME   HeaderType = "accept-language"
	CONTENT_LENGTH_NAME    HeaderType = "content-length"
	CONTENT_TYPE_NAME      HeaderType = "content-type"
	TRANSFER_ENCODING_NAME HeaderType = "Transfer-Encoding"
	COOKIE_NAME            HeaderType = "cookie"
	JSESSIONID_NAME        HeaderType = "jsessionid"
)

type HeaderValue string

const CONNECTION_CLOSE = "close"

type HttpHeader struct {
	INITIAL_NAME_SIZE  int64
	INITIAL_VALUE_SIZE int64
	MAX_NAME_SIZE      int64
	MAX_VALUE_SIZE     int64
	Name               []byte
	NameEnd            int64
	Value              []byte
	ValueEnd           int64
}

func NewHttpHeader() *HttpHeader {
	hh := &HttpHeader{
		INITIAL_NAME_SIZE:  64,
		INITIAL_VALUE_SIZE: 512,
		MAX_NAME_SIZE:      128,
		MAX_VALUE_SIZE:     1024,
	}
	hh.Name = make([]byte, hh.MAX_NAME_SIZE)
	hh.Value = make([]byte, hh.MAX_VALUE_SIZE)
	return hh
}

func (hh *HttpHeader) Recycle() {
	hh.NameEnd = 0
	hh.ValueEnd = 0
}
