# BoredGamz production design and verification report

## Executive summary

The finished product is a six-table online game room, not one generic board component wearing six labels. The repository’s original strengths—warm tactile surfaces, existing game artwork, game-oriented folders, and server-side validation—were retained. Incomplete stores, duplicated legacy lobby UI, solo placeholders, and client-authoritative flows were replaced with modular production implementations.

The common contract is consistent: real player accounts, exact mode matchmaking, optional private codes, two seated humans, chat, observers, reconnect grace, authoritative timers/rules, and rematches. The visual treatment is intentionally not uniform. Each game has its own lobby composition, table hierarchy, motion language, controls, settings, and thematic materials.

The resulting TypeScript/CSS/Go surface is about 17,400 lines across 156 source files. No source file exceeds 1,000 lines; the largest is an 828-line table stylesheet. That size is a consequence of six complete experiences and three server game systems, not repeated filler.

## Design system choices

### Shared foundation

The foundation returns to the original BoredGamz character:

- warm near-black `#1b1918` instead of blue-black SaaS chrome;
- layered charcoal/wood panels (`#302e2e`, `#433d3a`);
- parchment copy (`#c3b299` and lighter variants);
- restrained mint (`#7dcfb6`) for live/positive actions;
- Georgia display type paired with a compact system sans-serif UI face;
- tactile borders, inset shadows, grain-like gradients, and physical board depth.

The homepage and catalog remain coherent, while immersive game routes remove the global navigation/footer so the table can use the viewport. A skip link, semantic landmarks, visible focus states, button state semantics, readable live status, and `prefers-reduced-motion` fallbacks are shared accessibility decisions.

### Motion

Animation communicates state rather than running as decoration:

- table and board entrance sequences;
- piece drops, stone placement, card dealing, and tile removal;
- active-player glows and low-clock pulses;
- legal-move/hint reveals;
- lobby artwork drift and matchmaking loops;
- dialog/settings entrance transitions;
- all significant motion disabled for reduced-motion users.

### Asset policy

Existing repository art was retained for Connect Four, Chess, Poker, and Mahjong. The Tic-Tac-Toe cover was generated specifically for BoredGamz. Gomoku’s board textures were derived as local maple, walnut, and oak WebP assets. Runtime game pieces use CSS, Unicode, and locally bundled images; the game pages make no third-party asset requests and introduce no copyrighted franchise art.

## Game-by-game decisions

### Tic-Tac-Toe

Tic-Tac-Toe is presented as a fast editorial “quick table,” using coral marks, parchment panels, and compact motion.

- Classic: 3×3, connect three.
- Four Square: 4×4, connect three, creating more fork space.
- Misère: completing three loses.
- Best-of 1, 3, 5, or 7.
- Untimed, 15-, 30-, or 60-second move clocks.

The server owns placement, turn order, variant winner inversion, clocks, scores, starting-player rotation, and series completion.

### Connect Four

Connect Four is brighter and more mechanical: yellow/red discs, deep navy housing, drop-zone motion, and a vertical scoreboard.

- Classic Drop: 6×7, connect four.
- Quick Stack: 5×6, connect four.
- Long Game: 7×8, connect five.
- Best-of 1, 3, 5, or 7 with the same optional move clocks.

Clients submit columns. The server computes the landing row, rejects full columns, detects connections in every direction, and controls the series.

### Gomoku

Gomoku uses a quieter club aesthetic with locally bundled wood grains and a board designed around intersections rather than square cells.

- Freestyle: 15×15, five or more wins.
- Exact Five: 15×15, overlines do not win.
- Compact: 13×13 for faster pressure.
- Single game or best-of 3/5; untimed or 30/60-second moves.
- Pale maple, club walnut, and deep oak personal board themes.
- A server-backed move ledger with previous/next/live replay controls after moves begin.

The server enforces bounds, occupancy, exact-five semantics, turn clocks, series state, and rematches. Every accepted move is appended to the room snapshot, so replay is reconstructed from authoritative history rather than a client-only log. The empty history is serialized as `[]`, keeping the wire shape stable from the first snapshot onward.

### Accounts

Signup and login are first-class product flows again. PostgreSQL stores normalized email addresses, bcrypt password hashes, account IDs, and display names. Successful signup/login returns the public account identity and issues separate short-lived access and long-lived refresh cookies. The client checks the session before rendering routes, refreshes an expired access token once, and shows authenticated identity plus sign-out in the global navigation.

