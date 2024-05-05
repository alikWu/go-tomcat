package core

import (
	"context"

	"github.com/alikWu/go-tomcat/servlet"
)

type servletContext struct {
	ctx context.Context
}

func NewServletContext(ctx context.Context) servlet.ServletContext {
	return &servletContext{
		ctx: ctx,
	}
}

func (sc *servletContext) GetContext() context.Context {
	return sc.ctx
}
