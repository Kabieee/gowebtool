package main

import (
	"gowebtool/router"
	"gowebtool/task"
)

func main() {
	task.Start(new(task.EmailTask))
	router.InitRouter()
	router.Engine.Run(":8001")
}
