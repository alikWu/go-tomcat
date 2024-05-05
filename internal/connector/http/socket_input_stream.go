package http

import (
	"context"
	"io"
	"sync"

	"github.com/pkg/errors"
	"github.com/alikWu/go-tomcat/internal/logger"
)

const (
	lcOffset = byte('a') - byte('A')
	cr       = byte('\r')
	lf       = byte('\n')
	sp       = byte(' ')
	ht       = byte('\t')
	colon    = byte(':')
)

type SocketInputStream struct {
	mutex sync.Mutex
	buf   []byte
	//buf中的数据数量
	count int64
	//当前读到的位置
	pos         int64
	inputStream io.Reader
	ctx         context.Context
}

func NewSocketInputStream(ctx context.Context, inputStream io.Reader, bufferSize int64) *SocketInputStream {
	return &SocketInputStream{
		buf:         make([]byte, bufferSize),
		inputStream: inputStream,
		ctx:         ctx,
	}
}

//从输入流中解析出request line
func (sis *SocketInputStream) ReadRequestLine(requestLine *HttpRequestLine) error {
	chr, err := sis.readNextByte()
	if err != nil {
		return err
	}
	//跳过空行
	for chr == cr || chr == lf {
		if chr, err = sis.readNextByte(); err != nil {
			return err
		}
	}

	sis.pos--
	readCount := 0
	for true {
		if chr, err = sis.readNextByte(); err != nil {
			logger.LogError(sis.ctx, "SocketInputStream#ReadRequestLine readNextByte err", err)
			return errors.Wrap(err, "SocketInputStream#ReadRequestLine read method err")
		}

		if chr == sp {
			break
		}
		requestLine.Method[readCount] = chr
		readCount++
	}
	requestLine.MethodEnd = int64(readCount)

	readCount = 0
	for true {
		if chr, err = sis.readNextByte(); err != nil {
			logger.LogError(sis.ctx, "SocketInputStream#ReadRequestLine readNextByte err", err)
			return errors.Wrap(err, "SocketInputStream#ReadRequestLine read uri err")
		}

		if chr == sp {
			break
		}
		requestLine.Uri[readCount] = chr
		readCount++
	}
	requestLine.UriEnd = int64(readCount)

	readCount = 0
	for true {
		if chr, err = sis.readNextByte(); err != nil {
			logger.LogError(sis.ctx, "SocketInputStream#ReadRequestLine readNextByte err", err)
			return errors.Wrap(err, "SocketInputStream#ReadRequestLine read protocol err")
		}
		if chr == cr {
			continue
		}
		if chr == lf {
			break
		}

		requestLine.Protocol[readCount] = chr
		readCount++
	}
	requestLine.ProtocolEnd = int64(readCount)
	return nil
}

//从输入流中解析出request header
func (sis *SocketInputStream) ReadHeader(header *HttpHeader) error {
	b, err := sis.readNextByte()
	if err != nil {
		return err
	}
	if b == cr {
		_, err = sis.readNextByte()
		if err != nil {
			return err
		}
		header.NameEnd = int64(0)
		header.ValueEnd = int64(0)
		return nil
	} else if b == lf {
		header.NameEnd = int64(0)
		header.ValueEnd = int64(0)
		return nil
	} else {
		sis.pos--
	}

	//read header name
	readCount := 0
	for true {
		if b, err = sis.readNextByte(); err != nil {
			logger.LogError(sis.ctx, "SocketInputStream#ReadHeader readNextByte err", err)
			return errors.Wrap(err, "SocketInputStream#ReadHeader read head name err")
		}

		if b == colon {
			break
		}
		header.Name[readCount] = toLower(b)
		readCount++
	}
	header.NameEnd = int64(readCount)

	//read header value（a value may cross multiple rows）
	readCount = 0
	for true {
		//remove leading space
		for true {
			if b, err = sis.readNextByte(); err != nil {
				logger.LogError(sis.ctx, "SocketInputStream#ReadHeader readNextByte err", err)
				return errors.Wrap(err, "SocketInputStream#ReadHeader remove leading space err")
			}
			if b == sp || b == ht {
			} else {
				sis.pos--
				break
			}
		}

		//read a line
		for true {
			if b, err = sis.readNextByte(); err != nil {
				logger.LogError(sis.ctx, "SocketInputStream#ReadHeader readNextByte err", err)
				return errors.Wrap(err, "SocketInputStream#ReadHeader read head value err")
			}
			if b == cr {
			} else if b == lf {
				break
			} else {
				header.Value[readCount] = b
				readCount++
			}
		}

		//对于多行value，换行后需要缩进
		if b, err = sis.readNextByte(); err != nil {
			logger.LogError(sis.ctx, "SocketInputStream#ReadHeader readNextByte err", err)
			return errors.Wrap(err, "SocketInputStream#ReadHeader read head value err")
		}
		if b != sp && b != ht {
			sis.pos--
			break
		}
		header.Value[readCount] = ' '
		readCount++
	}
	header.ValueEnd = int64(readCount)
	return nil
}

func (sis *SocketInputStream) Read(p []byte) (n int, err error) {
	sis.mutex.Lock()
	defer func() {
		sis.mutex.Unlock()
	}()

	if sis.pos < sis.count {
		n = copy(sis.buf[sis.pos:sis.count], p)
		sis.pos += int64(n)
		return
	}
	if err = sis.fill(); err != nil {
		return 0, err
	}
	n = copy(sis.buf[sis.pos:sis.count], p)
	sis.pos += int64(n)
	return n, nil
}

func (sis *SocketInputStream) readNextByte() (byte, error) {
	if sis.pos >= sis.count {
		if err := sis.fill(); err != nil {
			return 0, err
		}
		if sis.pos >= sis.count {
			return byte(0), nil
		}
	}
	b := sis.buf[sis.pos]
	sis.pos++
	return b, nil
}

func (sis *SocketInputStream) fill() error {
	sis.pos = 0
	sis.count = 0
	n, err := sis.inputStream.Read(sis.buf)
	if err != nil {
		return err
	}
	if n > 0 {
		sis.count = int64(n)
	}
	logger.LogInfof(sis.ctx, "SocketInputStream#fill get request=%s", string(sis.buf[:sis.count]))
	return nil
}

func toLower(b byte) byte {
	if b >= 'A' && b <= 'Z' {
		return b + lcOffset
	}
	return b
}
