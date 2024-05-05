package test

import (
	"fmt"

	"github.com/alikWu/go-tomcat/servlet"
)

type TestFilter struct {
}

func (tf *TestFilter) GetFilterName() string {
	return "TestFilter"
}

func (tf *TestFilter) DoFilter(request servlet.ServletRequest, response servlet.ServletResponse, chain servlet.FilterChain) error {
	fmt.Println("TestFilter, The first Filter!!")
	return chain.DoFilter(request, response)
}

func (tf *TestFilter) GetFilterMatch() []*servlet.FilterMatch {
	return []*servlet.FilterMatch{servlet.NewFilterMatchServlet("HelloServlet")}
}
