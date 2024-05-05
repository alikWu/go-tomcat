package core

import (
	"github.com/pkg/errors"
	"go-tomcat/internal"
	"go-tomcat/internal/logger"
	"go-tomcat/internal/valve"
)

type StandardContextValve struct {
	valve.ValveBase
}

func NewStandardContextValve() *StandardContextValve {
	return &StandardContextValve{}
}

func (scv *StandardContextValve) GetInfo() string {
	return "StandardContextValve 0.1"
}

func (scv *StandardContextValve) Invoke(request internal.HttpServletRequest, response internal.HttpServletResponse, ctx internal.ValveContext) error {
	logger.LogInfo(ctx.GetContext(), "StandardContextValve invoke")
	sc, ok := scv.ValveBase.GetContainer().(internal.Context)
	if !ok {
		err := errors.New("StandardContextValve convert container type fail")
		logger.LogErrorf(ctx.GetContext(), "StandardContextValve#Invoke err", err)
		return err
	}

	url := request.GetPathInfo()
	wrapper := sc.GetWrapper(url)
	if wrapper == nil {
		logger.LogWarnf(ctx.GetContext(), "StandardContextValve invoke can't find the url=%s", url)
		response.SetStatus(int32(internal.SC_NOT_FOUND))
		return nil
	}
	return wrapper.Invoke(ctx.GetContext(), request, response)
}
