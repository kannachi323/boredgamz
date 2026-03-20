CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name varchar(255) NOT NULL,
    email varchar(255) NOT NULL UNIQUE,
    password TEXT NOT NULL,
    is_admin boolean NOT NULL DEFAULT false
);

CREATE TABLE IF NOT EXISTS games (
    id             UUID PRIMARY KEY,
    game_type      TEXT NOT NULL,
    player1_id     TEXT NOT NULL,
    player2_id     TEXT NOT NULL,
    winner         TEXT,
    finished_at    TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_at     TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at     TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    game_state     JSONB NOT NULL
    game_state     JSONB NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_games_type_created_at
    ON games (game_type, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_games_player1
    ON games (player1_id);

CREATE INDEX IF NOT EXISTS idx_games_player2
    ON games (player2_id);

CREATE TABLE IF NOT EXISTS messages (
    id           BIGSERIAL PRIMARY KEY,
    game_id      UUID NOT NULL REFERENCES games(id) ON DELETE CASCADE,
    game_type    TEXT NOT NULL,
    sender_id    TEXT NOT NULL,
    sender_name  TEXT NOT NULL,
    content      TEXT NOT NULL,
    sent_at      TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at   TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_messages_game_id_sent_at
    ON messages (game_id, sent_at);

CREATE INDEX IF NOT EXISTS idx_messages_game_type_sent_at
    ON messages (game_type, sent_at);
