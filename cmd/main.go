package main

import (
	initModule "github.com/sheginabo/go-quick-grpc/init"
)

func main() {
	initProcess := initModule.NewMainInitProcess("./")
	initProcess.Run()
}
