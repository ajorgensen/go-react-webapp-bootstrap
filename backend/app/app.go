package app

import (
	"github.com/ajorgensen/go-react-webapp-bootstrap/backend/db"
	"github.com/sirupsen/logrus"
)

type App struct {
	Config   *Config
	Database *db.Database
}

func (a *App) NewContext() *Context {
	return &Context{
		Logger:   logrus.StandardLogger(),
		Database: a.Database,
	}
}

func New() (app *App, err error) {
	app = &App{}
	app.Config, err = InitConfig()
	if err != nil {
		return nil, err
	}

	dbConfig, err := db.InitConfig()
	if err != nil {
		return nil, err
	}

	app.Database, err = db.New(dbConfig)
	if err != nil {
		return nil, err
	}

	return app, err
}

func (a *App) Close() error {
	return a.Database.Close()
}
