package http

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"sync"

	"go-tomcat/internal/logger"
)

//chunk size 用16进制表示
const chunkedFormat = "%s\r\n%s\r\n"

type httpResponseWriter struct {
	hr   *HttpResponse
	conn net.Conn
	ctx  context.Context

	mutex1      sync.Mutex
	buf         []byte
	mutex2      sync.Mutex
	sendHeaders bool
}

func NewHttpResponseWriter(ctx context.Context, hr *HttpResponse, conn net.Conn) *httpResponseWriter {
	return &httpResponseWriter{
		hr:   hr,
		ctx:  ctx,
		conn: conn,
	}
}

func (hpw *httpResponseWriter) Write(p []byte) (n int, err error) {
	if !hpw.sendHeaders {
		hpw.mutex2.Lock()
		defer func() {
			hpw.mutex2.Unlock()
		}()
		if !hpw.sendHeaders {
			if err = hpw.hr.SendHeaders(); err != nil {
				logger.LogError(hpw.ctx, "HttpResponseWriter#Write SendHeaders err", err)
				return 0, err
			}
			hpw.sendHeaders = true
		}
	}

	hpw.mutex1.Lock()
	defer func() {
		hpw.mutex1.Unlock()
	}()
	if hpw.hr.GetContentLength() >= 0 {
		if len(p) == 0 {
			return len(p), nil
		}
		hpw.buf = append(hpw.buf, p...)
		if n, err = hpw.conn.Write(p); err != nil {
			logger.LogError(hpw.ctx, "HttpResponseWriter#Write write err", err)
		}
		return
	}

	data := fmt.Sprintf(chunkedFormat, strconv.FormatInt(int64(len(p)), 16), string(p))
	hpw.buf = append(hpw.buf, []byte(data)...)
	if n, err = hpw.conn.Write([]byte(data)); err != nil {
		logger.LogError(hpw.ctx, "HttpResponseWriter#Write write err", err)
	}
	return
}

func (hpw *httpResponseWriter) WriteHeaders(p []byte) (n int, err error) {
	hpw.mutex1.Lock()
	defer func() {
		hpw.mutex1.Unlock()
	}()

	hpw.buf = append(hpw.buf, p...)
	if n, err = hpw.conn.Write(p); err != nil {
		logger.LogError(hpw.ctx, "HttpResponseWriter#Write write err", err)
	}
	return
}
