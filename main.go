package main

import (
	"gowebtool/common"
	"gowebtool/router"
	"gowebtool/task"
)

func main() {
	common.InitTranslator()
	task.Start(new(task.EmailTask))
	router.InitRouter()
	router.Engine.Run(":8001")
}
