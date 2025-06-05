# VENDIQ2

VENDIQ2 is a set of small services for managing prices on Amazon. The repository contains a Go backend, a job worker and a React front‑end.

## Directory overview

- **api** – simple REST API exposing pricing data from MySQL.
- **pricer** – main service implementing price adjustment logic, RabbitMQ worker and HTTP endpoints.
- **webui** – Vite + React front‑end to view and edit pricing and queue jobs.
- **compose.yaml** – Docker Compose configuration for local development (frontend, pricer and RabbitMQ).
- **Makefile** – helper targets to build and run the Docker environment.

## Getting started

Build the containers:

```bash
make build
```

Start the stack:

```bash
make up
```

The front‑end will be available at `http://localhost:5173`.

To stop and remove containers:

```bash
make down
```

Database credentials and Amazon API tokens are loaded from Docker secrets defined in `compose.yaml`.
