# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Ride-sharing backend (Uber-style) built with Go microservices, Docker, Kubernetes, and a Next.js frontend. Module name: `ride-sharing`, Go 1.23.0.

## Commands

### Proto code generation
```bash
make generate-proto
# Runs: protoc --proto_path=proto --go_out=. --go-grpc_out=. proto/*.proto
# Requires: protoc, protoc-gen-go, protoc-gen-go-grpc
```

### Build Go services (local, for Tilt dev loop)
```bash
# API Gateway
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/api-gateway ./services/api-gateway

# Trip Service
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/trip-service ./services/trip-service/cmd/main.go
```

### Run locally with Tilt (primary dev workflow)
```bash
tilt up
# Requires: Docker, kubectl, a K8s cluster (e.g. minikube)
# Port forwards: API Gateway → 8081, Web → 3000
```

### Frontend (web/)
```bash
cd web
npm install
npm run dev       # Dev server on port 3000
npm run build     # Production build
```

### Tests
No tests exist yet. Standard Go testing would be:
```bash
go test ./...
```

## Architecture

### Services

**API Gateway** (`services/api-gateway/`) — HTTP/WebSocket BFF on port 8081. Uses plain `net/http` (no framework). Routes: `POST /trip/preview`, `/ws/drivers`, `/ws/riders`. Translates HTTP/WS to gRPC calls to trip-service.

**Trip Service** (`services/trip-service/`) — gRPC server on port 9093. Follows Clean Architecture:
- `internal/domain/` — interfaces and models
- `internal/service/` — business logic
- `internal/infrastructure/grpc/` — gRPC handler
- `internal/infrastructure/repository/` — in-memory repository (currently)
- `pkg/types/` — public types with `ToProto()` conversion

Entry point: `services/trip-service/cmd/main.go`

### Communication Flow
```
Frontend (Next.js) → HTTP/WS → API Gateway → gRPC → Trip Service → HTTP → OSRM (routing API)
```

Planned: RabbitMQ for async inter-service events (contracts defined in `shared/contracts/`, not yet implemented).

### Proto definitions
- Source: `proto/trip.proto`
- Generated Go code: `shared/proto/trip/`
- Always regenerate after editing `.proto` files: `make generate-proto`

### Shared packages (`shared/`)
- `contracts/` — AMQP, HTTP, WebSocket message contracts (event/command definitions for planned services)
- `types/` — shared domain types (Coordinate, Route, Geometry)
- `env/` — environment variable helpers
- `retry/` — exponential backoff utility
- `proto/trip/` — generated protobuf/gRPC Go code

### Infrastructure
- `infra/development/` — dev Dockerfiles (copy pre-built binaries) + K8s manifests
- `infra/production/` — multi-stage Dockerfiles (build in container) + K8s manifests (GCP Artifact Registry)
- `Tiltfile` — orchestrates local dev: compiles Go, builds Docker, deploys to K8s, live-reload via `restart_process` extension

### Scaffolding new services
```bash
go run tools/create_service.go <service-name>
```
Generates a service directory following the Clean Architecture pattern.

## Environment Variables

| Variable | Default | Used by |
|---|---|---|
| `HTTP_ADDR` | `:8081` | API Gateway listen address |
| `TRIP_SERVICE_URL` | `trip-service:9093` | API Gateway → Trip Service gRPC |
| `ENVIRONMENT` | `development` | All services |
| `NEXT_PUBLIC_API_URL` | — | Frontend API endpoint |
| `NEXT_PUBLIC_WEBSOCKET_URL` | — | Frontend WebSocket endpoint |
