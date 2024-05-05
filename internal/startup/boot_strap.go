package startup

import (
	"fmt"

	"github.com/alikWu/go-tomcat/internal/connector/http"
	"github.com/alikWu/go-tomcat/internal/core"
	"github.com/alikWu/go-tomcat/servlet"
)

type BootStrap struct {
	port      string
	connector *http.HttpConnectorImpl
	context   *core.StandardContext
}

func NewBootStrap(port string) *BootStrap {
	hc := http.NewHttpConnector(port)
	sc := core.NewStandardContext()
	hc.SetContainer(sc)
	sc.SetConnector(hc)

	return &BootStrap{
		port:      port,
		connector: hc,
		context:   sc,
	}
}

func (bs *BootStrap) SetMaxConnections(maxConnections int32) {
	bs.connector.SetMaxConnections(maxConnections)
}

func (bs *BootStrap) SetConnectionTimeout(connectionTimeout int32) {
	bs.connector.SetConnectionTimeout(connectionTimeout)
}

func (bs *BootStrap) RegisterServlets(servlets []servlet.Servlet) {
	for _, s := range servlets {
		bs.context.RegisterServlet(s)
	}
}

func (bs *BootStrap) RegisterFilters(filters []servlet.Filter) {
	for _, filter := range filters {
		bs.context.RegisterFilter(filter)
	}
}

func (bs *BootStrap) Start() error {
	bs.context.SetFactory(core.NewStandardFactory())
	bs.context.Start()

	bs.connector.Initialize()
	err := bs.connector.ListenConnect()
	if err != nil {
		fmt.Println(err)
	}
	return err
}
