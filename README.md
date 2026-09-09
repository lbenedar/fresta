# Fresta - Foundry REST API (WORK IN PROGRESS)

## 📄 **About**

Fresta (Foundry REST API) - service written in Go, that provides access to FoundryVTT functions via a REST API and enables bot integration using gRPC.
Data is accessed directly through WebSocket requests to FoundryVTT and stored in a database, ensuring availability when switching between FoundryVTT operating modes.

---

## 🚀 **Quick Start**

Get the entire stack up and running with the single command:

```bash
docker compose up
```

The application will be accessible at **`http://localhost:9090`** by default

---

## 🛠️ **Make commands** 

```bash
make api/run                                    # Runs Fresta main API
make discord/run                                # Runs Fresta Discord bot
make docker/start                               # Runs API from docker compose
make db/migration/new name=NEW_MIGRATION_NAME   # Creates new migration file with NEW_MIGRATION_NAME
make db/migration/up                            # Apply all migrations
make proto/generate name=PROTO_FILE_NAME        # Generate go files from proto file with PROTO_FILE_NAME
make env/create                                 # Creates all .env files from .env.example
make help                                       # Return help from main API
```

---

## 📋 **Prerequisites**
Ensure you have the following installed before starting:
* **Docker Desktop** (v4.0.0+) or Docker Engine (v24.0.0+) with the **Docker compose V2** plugin
* **Go** (v1.25+) *-optional, only required for building and running binaries natively outside of Docker*
* **Make** (v4.0.0+)

---

## 🌲 **Project Structure**

```text
├── cmd/                                    
│   ├── api/                        # Main API entry point
│   │   ├── api.env.example         # Example env file for API
│   └── discord_bot/                # Discord Bot entry point
│       ├── discord.env.example     # Example env file for API
├── db/
│   └── migrations/                 # SQL migration files
├── docker-compose.yml              # Docker Compose orchestration file
├── Dockerfile                      # Multi-stage optimized Dockerfile
├── go.mod                          # Golang configuration file of module
├── internal/                       # Private application code
│   ├── foundry/                    # Code for FoundryVTT integration
│   ├── grpc/
│   │   └── bot/                    # gRPC server for bots
│   └── models/
│       ├── db/                     # Models for DB insertion of Foundry VTT data
│       └── json/                   # Models for data received from Foundry VTT 
├── Makefile                        # Project Makefile
├── proto/
│   ├── bots.proto                  # Proto file for bots gRPC
│   └── health.proto                # Proto file for healthcheck gRPC
└── README.md                       # Project Documentation
```

---

## ⚙️ **Environment Configuration**

The application uses environment variables for configuration. Copy the example file to create your local configurations:

```bash
cp -n cmd/api/api.env.example cmd/api/api.env;\
cp -n cmd/discord_bot/discord.env.example cmd/discord_bot/discord.env;\
cp -n .env.example .env
```

### Essential `api_docker.env` Settings

```api_docker.env
API_PORT=9090
GRPC_PORT=40404
```

### Essential `api.env` Settings

```api.env
GO_ENV=development # development|staging|production

SERVICE_MODE=api # api

FOUNDRY_PASS=admin_pass
FOUNDRY_HOST=localhost:9090
FOUNDRY_WORLDS=world1 # format: "world", "world1,world2,world3", "world1,world2,world3"
WORLD_USER=user1 # format: "user", "user1,user2,user3", "userForAll"
WORLD_PASS=pass1 # format: "pass", "pass1,pass2,pass3", "passForAll"

DB_DSN="file:db/default.db?cache=shared"
DB_DSN_MIGRATE=sqlite3://db/default.db?query

LOG_LEVEL=info # debug|info|warn|error
NO_FOUNDRY=true
```

### Essential `discord.env` Settings

```discord.env
API_HOST=localhost:40404

DISCORD_API_KEY="T1gQxQr72TU4NDGYrVOj.ajDk2MAwgo9TozOOKLIoNDTNU93aT-snZ4hY2SOgTMQx.k0NTf_bE"
```

---

## ⚖️ **MIT License**

Copyright (c) 2026 lbenedar

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
