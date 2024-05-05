package webapps

import (
	"fmt"

	"go-tomcat/internal/logger"
	"go-tomcat/servlet"
)

type HelloServlet struct {
	count int
}

func (h *HelloServlet) GetServletName() string {
	return "HelloServlet"
}

func (h *HelloServlet) Service(request servlet.ServletRequest, response servlet.ServletResponse) error {
	doc1 := "<!DOCTYPE html> \n" + "<html>\n" + "<head><meta charset=\"utf-8\"><title>Test</title></head>\n" +
		"<body bgcolor=\"#f0f0f0\">\n" + "<h1 align=\"center\">"
	_, err := response.GetWriter().Write([]byte(doc1))
	if err != nil {
		logger.LogError(request.GetServletContext().GetContext(), "HelloServlet#Service writeBack err:", err)
	}
	doc2 := "Hello servlet!!" + "</h1>\n"
	_, err = response.GetWriter().Write([]byte(doc2))
	if err != nil {
		logger.LogError(request.GetServletContext().GetContext(), "HelloServlet#Service writeBack err:", err)
	}
	fmt.Println("att1", request.GetAttribute("att1"))
	request.SetAttribute("att1", h.count)
	h.count++
	return err
}

func (h *HelloServlet) GetMatchedUrlPattern() []string {
	return []string{"/servlet/*"}
}

func (h *HelloServlet) GetServletInfo() string {
	return ""
}
func (h *HelloServlet) Destroy() {
}
