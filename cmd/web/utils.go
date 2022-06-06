package web

import (
	"L0/pkg/model"
	"encoding/json"
	"fmt"
	"github.com/nats-io/stan.go"
)

func (a *Application) PutAll(m *stan.Msg) {
	var data model.Model

	if err := json.Unmarshal(m.Data, &data); err != nil {
		a.Logger.ErrorLog.Println("Error decode Json-msg:%+v", err)
		return
	}
	fmt.Sprintf("insert into public.orders (order_data) values ('$1'::jsonb);", m.Data)
	//a.ConnectionDB.ExecParams(context.Background(), "insert into public.orders (order_data) values ($1::jsonb);", nil, nil, nil, nil)
	//a.ConnectionDB.Exec(context.Background(), "insert into public.orders (order_data) values ($1::jsonb);")
	// TODO: add to BD return id!!
	//a.Cache.Set()
}
