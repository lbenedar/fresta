# ===========================================
# BUILD STAGE
# ===========================================
FROM golang:1.26-alpine AS builder

RUN apk update && apk add --no-cache ca-certificates tzdata gcc musl-dev wget
ENV GRPC_HEALTH_PROBE_VERSION=v0.4.24
RUN wget -qO/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x /grpc_health_probe

WORKDIR /opt/app

COPY go.mod go.sum ./
RUN go mod download

ARG APP_VERSION
ARG TZ

COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \ 
    CGO_ENABLED=1 GOOS=linux go build -ldflags="-X 'main.Version=${APP_VERSION}' -X 'main.BuildTime=$(date +%Y-%m-%dT%H:%M:%SZ)' -w" -o /bin/fresta ./cmd/api
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=1 GOOS=linux go build -ldflags="-X 'main.Version=${APP_VERSION}' -X 'main.BuildTime=$(date +%Y-%m-%dT%H:%M:%SZ)' -w" -o /bin/discord-fresta ./cmd/discord_bot

# ===========================================
# RUN API STAGE
# ===========================================
FROM alpine:3.24 as fresta

RUN apk update && apk add --no-cache make

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

COPY --from=builder /bin/fresta /api/
COPY --from=builder /grpc_health_probe /bin/grpc_health_probe

WORKDIR /api

CMD ./fresta

# ===========================================
# RUN DISCORD BOT STAGE
# ===========================================
FROM alpine:3.24 as discord-fresta

RUN apk update && apk add --no-cache make

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

COPY --from=builder /bin/discord-fresta /bot/

WORKDIR /bot

RUN chmod +x discord-fresta

CMD ./discord-fresta