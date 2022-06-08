package web

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"runtime/debug"
	"time"
)

type Server struct {
	Host    string `yaml:"host"`
	Port    string `yaml:"port"`
	Timeout struct {
		Server time.Duration `yaml:"Server"`
		Write  time.Duration `yaml:"write"`
		Read   time.Duration `yaml:"read"`
		Idle   time.Duration `yaml:"idle"`
	} `yaml:"timeout"`
}

func (a *Application) NewRouter() *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("/", a.home)
	router.HandleFunc("/show", a.showOrder)
	return router
}

//Path for routers
func (a *Application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		a.notFound(w)
		return
	}

	files := []string{
		"ui/html/index.html",
		"ui/html/header.html",
		"ui/html/closer.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		a.serverError(w, err)
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		a.serverError(w, err)
	}
}

func (a *Application) showOrder(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	data := struct {
		Count int
		UUIDs []string
		Model string
		Msg   string
	}{Count: a.Cache.Length(),
		UUIDs: a.Cache.GetAllUUID()}
	if id != "" {
		k, v := a.Cache.Get(id)
		if !v {
			data.Msg = "Not Found"
		} else {
			temp, err := json.MarshalIndent(k, "", "  ")
			if err != nil {
				fmt.Println("Error Marshalling")
				a.serverError(w, err)
			}
			data.Model = string(temp)
			data.Msg = id
		}
	}

	t := template.New("index")
	tmpl, err := t.ParseFiles("ui/html/index.html", "ui/html/header.html", "ui/html/closer.html")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		fmt.Println(err)
		return
	}
}

// Server Error
func (a *Application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	a.Logger.ErrorLog.Println(trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (a *Application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (a *Application) notFound(w http.ResponseWriter) {
	a.clientError(w, http.StatusNotFound)
}

//end server error
