package main

import (
	"gowebtool/router"
)

func main() {
	router.InitRouter()
	router.Engine.Run(":8001")
}
