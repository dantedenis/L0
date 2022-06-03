package web

import (
	"fmt"
	"net/http"
	"time"
)

type Server struct {
	Host    string `yaml:"host"`
	Port    string `yaml:"port"`
	Timeout struct {
		Server time.Duration `yaml:"server"`
		Write  time.Duration `yaml:"write"`
		Read   time.Duration `yaml:"read"`
		Idle   time.Duration `yaml:"idle"`
	} `yaml:"timeout"`
}

func NewRouter() *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("/welcome", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, "Hello, you have request: %s\n", request.URL.Path)
	})
	return router
}

func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != '/' {
		app.NotFound(w)
		return
	}
	
	files := []string {
		"ui/html/home.page.tmpl"
		"ui/html/base.layout.tmpl"
	}
	
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}
	
	err = ts.Execute(w, nil)
	if err != nil {
		app.serverError(w, err)
	}
}
