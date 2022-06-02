package web

import (
	"L0/cmd/logger"
	"L0/pkg/model"
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

type Application struct {
	server *Config        `yaml:"server"`
	logger *logger.Logger `yaml:"-"`
	db     *model.Model   `yaml:"-"`
}

func NewApplication(configPath string) (*Application, error) {
	app := &Application{
		logger: &logger.Logger{
			ErrorLog: log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
			InfoLog:  log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime),
		},
	}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	dec := yaml.NewDecoder(file)
	if err = dec.Decode(&app); err != nil {
		return nil, err
	}

	return app, nil
}
