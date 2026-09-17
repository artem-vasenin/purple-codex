# Backend

Go backend scaffold.

## Run

```sh
docker compose -f ../docker-compose.yml up -d postgres
go run ./cmd/server
```

Set `JWT_SECRET` and, when needed, `DATABASE_URL`; optional settings are `APP_ADDR`, `ACCESS_TTL`, `REFRESH_TTL`, and `COOKIE_SECURE`. Apply `migrations/*.sql` with `golang-migrate` before starting the server. Auth endpoints are `POST /auth/register`, `/auth/login`, `/auth/refresh`, and `/auth/logout`; refresh tokens are stored in an HttpOnly cookie.

From the repository root, run both applications together with `./scripts/dev.sh`.
The script requires Docker Compose, starts the Postgres container automatically,
and waits until it accepts connections. Database data is stored in the named
Docker volume `purple-codex_postgres_data` and remains after the script exits.
Stop the database manually with `docker compose down` when it is no longer needed.
