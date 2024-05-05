package http

import (
	"context"
	"fmt"
	"net"

	"github.com/alikWu/go-tomcat/internal/logger"
)

type SocketProcessorRecycler interface {
	Recycle(hp *SocketProcessor) error
}

type SocketProcessor struct {
	protocol string
	connCh   chan net.Conn
	ctx      context.Context
	hc       *HttpConnectorImpl

	hpc SocketProcessorRecycler
}

func NewSocketProcessor(hpc SocketProcessorRecycler, hc *HttpConnectorImpl) *SocketProcessor {
	c := make(chan net.Conn, 1)
	return &SocketProcessor{
		connCh:   c,
		hpc:      hpc,
		hc:       hc,
		protocol: hc.GetProtocol(),
	}
}

func (hp *SocketProcessor) Process(ctx context.Context, conn net.Conn) {
	hp.ctx = ctx
	hp.connCh <- conn
}

func (hp *SocketProcessor) Start() {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				logger.LogWarn(hp.ctx, fmt.Sprintf("SocketProcessor#Start recover get err:%v", err))
			}
		}()

		for true {
			conn := <-hp.connCh
			if conn == nil {
				logger.LogWarn(hp.ctx, "SocketProcessor#Start get a nil conn")
				continue
			}
			logger.LogInfof(hp.ctx, "SocketProcessor start to handle new conn remoteAddr=%v", conn.RemoteAddr())

			//now only support http 1.1
			if hp.protocol != "HTTP/1.1" {
				panic(fmt.Sprintf("unsupported protocol=%s", hp.protocol))
			}

			http11Processor := NewHttp11Processor(hp.ctx, conn, hp.hc)
			if err := http11Processor.handle(); err != nil {
				logger.LogError(hp.ctx, "SocketProcessor#Start http11Processor.handle err:", err)
			}

			if err := conn.Close(); err != nil {
				logger.LogError(hp.ctx, "SocketProcessor#Start close conn err:", err)
			}

			if err := hp.hpc.Recycle(hp); err != nil {
				logger.LogWarn(hp.ctx, "SocketProcessor#Start recycle fail")
			}
			logger.LogInfof(hp.ctx, "SocketProcessor end the conn")
		}
	}()
}
