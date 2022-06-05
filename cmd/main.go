package main

import (
	"L0/cmd/web"
	"fmt"
)

func main() {
	app, err := web.NewApplication("config.yml")
	if err != nil {
		fmt.Println(err)
		return
	}
	app.Run()
}
