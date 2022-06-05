package cache

import (
	"L0/pkg/model"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
)

type Cache struct {
	data map[string]model.Model
}

func (c *Cache) FillCache(conn *pgx.Conn, request string) error {
	response, err := conn.Exec(context.Background(), request)
	if err != nil {
		return err
	}
	for line := range response {
		fmt.Println(line)
	}
	return nil
}
