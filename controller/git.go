package controller

import (
	"fmt"
	"gowebtool/common"

	"github.com/gin-gonic/gin"
)

type GitController struct {
	BaseController
}

func (g *GitController) Self(c *gin.Context) {
	g.ParseRequest(c)
	cmd := fmt.Sprintf("%s/git_self.sh", execPath)
	execCmd := common.ExecCmd(cmd)
	g.Success(c, &Success{Data: temp, Extend: execCmd})
}
