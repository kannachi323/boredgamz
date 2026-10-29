# BoredGamz

BoredGamz is an account-backed online game room for six familiar tabletop games. A React/TypeScript client provides six deliberately different table experiences; one Go service serves the production bundle, owns every live match over WebSockets, and persists player accounts in PostgreSQL.

Players can create an account, sign in, sign out, and resume a cookie-backed session after a reload. Their account identity and display name are used by the multiplayer server instead of trusting URL parameters. Every game also supports two human players, private codes, table chat, spectators, reconnect protection, and server-authoritative rules.

## Games

| Game | Online modes | Match options | Server authority |
| --- | --- | --- | --- |
| Tic-Tac-Toe | Classic 3×3, Four Square 4×4, Misère | Best of 1/3/5/7; untimed or 15/30/60-second moves | Placement, turns, series, misère result, clock |
| Connect Four | Classic 6×7, Quick 5×6, Connect Five 7×8 | Best of 1/3/5/7; untimed or 15/30/60-second moves | Gravity, full columns, connections, series, clock |
| Gomoku | Freestyle, Exact Five, Compact 13×13 | Best of 1/3/5; untimed or 30/60-second moves | Placement, five/overline rules, series, clock |
| Chess | Blitz 3+2, Rapid 10, Classical 15+10 | Public or private challenge; four board themes | Legal moves, clocks, SAN/FEN, checkmate, draws, resignations |
| Poker | Heads-up No-Limit Hold’em, Five-Card Draw | Public or private table; 1,000 free-play chips | CSPRNG deck, private cards, betting, pots, draws, showdown |
| Mahjong Solitaire | Shared Garden, Mirror Duel, Zen Partners | 144-tile timed co-op, 72-tile race, or untimed co-op | Free tiles, tile families, timers, hints, deadlock reshuffles |

Poker uses free-play chips only. There are no purchases, cash value, or gambling integrations.

## Repository layout

```text
app/
  src/components/games/        shared chat and table UI
  src/hooks/                   Arena, Poker, and Mahjong socket clients
  src/pages/Games/<Game>/      game-specific lobby, table, board, controls
  src/styles/                  modular game lobby/table/motion styles
server/
  core/arena/                  grid games and Chess hub/room/rules
  core/poker/                  cards, evaluator, table, room, hub, socket
  core/mahjong/                layouts, solver, room, hub, socket
  server/                      HTTP routing, security, static SPA, account auth
  tests/                       HTTP, WebSocket integration, and load tests
docs/PRODUCTION_REPORT.md      design and engineering decision report
```

The largest source file is under 1,000 lines. UI, protocol, rules, room lifecycle, and hub responsibilities are kept separate so game logic can be tested without a browser or network.

## Local development

Use Node.js 22+ and Go 1.26.5. Go 1.26.5 is the pinned release toolchain because it includes the current standard-library security fixes.

```bash
cd app
npm ci
VITE_WS_ROOT=ws://localhost:3000/api/v1/arena npm run dev
```

Start PostgreSQL for account signup/login:

```bash
cd server
docker compose up -d db
```

In another terminal:

```bash
cd server
ALLOWED_ORIGINS=http://localhost:5173 go run .
```

Open [http://localhost:5173](http://localhost:5173). `VITE_WS_ROOT` is an Arena root; the Poker and Mahjong clients derive their sibling WebSocket routes from it.

## Verification

```bash
cd app
npm run lint
npm test
npm run build
npm audit --audit-level=high

cd ../server
go test -race ./...
go vet ./...
```

The repository currently has 15 client tests and 70 Go tests. Coverage includes all game engines, matchmaking isolation, private codes, invalid/out-of-turn input, chat limits, clocks, reconnect generations, spectators, rematches, Gomoku replay history, Poker evaluation and chip conservation, Mahjong solvability and simultaneous claims, JWT separation, auth validation and rate limits, database fail-fast startup, HTTP security, SPA caching, and 24 concurrent real WebSocket matches.

The release vulnerability check is run with the same toolchain as the container:

```bash
docker run --rm -v "$PWD/server:/src:ro" -w /src golang:1.26.5-alpine \
  sh -c 'go run golang.org/x/vuln/cmd/govulncheck@latest ./...'
```

## Deploy with Docker

```bash
cp .env.example .env
# Set PUBLIC_ORIGIN to the exact browser-visible https origin.
docker compose up --build -d
curl http://localhost:3000/readyz
```

The image is multi-stage and contains only the compiled Go binary and hashed web bundle. Runtime characteristics:

- non-root `boredgamz` user (UID 10001);
- read-only root filesystem and temporary `/tmp` mount;
- all Linux capabilities dropped and `no-new-privileges` enabled;
- health check against `/readyz`;
- immutable one-year caching for hashed assets and no-cache SPA HTML;
- graceful SIGTERM shutdown with a ten-second deadline;
- production startup refusal when `ALLOWED_ORIGINS` is missing.

Terminate TLS at the platform or reverse proxy and forward WebSocket upgrades for:

- `/api/v1/arena/*`
- `/api/v1/poker`
- `/api/v1/mahjong`

`ALLOWED_ORIGINS` accepts a comma-separated list of exact origins without trailing slashes. Do not use wildcard origins with credentialed requests.

The live-room stores are intentionally in memory. Deploy one application replica. Horizontal scaling requires shared matchmaking/room ownership (for example Redis plus a durable pub/sub transport), not merely additional stateless replicas.

## Accounts and sessions

The production Compose stack includes PostgreSQL for account-backed display names and secure sessions. Production refuses to start without `DATABASE_URL`, without a 32-byte `JWT_SECRET_KEY`, or when the database/schema cannot be initialized. The database schema is embedded in the Go binary and applied idempotently during startup.

Authentication uses distinct access/refresh token claims, Secure/HttpOnly/SameSite cookies, POST-only refresh/logout, strict bounded JSON, generic credential errors, no-store responses, bcrypt passwords, and a per-client attempt limit. In production, terminate TLS before the service so its `Secure` cookies are accepted by browsers.

## HTTP and WebSocket surface

- `GET /healthz` — process liveness, uptime, optional database state
- `GET /readyz` — deployment readiness
- `GET /api/v1/catalog` — six-game catalog
- `GET /api/v1/status` — non-sensitive queue/room/connection counts
- `POST /api/v1/auth/signup` — create an account and issue a session
- `POST /api/v1/auth/login` — verify credentials and issue a session
- `POST /api/v1/auth/logout` — clear both session cookies
- `POST /api/v1/auth/refresh` — rotate the short-lived access cookie
- `GET /api/v1/auth/session` — return the current account identity
- `GET /api/v1/arena/{game}` — Tic-Tac-Toe, Connect Four, Gomoku, and Chess socket
- `GET /api/v1/poker` — Poker socket
- `GET /api/v1/mahjong` — Mahjong socket

All WebSocket handshakes validate the configured Origin. Connections have bounded messages, read/write deadlines, ping/pong liveness, bounded send queues, and generation-safe reconnect handling.

See [the production report](docs/PRODUCTION_REPORT.md) for the complete design rationale and browser verification record.
