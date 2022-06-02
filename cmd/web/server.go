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
