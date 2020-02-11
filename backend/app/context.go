package app

import (
	"github.com/ajorgensen/go-react-webapp-bootstrap/backend/db"
	"github.com/sirupsen/logrus"
)

type Context struct {
	Logger        logrus.FieldLogger
	RemoteAddress string
	Database      *db.Database
}

func (ctx *Context) WithLogger(logger logrus.FieldLogger) *Context {
	ret := *ctx
	ret.Logger = logger
	return &ret
}

func (ctx *Context) WithRemoteAddress(remoteAddress string) *Context {
	ret := *ctx
	ret.RemoteAddress = remoteAddress
	return &ret
}
