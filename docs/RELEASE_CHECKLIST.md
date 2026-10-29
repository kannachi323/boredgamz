# Release checklist

Status as of 2026-10-28, branch `production-ready`. Nothing has been pushed or deployed.

## Done

### Triage
- The `temp shutdown` commit (`fffbc9c`) put the site behind a lockdown flag (`SITE_LOCKDOWN_ENABLED`) because of server costs. The owner's uncommitted work had already turned it off. That work is committed as-is in its own commit, and the lockdown component and flag are now removed.

### Server hardening
- **Auth on every socket.** With accounts enabled (always in production), a WebSocket handshake without a valid access cookie gets 401. Display names come from the account, never from the URL.
- **Origin checks.** WebSocket handshakes require an allowed `Origin`. Every POST requires an allowed `Origin`/`Referer` (CSRF), on top of `SameSite=Lax` cookies. Auth bodies must be `application/json`.
- **Input limits.** 4 KB socket messages, bounded query parameters, 8 KB auth bodies, unknown JSON fields rejected, chat capped at 240 characters, control characters stripped.
- **Rate limits.** 10/min on signup and login, 60/min on session refresh and socket handshakes, all per client IP. A chat flood guard (750 ms) runs per player. `X-Forwarded-For` is honoured only with `TRUST_PROXY=true`.
- **Passwords.** bcrypt cost 12. Length is 12–72 characters: bcrypt ignores bytes after 72 and the library errors on them, which previously made long passwords fail signup with a misleading "already registered". Unknown-email logins spend the same bcrypt time as real ones.
- **Cookies.** HttpOnly, Secure in production, SameSite=Lax. Access tokens last 1 hour and are refreshed every 20 minutes while a tab is open. Refresh tokens last 30 days. Access and refresh tokens have distinct claims.
- **Operations.** Versioned embedded migrations under an advisory lock. A bounded pgx pool (`DB_MAX_CONNS`). `/healthz` for liveness, and `/readyz` returning 503 while the DB is down. JSON `slog` logs with one line per request. Graceful SIGTERM shutdown.
- **Config.** Environment variables only. `.env` files are gitignored and loaded only outside production. No secrets in the repo.
- **Dead code removed.** The legacy Gomoku/Connect Four lobby packages, their DB queries, `/api/hello`, Tailwind, and unused pages and assets.
- **Dependencies.** Go modules upgraded per `govulncheck`. `npm audit` reports 0.

### Bugs found and fixed while playing
1. Signed-in Poker players were treated as spectators: seat lookup compared the client's local ID with the server's account ID. Snapshots now carry `viewerId`/`viewerSeat`.
2. Chat never marked your own messages, for the same reason.
3. A player who reloaded showed as "reconnecting" to the opponent forever. Rejoining now broadcasts state.
4. A socket that dropped right after being matched was treated as leaving the queue, leaving the opponent facing a ghost with no timeout.
5. A move sent immediately after the first state failed with "wait for an opponent" (race with the match poller). Also fixed a data race on the same channel.
6. After a game, creating or joining another table dropped you back into the finished room.
7. Reloading a tab returned you to the lobby instead of your seat. Seats now resume from `sessionStorage`. Leaving the page no longer forfeits; only the Leave button does.
8. A network error during session restore left the app spinning forever.
9. The opponent of a player who pressed Leave saw "reconnecting" and a Next round button that could never complete. They now see "Left the table".
10. Private codes only matched if both players picked identical options. A code now identifies the table, the host's settings win, and the code's first letter identifies the game, so it works from the home page.
11. Lobby tiles could stay invisible if the entrance animation was paused (the old `/games` page had the same symptom).

### Tests
- Go: rule engines for all six games (every grid mode, Connect Four in every direction, Connect Five, Exact Five, misère, chess mate/draws/clock, poker evaluation, chip conservation and private cards, mahjong solvability and claims). Also hubs, reconnect and leave semantics, security middleware, migrations, and a Postgres-backed signup → login → authenticated socket → CSRF → logout test. `go test -race ./...` passes.
- Web: typecheck, lint and Vitest (codes, safe redirects, socket endpoints, identity) pass. Production build passes.

