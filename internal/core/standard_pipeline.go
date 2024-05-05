package core

import (
	"context"
	"fmt"
	"sync"

	"github.com/pkg/errors"
	"github.com/alikWu/go-tomcat/internal"
	"github.com/alikWu/go-tomcat/internal/logger"
)

type StandardPipeline struct {
	basic     internal.Valve
	container internal.Container

	mutex  sync.Mutex
	valves []internal.Valve
}

func NewStandardPipeline(container internal.Container) *StandardPipeline {
	return &StandardPipeline{
		container: container,
	}
}

func (s *StandardPipeline) GetBasic() internal.Valve {
	return s.basic
}

func (s *StandardPipeline) SetBasic(valve internal.Valve) {
	oldBasic := s.basic
	if oldBasic == valve {
		return
	}

	if valve == nil {
		return
	}
	valve.SetContainer(s.container)
	s.basic = valve
}

func (s *StandardPipeline) AddValve(valve internal.Valve) {
	s.mutex.Lock()
	defer func() {
		s.mutex.Unlock()
	}()
	valve.SetContainer(s.container)
	s.valves = append(s.valves, valve)
}

func (s *StandardPipeline) GetValves() []internal.Valve {
	if s.basic == nil {
		return s.valves
	}

	return append(s.valves, s.basic)
}

func (s *StandardPipeline) Invoke(ctx context.Context, request internal.HttpServletRequest, response internal.HttpServletResponse) error {
	logger.LogInfo(ctx, "StandardPipeline Invoke")
	//转而调用valveContext中的invokeNext，发起职责链调用
	return NewStandardPipelineValveContext(ctx, s).InvokeNext(request, response)
}

func (s *StandardPipeline) RemoveValve(valve internal.Valve) {
	s.mutex.Lock()
	defer func() {
		s.mutex.Unlock()
	}()
	j := 0

	for ; j < len(s.valves); j++ {
		if s.valves[j] == valve {
			break
		}
	}
	if j >= len(s.valves) {
		return
	}

	valve.SetContainer(nil)
	s.valves = append(s.valves[:j], s.valves[j+1:]...)
}

type standardPipelineValveContext struct {
	stage int
	sp    *StandardPipeline
	ctx   context.Context
}

func NewStandardPipelineValveContext(ctx context.Context, sp *StandardPipeline) *standardPipelineValveContext {
	return &standardPipelineValveContext{
		sp:  sp,
		ctx: ctx,
	}
}

func (spvc *standardPipelineValveContext) GetInfo() string {
	return "standardPipelineValveContext 0.1"
}

func (spvc *standardPipelineValveContext) InvokeNext(request internal.HttpServletRequest, response internal.HttpServletResponse) error {
	logger.LogInfo(spvc.ctx, "standardPipelineValveContext InvokeNext")
	curStage := spvc.stage
	spvc.stage++

	var err error
	if curStage < len(spvc.sp.valves) {
		err = spvc.sp.valves[curStage].Invoke(request, response, spvc)
	} else if curStage == len(spvc.sp.valves) {
		err = spvc.sp.basic.Invoke(request, response, spvc)
	} else {
		err = errors.New(fmt.Sprintf("abnormal stage, curStage=%d, len(Stage)=%d", curStage, len(spvc.sp.valves)))
	}
	return err
}

func (spvc *standardPipelineValveContext) GetContext() context.Context {
	return spvc.ctx
}
