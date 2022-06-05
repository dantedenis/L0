package web

import (
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
	router.HandleFunc("/1", a.getID)
	return router
}

//Path for routers
func (a *Application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		a.notFound(w)
		return
	}

	files := []string{
		"ui/html/home.page.tmpl",
		"ui/html/base.layout.tmpl",
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

func (a *Application) getID(w http.ResponseWriter, r *http.Request) {
	err := a.GetExecID(w, "select * from orders")
	if err != nil {
		a.Logger.ErrorLog.Println(err)
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
