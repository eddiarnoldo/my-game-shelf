CREATE TABLE board_games (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    min_players INTEGER NOT NULL,
    max_players INTEGER,
    play_time INTEGER,
    min_age INTEGER,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT check_min_players CHECK (min_players > 0),
    CONSTRAINT check_max_players CHECK (max_players IS NULL OR max_players >= min_players),
    CONSTRAINT check_play_time CHECK (play_time > 0)
);

-- Index for name searches/sorting
CREATE INDEX idx_board_games_name ON board_games(name);