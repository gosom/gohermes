package main

//go:generate sqlboiler psql

import (
	"os"

	"github.com/gosom/gohermes/pkg/container"
	"github.com/gosom/gohermes/pkg/server"

	"github.com/gosom/gohermes/examples/todo/routes"
	"github.com/gosom/gohermes/examples/todo/services"
)

func main() {
	di, err := container.NewDefault()
	if err != nil {
		panic(err)
	}

	if err := services.RegisterServices(di); err != nil {
		panic(err)
	}

	srv, err := server.New(di)
	if err != nil {
		panic(err)
	}

	if err := routes.SetUpRoutes(di, srv); err != nil {
		panic(err)
	}
	if err := srv.Run(); err != nil {
		di.Logger.Error().Msg(err.Error())
		os.Exit(2)
	} else {
		os.Exit(0)
	}
}
