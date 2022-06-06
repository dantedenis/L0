package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
)

type Postgres struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"pass"`
	NameDB   string `yaml:"db_name"`
}

func NewPostgres() *Postgres {
	return &Postgres{}
}

/*
# Example DSN
user=jack password=secret host=pg.example.com port=5432 dbname=mydb sslmode=verify-ca

# Example URL
postgres://jack:secret@pg.example.com:5432/mydb?sslmode=verify-c
*/

func (p *Postgres) GetConnection() (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), fmt.Sprintf("postgres://%s:%s@%s:%s/%s", p.User, p.Password, p.Host, p.Port, p.NameDB))
	if err != nil {
		return nil, err
	}
	return conn, nil
}
