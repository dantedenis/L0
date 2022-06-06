package cache

import (
	"L0/pkg/model"
	"context"
	"github.com/jackc/pgx/v4"
	"sync"
)

type Cache struct {
	sync.RWMutex
	data map[int]model.Model
}

func NewCache() *Cache {
	return &Cache{
		data: make(map[int]model.Model),
	}
}

func (c *Cache) RestoreCache(conn *pgx.Conn, request string) error {
	response, err := conn.Query(context.Background(), request)
	if err != nil {
		return err
	}
	for response.Next() {
		data := struct {
			Data model.Model `db:"order_data"`
			ID   int         `db:"id"`
		}{}
		err := response.Scan(&data.ID, &data.Data)
		if err != nil {
			return err
		}
		c.Set(data.ID, data.Data)
	}
	return nil
}

// Set - GET methods (Asynchronous with mutex)
func (c *Cache) Set(key int, value model.Model) {
	c.Lock()
	defer c.Unlock()

	c.data[key] = value
}

func (c *Cache) Get(key int) (model.Model, bool) {
	c.RLock()

	defer c.RUnlock()

	item, found := c.data[key]
	if !found {
		return model.Model{}, found
	}

	return item, true
}

// end methods
