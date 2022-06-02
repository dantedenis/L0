package web

import (
	"context"
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Config struct {
	Server struct {
		Host    string `yaml:"host"`
		Port    string `yaml:"port"`
		Timeout struct {
			Server time.Duration `yaml:"server"`
			Write  time.Duration `yaml:"write"`
			Read   time.Duration `yaml:"read"`
			Idle   time.Duration `yaml:"idle"`
		} `yaml:"timeout"`
	} `yaml:"server"`
}

func NewConfig(configPath string) (*Config, error) {
	config := &Config{}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	d := yaml.NewDecoder(file)
	if err = d.Decode(&config); err != nil {
		return nil, err
	}
	return config, nil
}

func NewRouter() *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("/welcome", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, "Hello, you have request: %s\n", request.URL.Path)
	})
	return router
}

func (config Config) Run() {
	runChan := make(chan os.Signal, 1)

	ctx, cancel := context.WithTimeout(context.Background(), config.Server.Timeout.Server)
	defer cancel()

	server := &http.Server{
		Addr:         config.Server.Host + ":" + config.Server.Port,
		Handler:      NewRouter(),
		ReadTimeout:  config.Server.Timeout.Read * time.Second,
		WriteTimeout: config.Server.Timeout.Write * time.Second,
		IdleTimeout:  config.Server.Timeout.Idle * time.Second,
	}

	signal.Notify(runChan, os.Interrupt, syscall.SIGTSTP)

	log.Printf("Server is starting on : %s\n", server.Addr)

	go func() {
		if err := server.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {

			} else {
				log.Fatalf("Server failed start: %v\n", err)
			}
		}
	}()

	interrupt := <-runChan

	log.Printf("Server is shutting down to %+v\n", interrupt)
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server was unable to gracefully shutdown due to err: %+v", err)
	}
}
