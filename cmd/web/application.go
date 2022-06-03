package web

import (
	"L0/cmd/logger"
	"L0/pkg/model"
	"context"
	"gopkg.in/yaml.v2"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Application struct {
	server Server         `yaml:"server"`
	logger *logger.Logger `yaml:"-"`
	db     *model.Model   `yaml:"-"`
}

func NewApplication(configPath string) (*Application, error) {

	app := &Application{
		server: Server{},
		logger: &logger.Logger{
			ErrorLog: log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
			InfoLog:  log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime),
		},
		db: &model.Model{},
	}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	dec := yaml.NewDecoder(file)
	if err = dec.Decode(&app); err != nil {
		return nil, err
	}

	return app, nil
}

func (config Application) Run() {
	runChan := make(chan os.Signal, 1)

	ctx, cancel := context.WithTimeout(context.Background(), config.server.Timeout.Server)
	defer cancel()

	server := &http.Server{
		Addr:         config.server.Host + ":" + config.server.Port,
		Handler:      NewRouter(),
		ReadTimeout:  config.server.Timeout.Read * time.Second,
		WriteTimeout: config.server.Timeout.Write * time.Second,
		IdleTimeout:  config.server.Timeout.Idle * time.Second,
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

func (app *Application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Println(trace)
 
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}
