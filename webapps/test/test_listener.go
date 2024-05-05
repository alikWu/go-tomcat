package test

import (
	"github.com/alikWu/go-tomcat/internal"
)

type TestListener struct {
}

func (tl *TestListener) ContainerEvent(event *internal.ContainerEvent) {
	//bytes, err := sonic.Marshal(event)
	//if err != nil {
	//	fmt.Println("TestListener err", err)
	//}
	//fmt.Println("TestListener ", string(bytes))
}
