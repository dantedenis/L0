package main

import (
	"L0/cmd/web"
	"fmt"
)

func main() {
	config, err := web.NewConfig("config.yml")
	if err != nil {
		fmt.Println(err)
	}
	config.Run()
}
