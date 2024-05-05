package core

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/alikWu/go-tomcat/internal"
	"github.com/alikWu/go-tomcat/internal/logger"
	"github.com/alikWu/go-tomcat/internal/util"
	"github.com/alikWu/go-tomcat/internal/valve"
	"github.com/alikWu/go-tomcat/servlet"
	"github.com/alikWu/go-tomcat/webapps/test"
)

//现在使用了容器技术之后，Engine 和 Host 的概念已经名存实亡类。
//实际上，当我们部署的时候，一个 Tomcat 一般就只用一个 Engine 和一个 Host，如果需要多个，就用多个容器
//因此，不打算实现 Engine 和 Host，只用Contexts和Wrapper两层Container结构.一个 Context 组件可以处理一个特定 Web 应用的所有请求
type StandardContext struct {
	hc internal.Connector
	ContainerBase

	//key=servletName
	servletMaps []*ServletMap

	//key=filterName
	filters    map[string]servlet.Filter
	filterMaps []*FilterMap

	listenerDefs    []*ContainerListenerDef
	listeners       []internal.ContainerListener
	listenerFactory internal.ListenerFactory
}

func NewStandardContext() *StandardContext {
	return &StandardContext{
		filters: make(map[string]servlet.Filter),
	}
}

func (sc *StandardContext) SetFactory(factory *standardFactory) {
	sc.listenerFactory = factory
}

func (sc *StandardContext) Start() {
	sc.SetName("StandardContext")

	pipeline := NewStandardPipeline(sc)
	pipeline.SetBasic(NewStandardContextValve())
	pipeline.AddValve(valve.NewLogValve())
	sc.SetPipeline(pipeline)

	//url匹配优先级： /* < /path1/* < /path1/path2
	sort.Slice(sc.servletMaps, func(i, j int) bool {
		if sc.servletMaps[i].GetURLPattern() == util.MatchAll {
			return false
		}
		if sc.servletMaps[i].GetURLPattern()[:2] == util.MatchAll {
			return sc.servletMaps[j].GetURLPattern() == util.MatchAll
		}
		if strings.Index(sc.servletMaps[j].GetURLPattern(), util.MatchAll) >= 0 {
			return true
		}
		return false
	})

	cld := &ContainerListenerDef{}
	cld.SetListenerName("TestListener")
	sc.AddContainerListener(&test.TestListener{})
	sc.AddListenerDef(cld)
	sc.ListenerStart()

	fmt.Println("Container start!!!!!")
	sc.fireContainerEvent("Container start", sc)
}

func (sc *StandardContext) RegisterServlet(servlet servlet.Servlet) {
	for _, urlPattern := range servlet.GetMatchedUrlPattern() {
		sc.servletMaps = append(sc.servletMaps, NewServletMap(servlet.GetServletName(), urlPattern))
	}
	servletWrapper := NewServletWrapper(servlet, sc, getWrapperName(servlet.GetServletName()))
	sc.AddChild(servletWrapper)
}

func (sc *StandardContext) AddContainerListener(listener internal.ContainerListener) {
	sc.listeners = append(sc.listeners, listener)
}

func (sc *StandardContext) RemoveContainerListener(listener internal.ContainerListener) {

}

func (sc *StandardContext) fireContainerEvent(kind string, data interface{}) {
	event := internal.NewContainerEvent(sc, kind, data)
	for _, listener := range sc.listeners {
		listener.ContainerEvent(event)
	}
}

func (sc *StandardContext) AddListenerDef(listenerDef *ContainerListenerDef) {
	sc.listenerDefs = append(sc.listenerDefs, listenerDef)
}

func (sc *StandardContext) ListenerStart() bool {
	for _, def := range sc.listenerDefs {
		listener := sc.listenerFactory.CreateContainerListener(def.GetListenerName())
		sc.AddContainerListener(listener)
	}
	return true
}

func (sc *StandardContext) SetConnector(connector internal.Connector) {
	sc.hc = connector
}

func (sc *StandardContext) GetConnector() internal.Connector {
	return sc.hc
}

func (sc *StandardContext) GetServletContext() servlet.ServletContext {
	//TODO implement me
	panic("implement me")
}

func (sc *StandardContext) GetSessionTimeout() int64 {
	//TODO implement me
	panic("implement me")
}

func (sc *StandardContext) SetSessionTimeout(timeout int64) {
	//TODO implement me
	panic("implement me")
}

func (sc *StandardContext) Reload() {
	//TODO implement me
	panic("implement me")
}

func (sc *StandardContext) Invoke(ctx context.Context, request internal.HttpServletRequest, response internal.HttpServletResponse) error {
	logger.LogInfo(ctx, "StandardContext invoke")
	return sc.ContainerBase.Invoke(ctx, request, response)
}

func (sc *StandardContext) GetWrapper(url string) internal.Container {
	var targetServletName string
	for _, servletMap := range sc.servletMaps {
		if util.MatchUrl(servletMap.GetURLPattern(), url) {
			targetServletName = servletMap.GetServletName()
			break
		}
	}
	return sc.FindChild(getWrapperName(targetServletName))
}

func (sc *StandardContext) GetFilterMap(filterName string) *FilterMap {
	for _, filterMap := range sc.filterMaps {
		if filterMap.GetFilterName() == filterName {
			return filterMap
		}
	}
	return nil
}

func (sc *StandardContext) GetFilterMaps() []*FilterMap {
	return sc.filterMaps
}

func (sc *StandardContext) RegisterFilter(filter servlet.Filter) {
	sc.filters[filter.GetFilterName()] = filter
	for _, filterMatch := range filter.GetFilterMatch() {
		filterName := filter.GetFilterName()
		servletName := filterMatch.GetServletName()
		urlPattern := filterMatch.GetUrlPattern()

		if len(servletName) == 0 && len(urlPattern) == 0 {
			logger.LogWarn(sc.GetServletContext().GetContext(), "StandardContext add filterMap without servletName and urlPattern")
			continue
		}
		if len(servletName) > 0 && len(urlPattern) > 0 {
			logger.LogWarn(sc.GetServletContext().GetContext(), "StandardContext add filterMap with servletName and urlPattern")
			continue
		}

		sc.filterMaps = append(sc.filterMaps, &FilterMap{
			filterName:  filterName,
			servletName: servletName,
			urlPattern:  urlPattern,
		})
	}
}

func (sc *StandardContext) GetFilter(filterName string) servlet.Filter {
	return sc.filters[filterName]
}

func getWrapperName(servletName string) string {
	return fmt.Sprintf("%sWrapper", servletName)
}
