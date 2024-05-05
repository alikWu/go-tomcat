package core

import (
	"context"
	"sync"

	"go-tomcat/internal"
	"go-tomcat/internal/logger"
	"go-tomcat/internal/valve"
	"go-tomcat/servlet"
)

type standardWrapper struct {
	instance    servlet.Servlet
	servletName string

	ContainerBase

	once sync.Once
}

func NewServletWrapper(instance servlet.Servlet, parent *StandardContext, wrapperName string) internal.Wrapper {
	sw := &standardWrapper{
		instance:    instance,
		servletName: instance.GetServletName(),
	}
	sw.SetParent(parent)
	sw.SetName(wrapperName)
	return sw
}

func (sw *standardWrapper) GetName() string {
	return sw.ContainerBase.GetName()
}

func (sw *standardWrapper) Load() {
	pipeline := NewStandardPipeline(sw)
	basic := NewStandardWrapperValve()
	pipeline.SetBasic(basic)
	pipeline.AddValve(valve.NewLogValve())
	sw.SetPipeline(pipeline)
}

func (sw *standardWrapper) GetServlet() servlet.Servlet {
	return sw.instance
}

func (sw *standardWrapper) Invoke(ctx context.Context, request internal.HttpServletRequest, response internal.HttpServletResponse) error {
	sw.once.Do(func() {
		logger.LogInfo(ctx, "standardWrapper load")
		sw.Load()
	})
	logger.LogInfo(ctx, "standardWrapper invoke")
	return sw.ContainerBase.Invoke(ctx, request, response)
}
