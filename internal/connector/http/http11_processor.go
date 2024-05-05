package http

import (
	"context"
	"net"
	"time"

	"go-tomcat/internal/connector"
	"go-tomcat/internal/logger"
)

type Http11Processor struct {
	conn              net.Conn
	keepAlive         bool  //默认长连接
	connectionTimeout int32 // seconds

	ctx context.Context
	hc  *HttpConnectorImpl
}

func NewHttp11Processor(ctx context.Context, conn net.Conn, hc *HttpConnectorImpl) *Http11Processor {
	return &Http11Processor{
		conn:              conn,
		keepAlive:         true,
		ctx:               ctx,
		hc:                hc,
		connectionTimeout: hc.GetConnectionTimeout(),
	}
}

func (h *Http11Processor) handle() error {
	logger.LogInfo(h.ctx, "http11Processor start")
	for h.keepAlive {
		//set connection timeout
		deadline := time.Now().Add(time.Duration(h.connectionTimeout) * time.Second)
		if err := h.conn.SetReadDeadline(deadline); err != nil {
			logger.LogError(h.ctx, "Http11Processor#handle SetReadDeadline err:", err)
			return err
		}

		request := NewHttpRequest(h.ctx, h.conn, h.hc.GetSessionKeeper())
		response := NewHttpResponse(h.ctx, h.conn)
		response.SetRequest(request)
		request.SetResponse(response)
		response.SetContentType("text/html; charset=utf-8")

		var err error
		if err = request.Parse(); err != nil {
			logger.LogError(h.ctx, "Http11Processor parse request err:", err)
			return err
		}

		requestFacade := connector.NewHttpRequestFacade(request)
		responseFacade := connector.NewHttpResponseFacade(response)
		if err = h.hc.GetContainer().Invoke(h.ctx, requestFacade, responseFacade); err != nil {
			logger.LogError(h.ctx, "http11Processor call container.Invoke err", err)
		}

		if err = response.FinishResponse(); err != nil {
			return err
		}

		if CONNECTION_CLOSE == request.GetHeader(string(CONNECTION_NAME)) || CONNECTION_CLOSE == response.GetHeader(string(CONNECTION_NAME)) {
			h.keepAlive = false
		}
	}
	logger.LogInfo(h.ctx, "http11Processor end")
	return nil
}
