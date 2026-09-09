package core

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lbenedar/fresta/internal/foundry"
	"github.com/lbenedar/fresta/internal/foundry/requests"
	"github.com/lbenedar/fresta/internal/foundry/transport"
	"github.com/lbenedar/fresta/internal/foundry/types"
	"github.com/lbenedar/fresta/proto/grpc_health"
)

type service_mode string

type db_config struct {
	dsn          string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

type config struct {
	Mode      service_mode
	NoFoundry bool
	Port      string
	GrpcPort  string
	Env       string
	db        db_config
}

type Application struct {
	Cfg        config
	GrpcStatus grpc_health.HealthCheckResponse_ServingStatus
	FoundryApp *foundry.FoundryApi
	Slogger    *slog.Logger
	Logger     *log.Logger
	Version    string
	BuildTime  string
}

type envConfig struct {
	host      string
	pass      string
	worlds    string
	worldUser string
	worldPass string
	mode      string
	logLevel  string
	noFoundry bool
	port      string
	grpcPort  string
	env       string

	db db_config
}

func CreateApplication(version string, buildTime string) *Application {
	return &Application{FoundryApp: &foundry.FoundryApi{}, Version: version, BuildTime: buildTime, GrpcStatus: grpc_health.HealthCheckResponse_NOT_SERVING}
}

func (app *Application) OpenDB(cfg config) (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite3", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(cfg.db.maxIdleConns)
	db.SetMaxOpenConns(cfg.db.maxOpenConns)

	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (app *Application) ParseEnv() *envConfig {
	cfg := &envConfig{}

	cfg.host = os.Getenv("FOUNDRY_HOST")
	if cfg.host == "" {
		cfg.host = "127.0.0.1"
	}

	cfg.pass = os.Getenv("FOUNDRY_PASS")
	cfg.worlds = os.Getenv("FOUNDRY_WORLDS")
	cfg.worldUser = os.Getenv("WORLD_USER")
	cfg.worldPass = os.Getenv("WORLD_PASS")

	cfg.mode = os.Getenv("SERVICE_MODE")
	if cfg.mode == "" {
		cfg.mode = "api"
	}

	cfg.port = os.Getenv("API_PORT")
	if cfg.port == "" {
		cfg.port = ":9090"
	} else {
		cfg.port = fmt.Sprintf(":%s", cfg.port)
	}

	cfg.grpcPort = os.Getenv("GRPC_PORT")
	if cfg.grpcPort == "" {
		cfg.grpcPort = ":40404"
	} else {
		cfg.grpcPort = fmt.Sprintf(":%s", cfg.grpcPort)
	}

	cfg.logLevel = os.Getenv("LOG_LEVEL")
	if cfg.logLevel == "" {
		cfg.logLevel = "info"
	}

	noFoundryStr := os.Getenv("NO_FOUNDRY")
	noFoundry, err := strconv.ParseBool(noFoundryStr)
	if err == nil {
		cfg.noFoundry = noFoundry
	} else {
		cfg.noFoundry = true
	}

	cfg.env = os.Getenv("GO_ENV")
	if cfg.env == "" {
		cfg.env = "development"
	}

	cfg.db.dsn = os.Getenv("DB_DSN")
	if cfg.db.dsn == "" {
		cfg.db.dsn = "file:db/default.db?cache=shared"
	}

	maxOpenConnsStr := os.Getenv("DB_MAX_OPEN_CONNS")
	maxOpenConns, err := strconv.Atoi(maxOpenConnsStr)
	if err == nil {
		cfg.db.maxOpenConns = maxOpenConns
	} else {
		cfg.db.maxOpenConns = 25
	}

	maxOpenIdleStr := os.Getenv("DB_MAX_IDLE_CONNS")
	maxIdleConns, err := strconv.Atoi(maxOpenIdleStr)
	if err == nil {
		cfg.db.maxIdleConns = maxIdleConns
	} else {
		cfg.db.maxIdleConns = 25
	}

	cfg.db.maxIdleTime = os.Getenv("DB_MAX_IDLE_TIME")
	if cfg.db.maxIdleTime == "" {
		cfg.db.maxIdleTime = "15m"
	}

	return cfg
}

func (app *Application) ParseFlags() *transport.FoundryTransportData {
	envCfg := app.ParseEnv()

	flag.StringVar(&app.Cfg.Port, "port", envCfg.port, "REST API service port")
	flag.StringVar(&app.Cfg.GrpcPort, "grpc-port", envCfg.grpcPort, "GRPC service port")

	flag.StringVar(&app.Cfg.Env, "env", envCfg.env, "Enivornment (development|staging|production)")

	flag.StringVar(&app.Cfg.db.dsn, "db-dsn", envCfg.db.dsn, "PostgreSQL DSN")
	flag.IntVar(&app.Cfg.db.maxOpenConns, "db-max-open-conns", envCfg.db.maxOpenConns, "PostgreSQL max open connections")
	flag.IntVar(&app.Cfg.db.maxIdleConns, "db-max-idle-conns", envCfg.db.maxIdleConns, "PostgreSQL max idle connections")
	flag.StringVar(&app.Cfg.db.maxIdleTime, "db-max-idle-time", envCfg.db.maxIdleTime, "PostgreSQL max connection idle time")

	var mode string
	flag.StringVar(&mode, "service_mode", envCfg.mode, "Type of service mode(api|discord|tg)")
	flag.BoolVar(&app.Cfg.NoFoundry, "no-foundry", envCfg.noFoundry, "Run server without connecting to foundry")

	var foundryTransportData transport.FoundryTransportData
	foundryTransportData.HttpConfig = &requests.FoundryHttpRequest{SessionID: nil}
	flag.StringVar(&foundryTransportData.HttpConfig.Host, "foundry_host", envCfg.host, "Address to connect to Foundry")
	flag.StringVar(&foundryTransportData.HttpConfig.Password, "foundry_pass", envCfg.pass, "Password to connect to Foundry")

	var logLevelStr string
	flag.StringVar(&logLevelStr, "log_level", envCfg.logLevel, "Password to connect to Foundry(debug|info|warn|error)")

	var worlds string
	flag.StringVar(&worlds, "foundry-worlds", envCfg.worlds, "Worlds that initialize in db on startup (format: \"test\", \"test1,test2,test3\", \"test1,test2,test3\")")
	var users string
	flag.StringVar(&users, "world-user", envCfg.worldUser, "Usernames to authenticate the world (format: \"test\", \"test1,test2,test3\", \"testForAll\")")
	var passwords string
	flag.StringVar(&passwords, "world-pass", envCfg.worldPass, "Passwords to authenticate the world (format: \"test\", \"test1,test2,test3\", \"testForAll\")")
	flag.Parse()

	app.Cfg.Mode = service_mode(mode)

	logLevel, ok := StringToLogLevel[logLevelStr]
	if !ok {
		logLevel = slog.LevelInfo
	}

	app.Slogger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
	if !ok {
		app.Slogger.Warn("Wrong log level. Started with log_level=info", "received", logLevelStr)
	}
	foundryTransportData.Logger = app.Slogger

	worldsData, err := types.CreateWorldDataSlice(worlds, users, passwords)
	if err != nil {
		app.Slogger.Warn("Got err on parsing authentication world data", "err", err)
	}
	foundryTransportData.Worlds = worldsData
	return &foundryTransportData
}
