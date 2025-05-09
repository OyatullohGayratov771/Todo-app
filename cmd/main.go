package main

import (
	"fmt"
	"todo-app/config"
)

func main() {
	config.LoadConfig()
	fmt.Print(config.AppConfig.Http.Host,config.AppConfig.Http.Port)
}
