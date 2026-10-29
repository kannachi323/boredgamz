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
Connect Five, Exact Five and Classical chess. Go tests cover their rules, and they use the same components as the verified modes. Also not covered: Lighthouse/axe audits, Safari, and real phones.

## Polish pass (second round)

Asked for: no gaps, minimal UI, no all-caps text, one clean font, consistent alignment and padding, production ready.

**UI**
- No uppercase anywhere: `text-transform: uppercase` and letter-spaced eyebrows are removed, and labels are sentence case.
- Private codes are six digits instead of letters, so no capitals appear and phones show a numeric keypad.
- One family (Geist) with a token type scale (`--fs-2xs`…`--fs-display`) and three weights. Every padding, gap and margin is on a 4 px spacing scale (`--s-0`…`--s-8`).
- **Mobile at 375 px.** Sign in and Sign up both fit in the top bar. Table-bar icons no longer collapse, game titles stay on one line, and players stack as compact rows. The code field spans the width.
- **Desktop at 1440 px.** The live-table count moved from the first tile to the header. The join row is balanced, stat labels no longer wrap, and the poker felt is a proper pill-shaped table.

**Gaps found and fixed**
1. Chess clocks ran from the moment of matching, so White lost time while the board loaded. Each side now gets a 30-second first move that doesn't touch the bank.
2. Joining a chess table by code showed "Classical 15+10" with a 10-minute bank. The mode now decides the clock on the server.
3. Resigning, agreeing a draw or disconnecting in chess didn't charge the running clock.
4. When the server was down, the app said "Your session ended". It now says it can't reach the server and keeps retrying.
5. Every signed-out page load logged two 401 errors in the console. `GET /session` now restores from the refresh cookie itself and answers 204 when signed out, so a normal load is error-free.
6. Signing in from an invite link lost the code. The top-bar links on the auth pages nested `next=` parameters.
7. Poker didn't say that stacks reset after someone busts.

**Verified in the browser this round:**
- Invite link → sign up → seated at the table
- Chess Rapid opening window and bank start, then resign
- Mahjong Mirror Duel (independent boards, Leave concedes) and Zen
- Tic-Tac-Toe Four Square with two-step Resign
- Poker fold and Next hand (dealer rotates)
- Wrong and right password, sign out
- 404 for unknown routes and games
- Console clean for signed-out and signed-in visitors
- No horizontal scroll at 375 px

Screenshots: `docs/screenshots/after/` shows the state before this pass (uppercase labels, letter codes). `docs/screenshots/polish/` shows the state after it.

## Server verification (third round)
- **Real-WebSocket tests for every service.** Poker checks private cards, turn enforcement and a fold. Mahjong checks seats, an invalid pair, and a hint followed by a claim. Chess checks an illegal and a legal move and the clock derived from the mode. Tic-Tac-Toe and Connect Four were already covered.
- **Fixed: closing a tab while waiting left a stale entry in the arena queue.** The lookup used the raw request instead of the normalised key. Nobody was ever paired with it, but `/status` kept counting it. Entries are now removed by player.
- **Fixed: abandoned best-of-N grid and chess rooms stayed in memory forever** and inflated the "tables live" count. A room is now removed 5 minutes after both players are gone. Poker and Mahjong already did this.
- **Production image smoke test with Postgres.** `/readyz` is ready, `/session` returns 204 then 200 after sign-up, and all six game sockets answer 101 when signed in and 401 when signed out. Queues return to 0 after disconnect, and no error logs.

## Queue fallback (owner approved)
- **Grid games:** the Tic-Tac-Toe, Connect Four and Gomoku public queues first match on exact settings. After 20 seconds, they pair anyone in the same game and mode, whatever series length or clock they chose.
- **Host settings:** whoever waited longer hosts, so their settings apply, and the table bar shows them.
- **Never mixed:** private codes never fall back, different modes never mix, and two tabs of one account never pair.
- **Waiting screen:** after 20 seconds it adds "Also matching any series length or clock in this mode".
- **Already effectively covered:** Chess (the mode fixes the clock) and Poker and Mahjong public queues (they match on mode alone).
- **Tests:** pairing after the delay, not before it, mode, private and same-account isolation, and both players receiving the new table. Verified in the browser: best of 3 with 30s paired with best of 5 with no clock.
- **Tuning:** the delay is `queueFallbackAfter` in `server/core/arena/types.go`.

## Decisions made without the owner
- **Sign-in is required to play.** The task asked for auth on every socket. Guests can browse the lobby and game pages. Guest sockets still work in local dev without a database.
- **Private codes are six digits.** The first digit is the game: 1 Tic-Tac-Toe, 2 Connect Four, 3 Gomoku, 4 Chess, 5 Poker, 6 Mahjong. There are 100,000 codes per game, which is fine for short-lived private tables. Custom codes are gone.
- **The chess mode fixes the clock**, and each side's first move has a 30-second window.
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
