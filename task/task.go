package task

import (
	"fmt"

	"github.com/fatih/color"
)

type Task interface {
	Run()
}

//var sy sync.WaitGroup

func Start(task ...Task) {
	for _, v := range task {
		go v.Run()
	}
	fmt.Println(color.BlueString("%s", "Start Task"))
	//sy.Wait()
}
