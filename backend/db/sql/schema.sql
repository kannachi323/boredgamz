CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name varchar(255) NOT NULL,
    email varchar(255) NOT NULL UNIQUE,
    password TEXT NOT NULL,
    is_admin boolean NOT NULL DEFAULT false
);

CREATE TABLE IF NOT EXISTS gomoku_games (
    id             UUID PRIMARY KEY,
    player1_id     TEXT NOT NULL,
    player2_id     TEXT NOT NULL,
    winner         TEXT,
    finished_at    TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_at     TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at     TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    game_state     JSONB NOT NULL
);
