# BoredGamz

Six classic tabletop games, played live in the browser: Tic-Tac-Toe, Connect Four, Gomoku, Chess, free-play Poker and Mahjong Solitaire. Sign in, pick a table, and play a stranger or a friend with a private code. Spectators can watch and everyone can chat. If you reload, you rejoin your seat. The server enforces every rule.

![Lobby](docs/screenshots/after/home-desktop-light.jpg)

| Game | Modes | Options |
| --- | --- | --- |
| Tic-Tac-Toe | Classic 3×3, Four Square 4×4, Misère | Best of 1/3/5/7, move clock |
| Connect Four | Classic 6×7, Quick 5×6, Connect Five 7×8 | Best of 1/3/5/7, move clock |
| Gomoku | Freestyle 15×15, Exact Five, Compact 13×13 | Best of 1/3/5, move clock, board style |
| Chess | Blitz 3+2, Rapid 10, Classical 15+10 | Draw offers, board style |
| Poker | Heads-up No-Limit Hold'em, Five-Card Draw | Free-play chips only, with no purchases or cash value |
| Mahjong Solitaire | Shared Garden (co-op), Mirror Duel (race), Zen | Hints, guaranteed-solvable layouts |

## Stack

- **Web** (`app/`): React 18 + TypeScript + Vite, with one hand-written design system (`src/styles/`) and SVG board art.
- **Server** (`server/`): Go with chi and gorilla/websocket. It owns matchmaking, rooms, clocks and rules (`core/arena`, `core/poker`, `core/mahjong`).
- **Database**: PostgreSQL for accounts, with embedded migrations.
- **Auth**: bcrypt passwords, HttpOnly/Secure/SameSite cookies, short-lived access tokens with refresh, Origin-checked writes and sockets.

## Run locally

Requirements: Node 22+, Go 1.26.5+ and Docker (for Postgres).

```bash
docker compose -f server/docker-compose.yml up -d db
```

```bash
cd server && DATABASE_URL='postgres://boredgamz:boredgamz-dev@127.0.0.1:5433/boredgamz?sslmode=disable' JWT_SECRET_KEY=dev-only-secret-at-least-32-bytes-long PORT=8201 go run .
```

```bash
cd app && npm ci && npm run dev
```

Open http://localhost:5201. Vite proxies `/api` (HTTP and WebSocket) to `:8201`, so cookies behave as they do in production. To play yourself, open a second browser profile, or use `http://player2.localhost:5201` (add that origin to `ALLOWED_ORIGINS`).

## Test

```bash
cd app && npm run typecheck && npm run lint && npm test && npm run build
```

```bash
cd server && go vet ./... && go test -race ./...
```

Set `TEST_DATABASE_URL` to also run the Postgres-backed auth and socket integration test, as CI does.

## Deploy

`docker compose up --build -d` runs the app and Postgres. See [docs/DEPLOY.md](docs/DEPLOY.md) for environment variables, migrations, TLS and proxy setup, and why the app runs as a single instance. [docs/RELEASE_CHECKLIST.md](docs/RELEASE_CHECKLIST.md) lists what's done and what's left before launch.
