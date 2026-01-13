package repository

import (
	"context"

	"github.com/eddiarnoldo/my-game-shelf/src/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BoardGameImageRepository struct {
	db *pgxpool.Pool
}

type BoardGameImageRepo interface {
	SaveImage(ctx context.Context, image *models.BoardGameImage) error
	GetAllImagesForBoardGame(ctx context.Context, boardGameId int64, imageType string) ([]*models.BoardGameImage, error)
	GetCoverThumbnail(ctx context.Context, boardGameId int64) (*models.BoardGameImage, error)
	DeleteImage(ctx context.Context, id int64) error
}

func NewBoardGameImageRepository(db *pgxpool.Pool) *BoardGameImageRepository {
	return &BoardGameImageRepository{db: db}
}

func (r *BoardGameImageRepository) SaveImage(ctx context.Context, image *models.BoardGameImage) error {
	query := `INSERT into board_game_images
	(board_game_id, image_data, image_mime_type, thumbnail_data, image_type, display_order, uploaded_at)
	VALUES ($1, $2, $3, $4, $5, $6, NOW()) RETURNING id, uploaded_at`

	err := r.db.QueryRow(ctx, query,
		image.BoardGameID,
		image.ImageData,
		image.ImageMimeType,
		image.ThumbnailData,
		image.ImageType,
		image.DisplayOrder,
	).Scan(&image.ID, &image.UploadedAt)

	return err
}

func (r *BoardGameImageRepository) GetAllImagesForBoardGame(ctx context.Context, boardGameId int64, imageType string) ([]*models.BoardGameImage, error) {
	query := `SELECT id, board_game_id, image_data, image_mime_type, thumbnail_data, image_type, display_order, uploaded_at
			FROM board_game_images`

	if imageType != "" {
		query += ` WHERE board_game_id = $1 AND image_type = $2 ORDER BY display_order ASC`
	} else {
		query += ` WHERE board_game_id = $1 ORDER BY display_order ASC`
	}

	var rows pgx.Rows
	var err error

	if imageType != "" {
		rows, err = r.db.Query(ctx, query, boardGameId, imageType)
	} else {
		rows, err = r.db.Query(ctx, query, boardGameId)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []*models.BoardGameImage
	for rows.Next() {
		var image models.BoardGameImage
		err := rows.Scan(&image.ID, &image.BoardGameID, &image.ImageData, &image.ImageMimeType,
			&image.ThumbnailData, &image.ImageType, &image.DisplayOrder, &image.UploadedAt)
		if err != nil {
			return nil, err
		}
		images = append(images, &image)
	}

	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return images, nil
}

// For list views - only get thumbnails
func (r *BoardGameImageRepository) GetCoverThumbnail(ctx context.Context, boardGameId int64) (*models.BoardGameImage, error) {
	query := `SELECT id, board_game_id, thumbnail_data, image_mime_type, image_type
			FROM board_game_images
			WHERE board_game_id = $1 AND image_type = 'cover'`

	var image models.BoardGameImage
	err := r.db.QueryRow(ctx, query, boardGameId).Scan(
		&image.ID,
		&image.BoardGameID,
		&image.ThumbnailData,
		&image.ImageMimeType,
		&image.ImageType,
	)

	if err != nil {
		return nil, err
	}

	return &image, nil
}

func (r *BoardGameImageRepository) DeleteImage(ctx context.Context, id int64) error {
	query := `DELETE FROM board_game_images WHERE id = $1`

	_, err := r.db.Exec(ctx, query, id)
	return err
}
