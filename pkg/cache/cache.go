package cache

import (
	"L0/pkg/model"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	"reflect"
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
		err = c.Set(data.ID, data.Data)
		if err != nil {
			fmt.Println(data.ID, "error type model")
		}
	}
	return nil
}

// Set - GET methods (Asynchronous with mutex)
func (c *Cache) Set(key int, value model.Model) error {
	var cmp model.Model
	if reflect.DeepEqual(value, cmp) {
		return errors.New("error type model")
	}
	c.Lock()
	defer c.Unlock()

	c.data[key] = value
	return nil
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
