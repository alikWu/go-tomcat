package servlet

import (
	"io"
)

type ServletResponse interface {
	GetWriter() io.Writer
	GetContentType() string
	SetContentType(s string)
	GetContentLength() int64
	SetContentLength(l int64)
}
