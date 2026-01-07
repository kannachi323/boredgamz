# boredgamz

A scalable, real-time web platform for multiplayer games. Play classic games like Gomoku, Connect Four, and Tic-Tac-Toe with a custom networking engine optimized for low-latency state synchronization and concurrent player management.

## About

**boredgamz** is built with a modular architecture that separates core game logic from individual game implementations. The platform features:

- **Real-time Multiplayer**: WebSocket-based communication for instant game state synchronization
- **Lobby & Matchmaking System**: Automated player management, room creation, and queue handling
- **Scalable Architecture**: Built with concurrency and extensibility in mind—add new games without touching core infrastructure
- **Modern UI**: Responsive, intuitive game interfaces built with React and Tailwind CSS
- **Persistent Storage**: PostgreSQL database for player data, game history, and statistics

### Current Games

- **Gomoku** (Five in a Row)
- **Connect Four**
- **Tic-Tac-Toe**

## Getting Started

### Prerequisites

- **Node.js** 18+ and npm
- **Go** 1.23+
- **Docker** and Docker Compose
- **PostgreSQL** (via Docker, or standalone)

### Installation

1. **Clone the repository**

   ```bash
   git clone https://github.com/mtcco/boredgamz.git
   cd boredgamz
   ```

2. **Start the database**

   ```bash
   cd service
   docker-compose up -d
   ```

   This starts PostgreSQL containers for development and testing on ports 5433 and 5434.

3. **Set up environment variables**

   Create a `.env` file in the `service/` directory:

   ```env
   POSTGRES_DB=boredgamz
   POSTGRES_USER=postgres
   POSTGRES_PASSWORD=yourpassword
   POSTGRES_DB_TEST=boredgamz_test
   POSTGRES_USER_TEST=postgres
   POSTGRES_PASSWORD_TEST=testpassword
   ```

4. **Start the backend server**

   ```bash
   cd service
   go mod tidy
   go run main.go
   ```

   The server runs on `http://localhost:3000`

5. **Start the frontend development server**

   ```bash
   cd app
   npm install
   npm run dev
   ```

   The frontend runs on `http://localhost:5173`

### Development Workflow

**Backend (Go)**
- API endpoints in `service/api/`
- Core game logic in `service/core/`
- Database queries in `service/db/`
- WebSocket upgrades and utility functions in `service/utils/`

**Frontend (React + TypeScript)**
- Game components in `app/src/pages/Games/`
- Shared components in `app/src/components/`
- State management (Zustand) in `app/src/stores/`
- Routing in `app/src/router.tsx`

### Building for Production

**Frontend**
```bash
cd app
npm run build
```

**Backend**
```bash
cd service
go build -o boredgamz main.go
```

## Contributing

We welcome contributions! Please follow these steps:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-game`)
3. Commit your changes (`git commit -m 'Add amazing game'`)
4. Push to the branch (`git push origin feature/amazing-game`)
5. Open a Pull Request

### Adding a New Game

To add a new game to the platform:

1. **Define game types** in `service/core/yourGame/types.go`
2. **Implement game logic** in `service/core/yourGame/game.go`
3. **Create a lobby handler** in `service/core/yourGame/lobby.go`
4. **Add WebSocket message handlers** in `service/api/yourGame/`
5. **Register the lobby** in `service/server/server.go`
6. **Build frontend components** in `app/src/pages/Games/YourGame/`

For detailed architectural guidance, see `service/core/README.md`.

## Tech Stack

**Frontend**
- React 18
- TypeScript
- Vite
- Tailwind CSS
- Zustand (state management)
- React Router

**Backend**
- Go 1.23
- Chi (HTTP router)
- Gorilla WebSocket
- PostgreSQL
- JWT (authentication)

**Infrastructure**
- Docker & Docker Compose
- PostgreSQL

## Support

For questions, issues, or feature requests:

- Open an [issue](../../issues) on GitHub
- Check existing documentation in the `service/core/README.md`
- Review game-specific code in `service/core/{game_name}/`

## License

This project is open source. See the LICENSE file for details.

## Maintainers

**boredgamz** is maintained by the community. Contributions are welcome!