Live-room identity is derived from the verified access cookie when present. A tab-scoped suffix preserves independent reconnect ownership without allowing a browser-provided query string to impersonate another account. Guest fallback remains available for low-friction local development when the account service is intentionally not configured. Production refuses to start without PostgreSQL, a strong JWT secret, or a successfully applied schema, so a deployment cannot silently lose signup/login.

### Chess

Chess was designed as a calm private club rather than a tournament dashboard. The ledger, clocks, board, controls, and conversation each have a clear visual zone.

- Blitz 3+2, Rapid 10, and Classical 15+10.
- Four persistent local board palettes.
- Legal destination guidance and last-move indication.
- FEN and SAN history from the authoritative Go engine.
- Castling, en passant, promotion, check/checkmate, stalemate, repetition, insufficient material, resignations, and agreed draws.
- Per-side authoritative clocks with increment handling.

The browser was tested through Scholar’s Mate, verified SAN/FEN and orientation for both seats, exercised chat, completed an agreed draw, started a fresh rematch with White correctly moving first, and rendered without overflow at 390×844.

### Poker

Poker has its own dark-green art-deco card room. It is explicitly free-play and heads-up; the interface never implies cash value.

- No-Limit Texas Hold’em with 10/20 blinds and 1,000-chip stacks.
- Five-Card Draw with a 10-chip ante and up to four replacements.
- OS cryptographic randomness for deck shuffling.
- Burn cards, betting rounds, check/call/fold/raise/all-in, minimum raises, unmatched-chip refunds, showdown evaluation, split pots, and stack reset when a player busts.
- Thirty-second authoritative action timer.
- Viewer-specific snapshots: a seat receives only its own private cards until showdown.

Browser verification completed a Hold’em hand through flop, turn, river, and a two-pair showdown; completed a Draw hand with three replacements versus a stand-pat opponent; checked chip conservation, chat, privacy from both DOMs, timer auto-action, rematch, and the 390×844 layout.

### Mahjong Solitaire

The repository’s original identity was Mahjong Solitaire, so the product did not silently become four-player wagering Mahjong. Multiplayer was added in ways that fit the existing game:

- Shared Garden: cooperative 144-tile clear in 12 minutes.
- Mirror Duel: two identical 72-tile layouts; first clear wins, or fewer remaining at eight minutes.
- Zen Partners: untimed shared 72-tile board.

The tile set includes the standard 144 instances: three numbered suits, winds, dragons, seasons, and flowers. Seasons match any season; flowers match any flower. A tile is free only when uncovered and open on its left or right edge.

The server generates a legal geometry-removal order before assigning shuffled face pairs, which proves every initial board has a solution. If player choices create a dead end, only remaining faces are cryptographically reshuffled onto a new provably removable order; already-cleared tiles stay cleared. Hints are authoritative, private, and limited.

Two-browser verification covered shared claims, a simultaneous race for the same pair (one accepted, one rejected), private hints, synchronized chat, independent mirrored progress, a complete 36-pair Duel clear through the UI, both-player rematch readiness, untimed Zen play, read-only spectators, hot-reload/reconnect behavior, four themes, and a 144-tile board measuring 374px inside a 390px viewport.

## Multiplayer architecture

### Arena service

Tic-Tac-Toe, Connect Four, Gomoku, and Chess share the Arena transport but not a client monolith. The Arena has:

1. a hub keyed by game, mode, series, clock, and private code;
2. rooms that serialize all state changes under a mutex;
3. parameterized grid rules plus a dedicated Chess engine;
4. personalized snapshots carrying a server-authoritative viewer role;
5. bounded channels, messages, chat history, and name lengths;
6. 20-second generation-safe reconnect grace and five-minute empty-room cleanup.

### Poker and Mahjong services

Poker and Mahjong have separate packages because their state and privacy models do not fit a grid-room abstraction. Each has its own protocol, hub, room, socket loop, and rule engine. This avoids a “universal game state” full of optional fields and keeps security decisions—private hole cards or per-seat mirrored boards—close to the relevant game.

### Presence and spectators

Spectator socket IDs are deliberately namespaced from seated browser IDs, preventing a watch tab from replacing a player’s reconnect mapping. Arena and Mahjong snapshots also carry the authoritative viewer seat/mark so the UI never infers permission from local storage. Spectators can chat but cannot move pieces, wager, remove tiles, request hints, or ready a rematch.

### Scope for deployment

Room state is in memory. This is a production-ready single-replica architecture: simple failure modes, no hidden broker dependency, and low latency. Multiple replicas require external matchmaking and room ownership plus pub/sub or actor routing. The documentation does not claim horizontal scaling that the code does not provide.

