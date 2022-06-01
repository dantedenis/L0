package main

import (
	"L0/server"
	"fmt"
)

func main() {
	config, err := server.NewConfig("config.yml")
	if err != nil {
		fmt.Println(err)
	}
	config.Run()
}
