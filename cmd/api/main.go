package main

import (
	"net/http/cookiejar"

	"github.com/lbenedar/fresta/cmd/api/core"
	restapi "github.com/lbenedar/fresta/cmd/api/rest_api"
	"github.com/lbenedar/fresta/internal/foundry/transport"
	grpcbot "github.com/lbenedar/fresta/internal/grpc/bot"
	"github.com/lbenedar/fresta/proto/grpc_health"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

const (
	API_MODE         = "api"
	DISCORD_BOT_MODE = "discord"
	TG_BOT_MODE      = "tg"
)

var Version = "0.1.0"
var BuildTime = ""

func grpcServe(app *core.Application) error {
	grpcServer, listener, err := grpcbot.CreateServer(app.Cfg.GrpcPort)
	if err != nil {
		return err
	}

	app.Slogger.Info("gRPC сервер запущен", "port", app.Cfg.GrpcPort)
	app.GrpcStatus = grpc_health.HealthCheckResponse_SERVING
	if err := grpcServer.Serve(listener); err != nil {
		app.Slogger.Error("Ошибка запуска сервера", "err", err)
		app.GrpcStatus = grpc_health.HealthCheckResponse_NOT_SERVING
		return err
	}
	app.GrpcStatus = grpc_health.HealthCheckResponse_NOT_SERVING

	return nil
}

func serve(app *core.Application) error {
	go grpcServe(app)

	switch app.Cfg.Mode {
	case API_MODE:
		return restapi.ListenAndServe((*restapi.AppServer)(app))
	default:
		app.Slogger.Error("Wrong Mode. Should be one of: 'api', 'discord', 'tg'", "mode", app.Cfg.Mode)
	}
	return nil
}

func main() {
	app := core.CreateApplication(Version, BuildTime)

	err := godotenv.Load()
	if err != nil {
		app.Slogger.Error("Error", "text", err.Error())
		return
	}

	foundryTransportData := app.ParseFlags()
	if foundryTransportData == nil {
		return
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		app.Slogger.Error("Error", "text", err.Error())
		return
	}
	foundryTransportData.HttpConfig.Jar = jar

	db, err := app.OpenDB(app.Cfg)
	if err != nil {
		app.Slogger.Error("Error when opening database connection", "err", err)
		return
	}
	defer db.Close()
	app.Slogger.Info("Database connection pool established")
	foundryTransportData.DbConn = db

	app.FoundryApp.SetTransport(transport.NewFoundryTransport(foundryTransportData))
	if !app.Cfg.NoFoundry {
		err = app.FoundryApp.PrepareDB()
		if err != nil {
			app.Slogger.Error("Error on preparing database", "err", err)
			return
		}

		go app.FoundryApp.StartListenFoundry()
	}

	err = serve(app)
	if err != nil {
		app.Slogger.Error("Error on serving application", "err", err)
	}
}