## Security decisions

- Exact production Origin allowlist for WebSocket handshakes and credentialed CORS.
- Production refuses to start without `ALLOWED_ORIGINS`.
- CSP, frame denial, MIME-sniff prevention, referrer policy, and browser permission policy on all responses.
- Client messages express intent only; clients never submit computed winners, legal-move sets, hands, pots, or tile availability.
- Message size caps, read/write deadlines, ping/pong liveness, bounded queues, and sanitized display/chat text.
- CSPRNG shuffling for Poker and Mahjong.
- No hole-card leakage before showdown.
- No secret invite codes or player content in the runtime status endpoint.
- Non-root, read-only, capability-dropped container with `no-new-privileges`.
- Accounts use bcrypt, distinct access/refresh JWT claims, exact HS256 verification, issuer/audience validation, 32-byte minimum secret, Secure/HttpOnly/SameSite cookies, bounded strict JSON, POST-only mutation routes, no-store responses, generic credential errors, and an in-memory per-client rate limit.
- `pgx` upgraded to 5.9.2 after vulnerability scanning identified the older reachable version.
- Go 1.26.5 is pinned in CI and Docker to include the standard-library fixes named by the current vulnerability database.

## Operations and observability

- `/healthz` reports process liveness, uptime, and whether optional Postgres is connected.
- `/readyz` is the orchestration readiness check.
- `/api/v1/status` reports non-sensitive waiting/room/connection counts for Arena, Poker, and Mahjong.
- SIGINT/SIGTERM trigger a ten-second graceful HTTP shutdown.
- HTTP server read-header, read, write, and idle timeouts are explicit.
- Hashed static assets use immutable one-year caching; SPA HTML uses `no-cache`.
- CI runs npm install/lint/tests/build/audit, Go race tests/vet, then builds and health-checks the container.

## Verification evidence

### Automated

- 15 frontend tests.
- 70 Go tests under the race detector.
- 24 concurrent private WebSocket matches, each connecting two clients, exchanging synchronized snapshots, playing a legal winning sequence, and agreeing on completion.
- Go vet clean.
- TypeScript production build and ESLint clean.
- npm audit: zero vulnerabilities.
- govulncheck under Go 1.26.5: zero reachable vulnerabilities.
- `git diff --check`: clean.
- Every TypeScript, TSX, CSS, and Go source file below 1,000 lines.

### Browser

- Real isolated-origin two-client sessions for all multiplayer systems.
- Two accounts were created through the signup UI (`Ada Stone` and `Grace Grid`), used simultaneously from isolated cookie jars, signed out, rejected with an incorrect password, signed back in, and retained across reload.
- Tic-Tac-Toe finished with a three-in-a-row result; Connect Four finished with a vertical connection; Gomoku finished with a five-stone line and its nine-move replay was navigated backward, forward, and back to live; Chess finished by checkmate after `f3 e5 g4 Qh4#`.
- Hold'em finished after both players acted preflop and on the flop, with the uncontested pot awarded correctly and private hole cards hidden from the opponent; Zen Mahjong finished after both accounts cleared 18 pairs each from the shared 72-tile layout.
- Full Chess, Hold’em, Five-Card Draw, and Mahjong flows described above.
- Invalid turn/action, occupied move, contested tile, reconnect, chat, draw/rematch, privacy, and timeout edge cases.
- Desktop sweep of homepage, catalog, and all six routes: expected heading, no horizontal overflow, no console warning/error.
- Phone sweep at 390×844 for homepage, catalog, Chess, Poker, and Mahjong: no horizontal overflow or console warning/error.

### Production container

The final multi-stage image was built with a fresh Go 1.26.5 builder and launched with the Compose restrictions. It reported:

- user `boredgamz`;
- read-only root filesystem;
- `ALL` capabilities dropped;
- `no-new-privileges` active;
- healthy readiness status;
- correct catalog and runtime counters;
- SPA fallback on `/games/chess`;
- `no-cache` for HTML and immutable caching for hashed CSS;
- expected CSP and security headers;
- intentional startup failure when production `ALLOWED_ORIGINS` was omitted.

## Remaining operational choices

These are deployment-owner choices, not unfinished game behavior:

- choose the public domain/TLS termination platform;
- set `PUBLIC_ORIGIN` and DNS;
- provision and back up the required PostgreSQL account database;
- add an external room transport before running more than one application replica;
- choose log/metric ingestion for the platform, using `/api/v1/status`, `/healthz`, and `/readyz` as the base signals.
