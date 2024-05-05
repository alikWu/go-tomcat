package core

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/alikWu/go-tomcat/internal"
	"github.com/alikWu/go-tomcat/internal/logger"
	"github.com/alikWu/go-tomcat/internal/util"
	"github.com/alikWu/go-tomcat/internal/valve"
	"github.com/alikWu/go-tomcat/servlet"
)

type StandardWrapperValve struct {
	valve.ValveBase
}

func NewStandardWrapperValve() *StandardWrapperValve {
	return &StandardWrapperValve{}
}

func (swv *StandardWrapperValve) Invoke(request internal.HttpServletRequest, response internal.HttpServletResponse, ctx internal.ValveContext) error {
	logger.LogInfo(ctx.GetContext(), "StandardWrapperValve Invoke")

	instance := swv.GetContainer().(internal.Wrapper).GetServlet()
	if instance == nil {
		err := errors.New("StandardWrapperValve can't find servlet")
		logger.LogError(ctx.GetContext(), "StandardWrapperValve Invoke", err)
		return err
	}

	filterChain := swv.createFilterChain(request, instance)
	if err := filterChain.DoFilter(request, response); err != nil {
		logger.LogError(ctx.GetContext(), "StandardWrapperValve DoFilter", err)
		response.SetStatus(int32(internal.SC_INTERNAL_SERVER_ERROR))
	}
	filterChain.Release()
	return nil
}

func (swv *StandardWrapperValve) createFilterChain(request internal.HttpServletRequest, servlet servlet.Servlet) *ApplicationFilterChain {
	filterChain := NewApplicationFilterChain()
	filterChain.SetServlet(servlet)
	wrapper := swv.GetContainer().(*standardWrapper)
	parent := wrapper.GetParent()
	context := parent.(*StandardContext)

	filterMaps := context.GetFilterMaps()
	if len(filterMaps) == 0 {
		return filterChain
	}

	//遍历filter Map，找到匹配URL模式的filter，加入到filterChain中
	requestPath := request.GetPathInfo()
	for _, filterMap := range filterMaps {
		if !util.MatchUrl(filterMap.GetURLPattern(), requestPath) {
			continue
		}
		filter := context.GetFilter(filterMap.GetFilterName())
		if filter == nil {
			fmt.Println("filter is nil")
			continue
		}

		filterChain.AddFilter(filter)
	}

	//遍历filter Map，找到匹配servlet的filter，加入到filterChain中
	servletName := wrapper.GetName()
	for _, filterMap := range filterMaps {
		if filterMap.GetServletName() != servletName {
			continue
		}
		filter := context.GetFilter(filterMap.GetFilterName())
		if filter == nil {
			fmt.Println("filter is nil")
			continue
		}

		filterChain.AddFilter(filter)
	}
	return filterChain
}
