-include .env

ifeq (${MAKECMDGOALS},api/run)
	-include cmd/api/api.env
endif

ifeq (${MAKECMDGOALS},discord/run)
	-include cmd/discord_bot/discord.env
endif

export

api/run: 
	go run ./cmd/api -service_mode=${SERVICE_MODE} -foundry_pass=${FOUNDRY_PASS} -foundry_host=${FOUNDRY_HOST} -port=${API_PORT} -env=development -log_level=${LOG_LEVEL} -db-dsn=${DB_DSN} -foundry-worlds=${FOUNDRY_WORLDS} -world-user=${WORLD_USER} -world-pass=${WORLD_PASS} -no-foundry=${NO_FOUNDRY}

discord/run: 
	go run ./cmd/discord_bot

docker/up:
	VERSION=$(shell git describe --tags --always) docker compose -f docker-compose.yml up -d

docker/build:
	VERSION=$(shell git describe --tags --always) docker compose -f docker-compose.yml build


db/migration/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}

db/migrations/up:
	@echo 'Running up migrations...'
	migrate -path ./db/migrations -database ${DB_DSN_MIGRATE} up

proto/generate:
	@echo 'Creating generating proto files from ${name}...'
	protoc --proto_path=./proto --go_out=./proto --go-grpc_out=./proto ${name}

env/create:
	cp -n .env.example .env
	cp -n ./cmd/api/api.env.example ./cmd/api/api.env
	cp -n ./cmd/discord_bot/discord.env.example ./cmd/discord_bot/discord.env

help:
	go run ./cmd/api -help

.PHONY: api/run discord/run start db/migration/new db/migrations/up help