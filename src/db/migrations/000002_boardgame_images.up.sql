-- Single table for all images with type differentiation
CREATE TABLE board_game_images (
    id SERIAL PRIMARY KEY,
    board_game_id INTEGER NOT NULL REFERENCES board_games(id) ON DELETE CASCADE,
    image_data BYTEA NOT NULL,
    image_mime_type VARCHAR(50) NOT NULL,
    thumbnail_data BYTEA,
    image_type VARCHAR(20) NOT NULL, -- 'cover' or 'gameplay'
    display_order INTEGER DEFAULT 0,
    uploaded_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_board_game FOREIGN KEY (board_game_id) REFERENCES board_games(id),
    CONSTRAINT check_image_type CHECK (image_type IN ('cover', 'gameplay'))
);

-- Index for quickly getting cover images in list views
CREATE INDEX idx_board_game_images_game_id ON board_game_images(board_game_id);
CREATE INDEX idx_board_game_images_cover ON board_game_images(board_game_id, image_type) WHERE image_type = 'cover';

-- Ensure only one cover image per game
CREATE UNIQUE INDEX idx_one_cover_per_game ON board_game_images(board_game_id) WHERE image_type = 'cover';