package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/xorwise/music-streaming-service/internal/domain"
)

type trackRepository struct {
	db      *sql.DB
	trackCh chan domain.TrackStatus
}

func NewTrackRepository(db *sql.DB, trackCh chan domain.TrackStatus) domain.TrackRepository {
	return &trackRepository{
		db:      db,
		trackCh: trackCh,
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

func (tr *trackRepository) Remove(ctx context.Context, track *domain.Track) error {
	stmt, err := tr.db.PrepareContext(ctx, "DELETE FROM tracks WHERE id = $1")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, track.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrTrackNotFound
		}
		return err
	}

	tr.trackCh <- domain.TrackStatus{
		ID:      track.ID,
		RoomID:  track.RoomID,
		IsReady: false,
	}
	return nil
}

func (tr *trackRepository) Update(ctx context.Context, track *domain.Track) error {
	oldTrack := domain.Track{}
	stmt, err := tr.db.PrepareContext(ctx, "SELECT id, title, artist, path, room_id, is_ready FROM tracks WHERE id = $1")
	if err != nil {
		return err
	}
	defer stmt.Close()
	row := stmt.QueryRowContext(ctx, track.ID)
	err = row.Scan(&oldTrack.ID, &oldTrack.Title, &oldTrack.Artist, &oldTrack.Path, &oldTrack.RoomID, &oldTrack.IsReady)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrTrackNotFound
		}
		return err
	}
	if track.IsReady && !oldTrack.IsReady {
		tr.trackCh <- domain.TrackStatus{
			ID:      track.ID,
			RoomID:  track.RoomID,
			Path:    track.Path,
			IsReady: track.IsReady,
		}
	}

	newStmt, err := tr.db.PrepareContext(ctx, "UPDATE tracks SET title = $1, artist = $2, path = $3, room_id = $4, is_ready = $5 WHERE id = $6")
	if err != nil {
		return err
	}
	defer newStmt.Close()
	_, err = newStmt.ExecContext(ctx, track.Title, track.Artist, track.Path, track.RoomID, track.IsReady, track.ID)
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

func (tr *trackRepository) ListByRoomID(ctx context.Context, roomID int64, params url.Values) ([]*domain.Track, error) {
	var tracks []*domain.Track
	query := generateSQLQuery(params)
	stmt, err := tr.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.QueryContext(ctx, roomID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		track := domain.Track{}
		if err := rows.Scan(&track.ID, &track.Title, &track.Artist, &track.Path, &track.RoomID, &track.IsReady); err != nil {
			return nil, err
		}
		tracks = append(tracks, &track)
	}
	return tracks, nil
}

func generateSQLQuery(params url.Values) string {
	baseQuery := "SELECT id, title, artist, path, room_id, is_ready FROM tracks WHERE room_id = $1"
	conditions := []string{}

	for key, value := range params {
		if len(value) > 0 {
			condition := ""
			switch key {
			case "title":
				condition = fmt.Sprintf("title LIKE '%%%s%%'", value[0])
			case "artist":
				condition = fmt.Sprintf("artist LIKE '%%%s%%'", value[0])
			case "is_ready":
				if value[0] == "true" {
					condition = "is_ready = true"
				} else {
					condition = "is_ready = false"
				}
			case "limit", "offset":
				continue
			}
			conditions = append(conditions, condition)
		}
	}

	limit := params.Get("limit")
	offset := params.Get("offset")

	if len(conditions) > 0 {
		baseQuery += " AND " + strings.Join(conditions, " AND ")
	} else {
		baseQuery = "SELECT id, title, artist, path, room_id, is_ready FROM tracks WHERE room_id = $1"
	}

	if limit == "" {
		baseQuery += " LIMIT 100"
	} else {
		baseQuery += fmt.Sprintf(" LIMIT %s", limit)
	}

	if offset == "" {
		baseQuery += " OFFSET 0"
	} else {
		baseQuery += fmt.Sprintf(" OFFSET %s", offset)
	}

	fmt.Println(baseQuery)
	return baseQuery
}
