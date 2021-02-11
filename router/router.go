package router

import (
	"gowebtool/controller"

	"github.com/gin-gonic/gin"
)

var (
	Engine *gin.Engine
	base   = new(controller.BaseController)
)

func InitRouter() {
	e := gin.Default()
	e.Any("/", base.Index)
	e.Any("/index", base.Index)
	e.GET("/user", base.User)
	e.POST("/user", base.User)
	gitGroup := e.Group("/git")
	{
		gitGroup.Use(CheckGitHubToken())
		git := new(controller.GitController)
		gitGroup.POST("/self", git.Self)
	}

	toolGroup := e.Group("/tool")
	{
		tool := new(controller.ToolController)
		toolGroup.GET("/SendEmail", tool.SendEmail)
		toolGroup.POST("/SendEmail", tool.SendEmail)
	}
	Engine = e
}
