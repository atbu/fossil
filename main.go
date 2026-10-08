package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/atbu/fossil/song"
)

func main() {
	filePath := flag.String("file", "", "Path to the Spotify audio streaming history file")
	flag.Parse()

	if _, err := os.Stat(*filePath); os.IsNotExist(err) {
		fmt.Printf("Error: File %s does not exist.\n", *filePath)
		os.Exit(1)
	}

	songs, err := parseAudioStreamingHistoryFile(*filePath)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Successfully loaded %d records.\n", len(songs))

	leaderboard := generateSongLeaderboard(songs)

	for index, entry := range leaderboard {
		duration := entry.MsPlayed.Round(time.Second)
		fmt.Printf("%d. %s - %s: %s\n", index, entry.TrackName, entry.ArtistName, duration)
	}
}

func parseAudioStreamingHistoryFile(filePath string) ([]song.Song, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	t, err := decoder.Token()
	if err != nil {
		return nil, fmt.Errorf("failed to read start token: %w", err)
	}
	if delim, ok := t.(json.Delim); !ok || delim != '[' {
		return nil, fmt.Errorf("expected opening array bracket '[', got %v", t)
	}

	songs := make([]song.Song, 0, 20000)

	for decoder.More() {
		var song song.Song
		if err := decoder.Decode(&song); err != nil {
			return nil, fmt.Errorf("error decoding object: %w", err)
		}
		songs = append(songs, song)
	}

	if _, err := decoder.Token(); err != nil {
		return nil, fmt.Errorf("failed to read end token: %w", err)
	}

	return songs, nil
}

type SongLeaderboardEntry struct {
	TrackName  string
	ArtistName string
	MsPlayed   time.Duration
}

func generateSongLeaderboard(songs []song.Song) []SongLeaderboardEntry {
	type trackKey struct {
		Track  string
		Artist string
	}

	entries := make(map[trackKey]int64, len(songs)/10)

	for _, song := range songs {
		if song.MasterMetadataTrackName == "" || song.MasterMetadataAlbumArtistName == "" {
			continue
		}

		key := trackKey{
			Track:  song.MasterMetadataTrackName,
			Artist: song.MasterMetadataAlbumArtistName,
		}

		entries[key] += song.MsPlayed
	}

	leaderboard := make([]SongLeaderboardEntry, 0, len(entries))
	for key, ms := range entries {
		leaderboard = append(leaderboard, SongLeaderboardEntry{
			TrackName:  key.Track,
			ArtistName: key.Artist,
			MsPlayed:   time.Duration(ms) * time.Millisecond,
		})
	}

	sort.Slice(leaderboard, func(i, j int) bool {
		return leaderboard[i].MsPlayed > leaderboard[j].MsPlayed
	})

	return leaderboard
}
