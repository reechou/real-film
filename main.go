package main

import (
	"github.com/reechou/real-film/config"
	"github.com/reechou/real-film/controller"
)

func main() {
	controller.NewLogic(config.NewConfig()).Run()
}