### Browser verification
Two accounts in separate cookie jars (`localhost` and `player2.localhost`), desktop and 375 px:

| Area | Verified |
| --- | --- |
| Accounts | Sign up with redirect back to the game, session restore on reload, sign out, refresh |
| Tic-Tac-Toe | Private code from the home page, moves, chat both ways, reload mid-game with seamless resume, spectator via watch link (count updates, no controls), best-of-3 with start rotation, move-clock timeout, match result, rematch, Leave forfeits |
| Connect Four | Quick 5×6 public match, gravity drops, vertical win, rematch button |
| Gomoku | Compact 13×13, Walnut board, stones on intersections, five-in-a-row win, move replay |
| Chess | Blitz, legal-move hints, flipped board for Black, increment, Fool's Mate checkmate, rematch, draw offer and accept |
| Poker | Hold'em: private hole cards, call/check/raise slider, showdown with correct winner and chip totals. Five-Card Draw: discard selection, draw, stand pat, showdown. Leave → "Left the table" |
| Mahjong | Shared Garden: hint, pair removal synced to both players, per-player hints, timer, mobile layout |

Screenshots: `docs/screenshots/before/` and `docs/screenshots/after/`.

### Not verified in a browser
Mirror Duel and Zen Mahjong, Four Square, Connect Five, Exact Five, and the Rapid and Classical chess clocks. Server tests cover their rules, and the UI shares the same components. Also not covered: Lighthouse/axe audits, Safari, and real phones.

## Decisions made without the owner
- **Sign-in is required to play.** The task asked for auth on every socket. Guests can browse the lobby and game pages. Guest sockets still work in local dev without a database.
- **Private codes are six characters** with a game prefix (`T`, `C`, `G`, `H` for chess, `P`, `M`). Custom codes are gone.
- **The host's settings win** on a private table. Public queues still match on exact settings.
- **Navigating away keeps your seat** for the server's grace period. The explicit Leave button forfeits the round, folds the hand or concedes the duel.
- **A game in progress always resumes.** You can't escape a live game by asking for another table, but between rounds you can.
- **Design.** Warm neutrals with one accent (vermilion) and one typeface (Geist, self-hosted). The old serif and lime look and the photo and PNG game art were replaced by SVG/CSS drawings. Player one is ink or paper and player two is the accent, in every game.
- **Removed `docs/PRODUCTION_REPORT.md`.** It described the previous UI. Its useful facts are in the README and DEPLOY.
- **Legacy database tables** (`gomoku_games`, `connectfour_games`) are left in place, not dropped.

## Owner to-do before launch
1. **Hosting.** Pick a single-instance container host with WebSocket support: Fly.io, Render, Railway, or a small VPS with Caddy. See `docs/DEPLOY.md` for why there's exactly one instance.
2. **Postgres.** Provision a managed database (Neon, Supabase, RDS, Fly Postgres) or use the bundled compose service with backups. Use `sslmode=require` for managed DBs.
3. **Domain and TLS.** Point DNS at the host, enable HTTPS, and set `PUBLIC_ORIGIN`/`ALLOWED_ORIGINS` to the exact origin.
4. **Secrets.** Generate `JWT_SECRET_KEY` (`openssl rand -base64 48`) and a long DB password in the host's secret store. Set `TRUST_PROXY=true` if the platform overwrites `X-Forwarded-For`.
5. **Smoke test** on the real domain with two browsers: sign up, play one game, check `/readyz`.
6. **Optional.** Add an error tracker (Sentry or similar), uptime monitoring on `/readyz`, a privacy note (only email and display name are stored), and a way to delete accounts.
7. **Optional.** Run Lighthouse and axe on the deployed site, and test on Safari/iOS.
8. Merge `production-ready` into `main` when you're happy, then push.
