board_game 

id seq
name string
min_users
max_users
play_time_minutes
description

```
CREATE TABLE board_game (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    min_users INTEGER NOT NULL,
    max_users INTEGER,
    play_time_minutes INTEGER,
    description TEXT,
    CONSTRAINT check_min_users CHECK (min_users > 0),
    CONSTRAINT check_max_users CHECK (max_users IS NULL OR max_users >= min_users),
    CONSTRAINT check_play_time CHECK (play_time_minutes > 0)
);
```