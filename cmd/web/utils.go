package web

import (
	"context"
	"fmt"
	"net/http"
)

func (a *Application) GetExecID(w http.ResponseWriter, request string) error {

	response := a.ConnectionDB.ExecParams(context.Background(), request, nil, nil, nil, nil)

	fmt.Fprintf(w, "Connection OK\n")
	fmt.Fprintf(w, "%s\n", response)

	for response.NextRow() {
		fmt.Println(string(response.Values()[1]))
	}

	return nil
}
