package main

import (
	"fmt"
	"html/template"
	"os"
)


type data struct {
	UUID []string
}

func main() {
/*
	file, err := os.Create("t.html")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
*/
	t := template.New("index")
	d := data{
		UUID: []string{"test1","test2","test3","test4"},
	}
	tmpl, err := t.ParseFiles("ui/html/index.html", "ui/html/header.html", "ui/html/close.html")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = tmpl.Execute(os.Stdout, d)
	if err != nil {
		fmt.Println(err)
		return
	}
}