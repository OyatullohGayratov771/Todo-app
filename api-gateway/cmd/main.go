package main

import (
	"fmt"
	"api-gateway/config"
)

func main() {
	config.LoadConfig()
	fmt.Print(config.AppConfig.Http.Host,config.AppConfig.Http.Port)
}
