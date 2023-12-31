// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"minerva/internal/biz"
	"minerva/internal/conf"
	"minerva/internal/data"
	"minerva/internal/server"
	"minerva/internal/service"
)

import (
	_ "go.uber.org/automaxprocs"
)

// Injectors from wire.go:

// wireApp init kratos application.
func wireApp(confServer *conf.Server, confData *conf.Data, logger log.Logger) (*kratos.App, func(), error) {
	dataData, cleanup, err := data.NewData(confData, logger)
	if err != nil {
		return nil, nil, err
	}
	minervaRepo := data.NewMinervaRepo(dataData, logger)
	minervaUsecase := biz.NewMinervaUsecase(minervaRepo, logger)
	minervaService := service.NewMinervaService(minervaUsecase, minervaRepo)
	grpcServer := server.NewGRPCServer(confServer, minervaService, logger)
	httpServer := server.NewHTTPServer(confServer, minervaService, logger)
	app := newApp(logger, grpcServer, httpServer)
	return app, func() {
		cleanup()
	}, nil
}
