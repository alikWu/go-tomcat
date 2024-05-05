package http

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/google/uuid"
	"go-tomcat/internal"
	"go-tomcat/internal/logger"
)

type HttpConnectorImpl struct {
	port              string
	protocol          string
	connectionTimeout int32 // seconds

	minProcessors int32
	maxProcessors int32 //default 200
	curProcessors int32

	mutex sync.Mutex
	hps   []*SocketProcessor

	sk *SessionKeeper

	ctx context.Context
	sc  internal.Context
}

func NewHttpConnector(port string) *HttpConnectorImpl {
	return &HttpConnectorImpl{
		port:              port,
		maxProcessors:     200,
		protocol:          "HTTP/1.1",
		connectionTimeout: 20,
		sk:                NewSessionKeeper(),
	}
}

func (hc *HttpConnectorImpl) SetMaxConnections(maxConnections int32) {
	hc.maxProcessors = maxConnections
}

func (hc *HttpConnectorImpl) SetConnectionTimeout(connectionTimeout int32) {
	hc.connectionTimeout = connectionTimeout
}

func (hc *HttpConnectorImpl) GetConnectionTimeout() int32 {
	return hc.connectionTimeout
}

func (hc *HttpConnectorImpl) SetContainer(container internal.Context) {
	hc.sc = container
}

func (hc *HttpConnectorImpl) GetContainer() internal.Context {
	return hc.sc
}

func (hc *HttpConnectorImpl) ListenConnect() error {
	ctx := context.Background()
	logger.LogInfo(ctx, "HttpConnector start")
	listen, err := net.Listen("tcp", hc.port)
	if err != nil {
		logger.LogError(ctx, "HttpConnectorImpl#Await net.Listen err", err)
		return err
	}

	for {
		conn, err := listen.Accept()
		if err != nil {
			logger.LogError(ctx, "HttpConnectorImpl#Await accept err", err)
			return err
		}

		fmt.Println("HttpConnectorImpl new conn, remotePort:", conn.RemoteAddr())
		curCtx := context.WithValue(ctx, logger.TraceID, uuid.New())
		processor := hc.getProcessor()

		//if run out of processor, deprecate new requests
		if processor == nil {
			logger.LogWarn(ctx, "HttpConnectorImpl#ListenConnect don't have enough processor, deprecate the conn")
			response := NewHttpResponse(curCtx, conn)
			response.SetStatus(int32(internal.SC_SERVICE_UNAVAILABLE))
			if err = response.FinishResponse(); err != nil {
				logger.LogError(ctx, "HttpConnectorImpl#ListenConnect FinishResponse err", err)
			}
			_ = conn.Close()
			continue
		}
		processor.Process(curCtx, conn)
	}
}

func (hc *HttpConnectorImpl) getProcessor() *SocketProcessor {
	hc.mutex.Lock()
	defer func() {
		hc.mutex.Unlock()
	}()
	if len(hc.hps) == 0 {
		if hc.curProcessors < hc.maxProcessors {
			hc.curProcessors++
			processor := NewSocketProcessor(hc, hc)
			processor.Start()
			return processor
		}
		return nil
	}
	res := hc.hps[len(hc.hps)-1]
	hc.hps = hc.hps[:len(hc.hps)-1]
	return res
}

func (hc *HttpConnectorImpl) Recycle(hp *SocketProcessor) error {
	hc.mutex.Lock()
	defer func() {
		hc.mutex.Unlock()
	}()
	hc.hps = append(hc.hps, hp)
	return nil
}

func (hc *HttpConnectorImpl) Initialize() {
	hc.ctx = context.WithValue(context.Background(), logger.TraceID, uuid.New())
	hc.hps = make([]*SocketProcessor, 0, hc.maxProcessors)
	hc.minProcessors = hc.maxProcessors / 4
	for i := int32(0); i < hc.minProcessors; i++ {
		hp := NewSocketProcessor(hc, hc)
		hp.Start()
		hc.hps = append(hc.hps, hp)
	}
	hc.curProcessors = hc.minProcessors

	//定期删除无效的sessions
	go func() {
		defer func() {
			if err := recover(); err != nil {
				logger.LogWarn(hc.ctx, fmt.Sprintf("HttpConnectorImpl#Initialize recover get err:%v", err))
			}
		}()
		for true {
			hc.sk.DeprecateInvalidSessions()
			time.Sleep(time.Minute)
		}
	}()
}

func (hc *HttpConnectorImpl) GetProtocol() string {
	return hc.protocol
}

func (hc *HttpConnectorImpl) SetProtocol(protocol string) {
	panic("implement me")
}

func (hc *HttpConnectorImpl) GetSessionKeeper() *SessionKeeper {
	return hc.sk
}
