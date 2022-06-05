package web

import (
	"L0/cmd/logger"
	"L0/pkg/cache"
	"L0/pkg/repository"
	"context"
	"github.com/jackc/pgconn"
	"gopkg.in/yaml.v2"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Application struct {
	Server       Server               `yaml:"server"`
	Logger       *logger.Logger       `yaml:"-"`
	DB           *repository.Postgres `yaml:"sql"`
	ConnectionDB *pgconn.PgConn       `yaml:"-"`
	Cache        cache.Cache          `yaml:"-"`
}

func NewApplication(configPath string) (*Application, error) {

	app := &Application{
		Logger: logger.NewLogger(),
		DB:     repository.NewPostgres(),
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

func (a Application) Run() {
	runChan := make(chan os.Signal, 1)

	ctx, cancel := context.WithTimeout(context.Background(), a.Server.Timeout.Server)
	defer cancel()

	var err error
	a.ConnectionDB, err = a.DB.GetConnection()
	if err != nil {
		a.Logger.ErrorLog.Printf("Error open connection DB: %+v", err)
		return
	}
	defer func(ConnectionDB *pgconn.PgConn, ctx context.Context) {
		err := ConnectionDB.Close(ctx)
		if err != nil {
			a.Logger.ErrorLog.Println("Error close DB connection:%+v", err)
		}
	}(a.ConnectionDB, context.Background())

	server := &http.Server{
		Addr:         a.Server.Host + ":" + a.Server.Port,
		Handler:      a.NewRouter(),
		ReadTimeout:  a.Server.Timeout.Read * time.Second,
		WriteTimeout: a.Server.Timeout.Write * time.Second,
		IdleTimeout:  a.Server.Timeout.Idle * time.Second,
	}

	signal.Notify(runChan, os.Interrupt, syscall.SIGTSTP)

	a.Logger.InfoLog.Printf("Server is starting on : %s\n", server.Addr)

	go func() {
		if err := server.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
			} else {
				a.Logger.ErrorLog.Printf("Server failed start: %v\n", err)
			}
		}
	}()

	interrupt := <-runChan

	a.Logger.InfoLog.Printf("Server is shutting down to %+v\n", interrupt)
	if err := server.Shutdown(ctx); err != nil {
		a.Logger.InfoLog.Printf("Server was unable to gracefully shutdown due to err: %+v", err)
		return
	}
}
