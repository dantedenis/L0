package web

import (
	"L0/cmd/logger"
	"L0/pkg/cache"
	"L0/pkg/repository"
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/nats-io/stan.go"
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
	ConnectionDB *pgx.Conn            `yaml:"-"`
	Cache        *cache.Cache         `yaml:"-"`
}

func NewApplication(configPath string) (*Application, error) {
	app := &Application{
		Logger: logger.NewLogger(),
		DB:     repository.NewPostgres(),
		Cache:  cache.NewCache(),
	}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			app.Logger.ErrorLog.Println("Error close config file:%+v", err)
			return
		}
	}(file)

	if err = yaml.NewDecoder(file).Decode(&app); err != nil {
		return nil, err
	}
	return app, nil
}

func (a Application) Run() {
	var err error
	runChan := make(chan os.Signal, 1)
	ctx, cancel := context.WithTimeout(context.Background(), a.Server.Timeout.Server)
	defer cancel()

	//Start connection DB
	a.ConnectionDB, err = a.DB.GetConnection()
	if err != nil {
		a.Logger.ErrorLog.Printf("Error open connection DB: %+v", err)
		return
	}
	defer func(ConnectionDB *pgx.Conn, ctx context.Context) {
		err := ConnectionDB.Close(ctx)
		if err != nil {
			a.Logger.ErrorLog.Println("Error close DB connection:%+v", err)
		}
	}(a.ConnectionDB, context.Background())

	//Restore cache in the BaseData
	err = a.Cache.RestoreCache(a.ConnectionDB, "select * from test_table;")
	if err != nil {
		a.Logger.ErrorLog.Println("Error restore cache")
		return
	}

	// Start connection NATS
	sConn, err := stan.Connect("test-cluster", "client-wb")
	if err != nil {
		a.Logger.ErrorLog.Println("Error open Nats-streaming connection:%+v", err)
		return
	}
	defer func(sConn stan.Conn) {
		err := sConn.Close()
		if err != nil {
			a.Logger.ErrorLog.Println("Error close Nats-streaming:%+v", err)
			return
		}
	}(sConn)

	//Subscribe on the channel NATS
	sub, err := sConn.Subscribe("test", func(m *stan.Msg) {
		a.PutAll(m)
	}, stan.StartAtTimeDelta(time.Minute))
	if err != nil {
		a.Logger.ErrorLog.Println("Error make subscribe on Nats-channel :%+v", err)
		return
	}
	defer func(sub stan.Subscription) {
		err := sub.Unsubscribe()
		if err != nil {
			a.Logger.ErrorLog.Println("Error unsubscribe Nats channel:%+v", err)
			return
		}
	}(sub)

	server := &http.Server{
		Addr:         a.Server.Host + ":" + a.Server.Port,
		Handler:      a.NewRouter(),
		ReadTimeout:  a.Server.Timeout.Read * time.Second,
		WriteTimeout: a.Server.Timeout.Write * time.Second,
		IdleTimeout:  a.Server.Timeout.Idle * time.Second,
	}

	signal.Notify(runChan, os.Interrupt, syscall.SIGTSTP)

	a.Logger.InfoLog.Printf("Server is starting on : %s\n", server.Addr)

	//Routine for listening
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
