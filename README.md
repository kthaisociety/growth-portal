# growth-portal

Dashboards for the KTHAIS growth team: who is in the Luma "Community" tier, broken down by their
registration answers (graduation year, areas of interest, …). A proof of concept: until identity-service
is live, a dev user stands in for sign-in ([PCR 0006](docs/pcr/0006-simulated-sign-in-until-identity-service.md)).

One repo, one image: a Go server serves the JSON API under `/api` and the frontend (Vite + React, embedded
in the binary) on every other path. Member data comes from [luma-service](../luma-service).

## Running locally

Needs Go 1.26, Node 24 and Docker, plus luma-service running locally.

```sh
# 1. luma-service, in ../luma-service (see its README): listens on 127.0.0.1:8001
# 2. this app's database and config
docker compose up -d                      # Postgres on 127.0.0.1:5435
cp .env.example .env                      # LUMA_SERVICE_SECRET = luma-service's
# 3. the backend
go run ./cmd/server                       # 127.0.0.1:8000
# 4. the frontend, with hot reload
cd frontend && npm install && npm run dev      # http://127.0.0.1:5173, /api proxied to 8000
```

Or build the frontend once (`cd frontend && npm run build`) and open http://127.0.0.1:8000: the Go server
serves it, as in production.

## Tests

```sh
make test             # both halves; run `npm install` in frontend/ first
make test-backend     # go vet, go test -race ./...; database tests need `docker compose up -d`
make test-frontend    # lint, Vitest, build, in frontend/
```

## Database

The schema is SQL files in `internal/database/migrations/`, applied by the server at startup
([PCR 0007](docs/pcr/0007-sql-migrations.md)).

## API

| Route | Returns |
|---|---|
| `GET /api/health` | `{"service":"growth-portal","status":"healthy"}`, or 503 if the database is unreachable |
| `GET /api/me` | the signed-in user: `{"sub":"…","email":"…","name":"…"}`; 401 if no one is |
| `GET /api/community/demographics` | `{"member_count":N,"questions":[{"label":"…","answered":n,"values":[{"value":"2027","count":12}]}]}` |

## Configuration

All from env, read by the Go server ([PCR 0001](docs/pcr/0001-config-from-env.md)); see `.env.example`.

| Variable | Secret | |
|---|---|---|
| `DATABASE_URL` | yes | Postgres connection string |
| `LUMA_SERVICE_URL` | no | luma-service's base URL |
| `LUMA_SERVICE_SECRET` | yes | luma-service's shared secret |
| `DEV_USER_EMAIL` | no | signs every request in as this user until identity-service is live; never in production |
| `HOST`, `PORT` | no | default all interfaces, `8000` |

## Delivery

Like every KTHAIS app ([infrastructure/docs/app-delivery.md](https://github.com/kthaisociety/infrastructure/blob/main/docs/app-delivery.md)):
PRs need a Conventional Commit title; every merge to `main` builds `ghcr.io/kthaisociety/growth-portal:sha-<7>`
and deploys it to staging; merging the release PR releases `X.Y.Z` for production. Frontend and backend
share one version.

## Docs

- [`CONTEXT.md`](CONTEXT.md): vocabulary
- [`docs/adr/`](docs/adr/): decisions about one part of the code
- [`docs/pcr/`](docs/pcr/): conventions for the whole project
