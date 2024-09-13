package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/xorwise/music-streaming-service/internal/domain"
)

type trackRepository struct {
	db *sql.DB
}

func NewTrackRepository(db *sql.DB) domain.TrackRepository {
	return &trackRepository{
		db: db,
	}
}

func (tr *trackRepository) Create(ctx context.Context, track *domain.Track) (int64, error) {
	stmt, err := tr.db.PrepareContext(ctx, "INSERT INTO tracks (title, artist, path, room_id) VALUES ($1, $2, $3, $4) RETURNING id")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	var id int64
	err = stmt.QueryRowContext(ctx, track.Title, track.Artist, track.Path, track.RoomID).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (tr *trackRepository) Remove(ctx context.Context, trackID int64) error {
	stmt, err := tr.db.PrepareContext(ctx, "DELETE FROM tracks WHERE id = $1")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, trackID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrTrackNotFound
		}
		return err
	}
	return nil
}

func (tr *trackRepository) Update(ctx context.Context, track *domain.Track) error {
	stmt, err := tr.db.PrepareContext(ctx, "UPDATE tracks SET title = $1, artist = $2, path = $3, room_id = $4, is_ready = $5 WHERE id = $6")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, track.Title, track.Artist, track.Path, track.RoomID, track.IsReady, track.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrTrackNotFound
		}
		return err
	}
	return nil
}

func (tr *trackRepository) GetByID(ctx context.Context, trackID int64) (*domain.Track, error) {
	var track *domain.Track = &domain.Track{}
	stmt, err := tr.db.PrepareContext(ctx, "SELECT id, title, artist, path, room_id, is_ready FROM tracks WHERE id = $1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	row := stmt.QueryRowContext(ctx, trackID)
	if err := row.Scan(&track.ID, &track.Title, &track.Artist, &track.Path, &track.RoomID, &track.IsReady); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrTrackNotFound
		}
		return nil, err
	}
	return track, nil
}
