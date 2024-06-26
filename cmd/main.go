package main

import (
	"fmt"

	GoTomcat "github.com/alikWu/go-tomcat"
	"github.com/alikWu/go-tomcat/servlet"
	"github.com/alikWu/go-tomcat/webapps"
	"github.com/alikWu/go-tomcat/webapps/test"
)

func main() {
	bs := GoTomcat.NewBootStrap(":8008")
	bs.SetMaxConnections(200)
	bs.RegisterServlets([]servlet.Servlet{&webapps.HelloServlet{}, &webapps.StaticResourceProcessor{}})
	bs.RegisterFilters([]servlet.Filter{&test.TestFilter{}})
	if err := bs.Start(); err != nil {
		fmt.Println("Start Fail ,err ", err)
	}
}
