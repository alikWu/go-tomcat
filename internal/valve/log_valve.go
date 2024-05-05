package valve

import (
	"fmt"

	"github.com/alikWu/go-tomcat/internal"
	"github.com/alikWu/go-tomcat/internal/logger"
	"github.com/bytedance/sonic"
)

type LogValve struct {
	ValveBase
}

func NewLogValve() *LogValve {
	return &LogValve{}
}

//这种职责链模式，方便每个节点（Valve）自行决定下一个节点业务逻辑与当前节点业务逻辑的先后顺序，
//本质上就是执行本业务逻辑是在Servlet执行前还是后，这样就不用设计执行前和执行后两条职责链了。
func (lv *LogValve) Invoke(request internal.HttpServletRequest, response internal.HttpServletResponse, ctx internal.ValveContext) error {
	//先调用context中的invokeNext，实现职责链调用
	//Pass this request on to the next valve in our pipeline
	if err := ctx.InvokeNext(request, response); err != nil {
		return err
	}

	//本valve本身的业务逻辑
	reqBytes, err := sonic.Marshal(request)
	if err != nil {
		logger.LogError(ctx.GetContext(), "LogValve err:", err)
		return err
	}
	resBytes, err := sonic.Marshal(response)
	if err != nil {
		logger.LogError(ctx.GetContext(), "LogValve err:", err)
		return err
	}
	logger.LogInfof(ctx.GetContext(), "LogValve request=%s, response=%s", string(reqBytes), string(resBytes))
	fmt.Println(fmt.Sprintf("LogValve request=%s, response=%s", string(reqBytes), string(resBytes)))
	return nil
}
