package router

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type AppRouter interface {
	RegisterEndpoint()
	Start(address string) error
}

type router struct {
	e *echo.Echo
}

func NewRouter(e *echo.Echo) router {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	return router{
		e: e,
	}
}

func (receiver router) RegisterEndpoint(httpMethod string, path string, handlerFunc echo.HandlerFunc) {
	receiver.e.Add(httpMethod, path, handlerFunc)
}

func (receiver router) Start(address string) error {
	return receiver.e.Start(address)
}
