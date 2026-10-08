package db

import (
	"database/sql"
	"fmt"

	"github.com/atbu/fossil/records"
	"github.com/atbu/fossil/utils"

	_ "modernc.org/sqlite" // Pure Go driver (no CGO needed)
)

func InitDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if _, err := db.Exec(`PRAGMA journal_mode = WAL;`); err != nil {
		return nil, fmt.Errorf("failed to enable WAL mode: %w", err)
	}

	schema := `
	CREATE TABLE IF NOT EXISTS audio_records (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		ts DATETIME,
		platform TEXT NOT NULL,
		ms_played INTEGER NOT NULL,
		conn_country TEXT NOT NULL,
		ip_addr TEXT NOT NULL,
		master_metadata_track_name TEXT NOT NULL,
		master_metadata_album_artist_name TEXT NOT NULL,
		master_metadata_album_album_name TEXT NOT NULL,
		spotify_track_uri TEXT NOT NULL,
		episode_name TEXT NOT NULL,
		episode_show_name TEXT NOT NULL,
		spotify_episode_uri TEXT NOT NULL,
		audiobook_title TEXT NOT NULL,
		audiobook_uri TEXT NOT NULL,
		audiobook_chapter_uri TEXT NOT NULL,
		audiobook_chapter_title TEXT NOT NULL,
		reason_start TEXT NOT NULL,
		reason_end TEXT NOT NULL,
		shuffle INTEGER NOT NULL,
		skipped INTEGER NOT NULL,
		offline INTEGER NOT NULL,
		offline_timestamp INTEGER,
		incognito_mode INTEGER
	);

	CREATE INDEX IF NOT EXISTS idx_track_artist ON audio_records(master_metadata_track_name, master_metadata_album_artist_name);
	`

	if _, err := db.Exec(schema); err != nil {
		return nil, fmt.Errorf("failed to create schema: %w", err)
	}

	return db, nil
}

func BulkInsertAudioRecords(db *sql.DB, audioRecords []records.AudioRecord) (int64, int64, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT OR IGNORE INTO audio_records (
			ts,
			platform,
			ms_played,
			conn_country,
			ip_addr,
			master_metadata_track_name,
			master_metadata_album_artist_name,
			master_metadata_album_album_name,
			spotify_track_uri,
			episode_name,
			episode_show_name,
			spotify_episode_uri,
			audiobook_title,
			audiobook_uri,
			audiobook_chapter_uri,
			audiobook_chapter_title,
			reason_start,
			reason_end,
			shuffle,
			skipped,
			offline,
			offline_timestamp,
			incognito_mode
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	var attemptedRows int64

	for _, r := range audioRecords {
		if r.MasterMetadataTrackName == "" || r.MasterMetadataAlbumArtistName == "" {
			continue
		}

		attemptedRows++

		_, err := stmt.Exec(
			r.Ts,
			r.Platform,
			r.MsPlayed,
			r.ConnCountry,
			r.IPAddr,
			r.MasterMetadataTrackName,
			r.MasterMetadataAlbumArtistName,
			r.MasterMetadataAlbumAlbumName,
			r.SpotifyTrackURI,
			r.EpisodeName,
			r.EpisodeShowName,
			r.SpotifyEpisodeURI,
			r.AudiobookTitle,
			r.AudiobookURI,
			r.AudiobookChapterURI,
			r.AudiobookChapterTitle,
			r.ReasonStart,
			r.ReasonEnd,
			utils.BoolToInt(r.Shuffle),
			utils.BoolToInt(r.Skipped),
			utils.BoolToInt(r.Offline),
			r.OfflineTimestamp,
			utils.BoolToInt(r.IncognitoMode),
		)
		if err != nil {
			return 0, 0, fmt.Errorf("failed to execute insert: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	var totalInserted int64
	err = db.QueryRow("SELECT total_changes()").Scan(&totalInserted)
	if err != nil {
		totalInserted = attemptedRows // fallback if total_changes isn't available
	}

	skippedRows := attemptedRows - totalInserted
	if skippedRows < 0 {
		skippedRows = 0
	}

	return totalInserted, skippedRows, nil
}

func GetAudioRecordsByYear(db *sql.DB, year int) ([]records.AudioRecord, error) {
	startBoundary := fmt.Sprintf("%d-01-01 00:00:00", year)
	endBoundary := fmt.Sprintf("%d-01-01 00:00:00", year+1)

	query := `
	SELECT
		ts,
		ms_played,
		master_metadata_track_name,
		master_metadata_album_artist_name
	FROM audio_records
	WHERE ts >= ? AND ts < ?;
	`

	rows, err := db.Query(query, startBoundary, endBoundary)
	if err != nil {
		return nil, fmt.Errorf("failed to query records for year %d: %w", year, err)
	}
	defer rows.Close()

	var audioRecords []records.AudioRecord

	for rows.Next() {
		var record records.AudioRecord

		var ts sql.NullTime
		var trackName, artistName sql.NullString

		err := rows.Scan(
			&ts,
			&record.MsPlayed,
			&trackName,
			&artistName,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		if ts.Valid {
			record.Ts = &ts.Time
		}
		record.MasterMetadataTrackName = trackName.String
		record.MasterMetadataAlbumArtistName = artistName.String

		audioRecords = append(audioRecords, record)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return audioRecords, nil
}
