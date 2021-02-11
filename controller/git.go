package controller

import (
	"fmt"
	"gowebtool/common"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GitController struct {
	SiteController
}

func (g *GitController) Self(c *gin.Context) {
	g.ParseRequest(c)
	cmd := fmt.Sprintf("git")
	execCmd := common.ExecCmd(cmd)
	c.JSON(http.StatusOK, gin.H{
		"Code": http.StatusOK,
		"Data": temp,
		"Exec": execCmd,
	})
}
