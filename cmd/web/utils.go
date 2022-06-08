package web

import (
	"L0/pkg/model"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/nats-io/stan.go"
)

func (a *Application) PutAll(m *stan.Msg) {
	var data model.Model
	tempUUID := uuid.New()
	if err := json.Unmarshal(m.Data, &data); err != nil {
		a.Logger.ErrorLog.Println("Error decode Json-msg:%+v", err)
		return
	}
	err := a.Cache.Set(tempUUID.String(), data)
	if err != nil {
		a.Logger.ErrorLog.Println("Error set to cache: %s", err)
		return
	}
	_, err = a.ConnectionDB.Exec(context.Background(), "insert into test_table (id, order_data) values ($1, $2);", tempUUID, m.Data)
	if err != nil {
		a.Logger.ErrorLog.Println("Error insert to BD: %s", err)
		return
	}
	fmt.Println("Insert to DB and Cache - ok")
}
