package core

import (
	"context"
	"sync"

	"github.com/alikWu/go-tomcat/internal"
	"github.com/alikWu/go-tomcat/internal/logger"
)

const (
	ADD_CHILD_EVENT string = "addChild"
)

type ContainerBase struct {
	name   string
	parent internal.Container

	mutex    sync.Mutex
	children map[string]internal.Container

	pipeline internal.Pipeline
}

func (c *ContainerBase) SetPipeline(pipeline internal.Pipeline) {
	c.pipeline = pipeline
}

func (c *ContainerBase) GetPipeline() internal.Pipeline {
	return c.pipeline
}

func (c *ContainerBase) GetName() string {
	return c.name
}

func (c *ContainerBase) SetName(name string) {
	c.name = name
}

func (c *ContainerBase) GetParent() internal.Container {
	return c.parent
}

func (c *ContainerBase) SetParent(parent internal.Container) {
	c.parent = parent
}

func (c *ContainerBase) AddChild(child internal.Container) {
	c.mutex.Lock()
	defer func() {
		c.mutex.Unlock()
	}()
	if c.children == nil {
		c.children = make(map[string]internal.Container)
	}
	c.children[child.GetName()] = child
}

func (c *ContainerBase) FindChild(name string) internal.Container {
	c.mutex.Lock()
	defer func() {
		c.mutex.Unlock()
	}()
	if c.children == nil {
		c.children = make(map[string]internal.Container)
	}
	return c.children[name]
}

func (c *ContainerBase) Invoke(ctx context.Context, request internal.HttpServletRequest, response internal.HttpServletResponse) error {
	logger.LogInfo(ctx, "ContainerBase invoke")
	return c.pipeline.Invoke(ctx, request, response)
}

func (c *ContainerBase) GetBasic() internal.Valve {
	return c.pipeline.GetBasic()
}

func (c *ContainerBase) SetBasic(valve internal.Valve) {
	c.pipeline.SetBasic(valve)
}

func (c *ContainerBase) AddValve(valve internal.Valve) {
	c.pipeline.AddValve(valve)
}

func (c *ContainerBase) GetValves() []internal.Valve {
	return c.pipeline.GetValves()
}

func (c *ContainerBase) RemoveValve(valve internal.Valve) {
	c.pipeline.RemoveValve(valve)
}
