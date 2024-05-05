package http

//头行格式
type HttpRequestLine struct {
	INITIAL_METHOD_SIZE   int64
	INITIAL_URI_SIZE      int64
	INITIAL_PROTOCOL_SIZE int64
	MAX_METHOD_SIZE       int64
	MAX_URI_SIZE          int64
	MAX_PROTOCOL_SIZE     int64
	//method uri protocol， eg: GET /hello.txt HTTP/1.1
	Method []byte
	//exclude the position
	MethodEnd   int64
	Uri         []byte
	UriEnd      int64
	Protocol    []byte
	ProtocolEnd int64
}

func NewHttpRequestLine() *HttpRequestLine {
	hql := &HttpRequestLine{
		INITIAL_METHOD_SIZE:   8,
		INITIAL_URI_SIZE:      128,
		INITIAL_PROTOCOL_SIZE: 8,
		MAX_METHOD_SIZE:       32,
		MAX_URI_SIZE:          2048,
		MAX_PROTOCOL_SIZE:     32,
	}
	hql.Method = make([]byte, hql.MAX_METHOD_SIZE)
	hql.Uri = make([]byte, hql.MAX_URI_SIZE)
	hql.Protocol = make([]byte, hql.MAX_PROTOCOL_SIZE)
	return hql
}

func (hrl *HttpRequestLine) Recycle() {
	hrl.MethodEnd = 0
	hrl.UriEnd = 0
	hrl.ProtocolEnd = 0
}
