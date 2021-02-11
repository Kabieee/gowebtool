package router

import (
	"gowebtool/controller"

	"github.com/gin-gonic/gin"
)

var Engine *gin.Engine

func InitRouter() {
	e := gin.Default()
	site := new(controller.SiteController)
	e.Any("/", site.Index)
	e.Any("/index", site.Index)
	gitGroup := e.Group("/git")
	{
		gitGroup.Use(CheckGitHubToken())
		git := new(controller.GitController)
		gitGroup.POST("/self", git.Self)
	}
	Engine = e
}
