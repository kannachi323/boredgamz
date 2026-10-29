# Deploying BoredGamz

One container serves everything: the Go binary hosts the built React app, the JSON API and the WebSockets. It needs PostgreSQL for accounts.

## Quick start (any Docker host)

```bash
cp .env.example .env        # fill in PUBLIC_ORIGIN, POSTGRES_PASSWORD, JWT_SECRET_KEY
docker compose up --build -d
curl -fsS http://localhost:3000/readyz
```

Put a TLS-terminating proxy (or your platform's edge) in front of port 3000.

## Environment variables

| Variable | Required | Default | Notes |
| --- | --- | --- | --- |
| `APP_ENV` | yes | `production` in the image | Enables Secure cookies, HSTS, JSON logs and strict startup checks. |
| `DATABASE_URL` | yes | – | `postgres://user:pass@host:5432/db?sslmode=require` for managed Postgres. Compose builds it from `POSTGRES_*`. |
| `JWT_SECRET_KEY` | yes | – | 32+ random bytes (`openssl rand -base64 48`). Rotating it signs everyone out. |
| `ALLOWED_ORIGINS` | yes | – | Exact browser origin(s), comma-separated, e.g. `https://boredgamz.com`. Compose sets it from `PUBLIC_ORIGIN`. Used for WebSocket Origin checks, CORS and CSRF. |
| `PORT` | no | `3000` | Listen port inside the container. |
| `TRUST_PROXY` | no | `false` | `true` only when a proxy you control overwrites `X-Forwarded-For`/`X-Real-IP`. Otherwise clients could spoof their IP and dodge rate limits. |
| `DB_MAX_CONNS` | no | `10` | Upper bound on the Postgres pool. Keep it below your plan's connection limit. |
| `LOG_LEVEL` | no | `info` | `debug` adds socket-close reasons. |
| `STATIC_DIR` | no | `/app/static` | Where the built SPA lives (set by the image). |

The server refuses to start in production without `DATABASE_URL`, `ALLOWED_ORIGINS` or a 32-byte `JWT_SECRET_KEY`.

## Database and migrations

Migrations are SQL files in `server/db/migrations/`, embedded in the binary and applied automatically at startup inside a transaction each, under a Postgres advisory lock (safe if two instances boot at once). Applied versions are recorded in `schema_migrations`.

- Add a migration: create `server/db/migrations/000N_description.sql`. Never edit one that has shipped.
- Existing databases created by the old `schema.sql` bootstrap upgrade cleanly (`0001` uses `IF NOT EXISTS`).
- The legacy `gomoku_games` and `connectfour_games` tables are no longer used. They are left in place; drop them by hand if you want.
- Back up the database. It holds only accounts (name, email, bcrypt hash), but losing it signs everyone out permanently.

## TLS, proxies and WebSockets

- Terminate TLS before the app. Session cookies are `Secure`, so they only work over HTTPS.
- The proxy must pass WebSocket upgrades (`Upgrade`/`Connection` headers) for `/api/v1/arena/*`, `/api/v1/poker` and `/api/v1/mahjong`, with an idle timeout of at least 75 seconds. The server pings every 54 seconds.
- `ALLOWED_ORIGINS` must match what the browser shows in the address bar, scheme included, with no trailing slash.
- Health checks: `GET /healthz` reports liveness. `GET /readyz` returns 503 while Postgres is unreachable, so use it for load balancer readiness.
- Shutdown: on SIGTERM the server stops accepting connections and drains HTTP for up to 10 seconds. Open game sockets close and clients reconnect to the new instance.

## Scaling limits

**Run exactly one instance.** Live game state (matchmaking queues, rooms, clocks, chat, spectators, reconnect grace) lives in process memory. Two replicas behind a load balancer would split players into rooms the other replica can't see. The in-memory rate limiters are per instance too.

One small instance handles a lot: the integration suite runs 24 concurrent matches over real sockets. Scale vertically first.

Consequences of the single-instance model:

- A deploy or crash ends in-progress games. Accounts survive because they're in Postgres.
- Rolling deploys can't hand matches over. Deploy when traffic is low, or accept the interruption.

Path to more than one instance, in order of effort:

1. **Sticky routing per room.** Give each room an owner instance and route `/api/v1/*?roomId=` and invite codes to it, using consistent hashing at the proxy or a small routing table. Matchmaking still needs one shared queue.
2. **Shared queue and registry in Redis.** Store waiting players, invite codes and the room→instance map in Redis. Forward moves from the instance holding the socket to the room's owner over Redis pub/sub or NATS. Move the rate limiters into Redis.
3. **Durable rooms.** Snapshot room state to Redis or Postgres after each move, so another instance can take over a room after a crash and games survive deploys.

## Operations checklist

- Logs are JSON on stdout: one `request` line per HTTP request, plus migration and startup lines.
- `GET /api/v1/status` shows queue, room and connection counts. It's public and contains no personal data.
- The container runs as a non-root user with a read-only root filesystem, `cap_drop: ALL` and `no-new-privileges` (see `compose.yaml`).
