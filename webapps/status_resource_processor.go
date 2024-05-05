package webapps

import (
	"io/ioutil"
	"os"

	"go-tomcat/internal/logger"
	"go-tomcat/servlet"
)

type StaticResourceProcessor struct {
}

func (p *StaticResourceProcessor) GetServletName() string {
	return "StaticResourceProcessor"
}

func (p *StaticResourceProcessor) GetMatchedUrlPattern() []string {
	return []string{"/hello.txt"}
}
func (p *StaticResourceProcessor) Service(request servlet.ServletRequest, response servlet.ServletResponse) error {

	ctx := request.GetServletContext().GetContext()
	file, err := os.Open("/Users/alik.wu/Documents/repo/GoTomcat/webapps/webroot/hello.txt")
	if err != nil {
		logger.LogError(ctx, "Processor#Process os.Open err", err)
	}
	defer file.Close()

	res, err := ioutil.ReadAll(file)
	if err != nil {
		logger.LogError(ctx, "Processor#Process ioutil.ReadAll err:", err)
		return err
	}
	writer := response.GetWriter()
	if _, err = writer.Write(res); err != nil {
		logger.LogError(ctx, "Processor#Process response.Write err:", err)
		return err
	}

	return nil
}

func (p *StaticResourceProcessor) GetServletInfo() string {
	//TODO implement me
	panic("implement me")
}

func (p *StaticResourceProcessor) Destroy() {
	//TODO implement me
	panic("implement me")
}
