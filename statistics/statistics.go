package statistics

import (
	"sort"
	"time"

	"github.com/atbu/fossil/records"
)

type AudioLeaderboardEntry struct {
	TrackName  string
	ArtistName string
	TimePlayed time.Duration
}

func GenerateAudioLeaderboard(audioRecords []records.AudioRecord) []AudioLeaderboardEntry {
	type trackKey struct {
		TrackName  string
		ArtistName string
	}

	entries := make(map[trackKey]int64, len(audioRecords)/10)

	for _, audioRecord := range audioRecords {
		if audioRecord.MasterMetadataTrackName == "" || audioRecord.MasterMetadataAlbumArtistName == "" {
			continue
		}

		key := trackKey{
			TrackName:  audioRecord.MasterMetadataTrackName,
			ArtistName: audioRecord.MasterMetadataAlbumArtistName,
		}

		entries[key] += audioRecord.MsPlayed
	}

	leaderboard := make([]AudioLeaderboardEntry, 0, len(entries))
	for key, timePlayed := range entries {
		leaderboard = append(leaderboard, AudioLeaderboardEntry{
			TrackName:  key.TrackName,
			ArtistName: key.ArtistName,
			TimePlayed: time.Duration(timePlayed) * time.Millisecond,
		})
	}

	sort.Slice(leaderboard, func(i, j int) bool {
		return leaderboard[i].TimePlayed > leaderboard[j].TimePlayed
	})

	return leaderboard
}

type ArtistLeaderboardEntry struct {
	ArtistName     string
	NumberOfTracks int
	TimePlayed     time.Duration
}

func GenerateArtistLeaderboard(audioRecords []records.AudioRecord) []ArtistLeaderboardEntry {
	leaderboard := GenerateAudioLeaderboard(audioRecords)

	artists := make(map[string][]int64)

	for _, entry := range leaderboard {
		artists[entry.ArtistName] = append(artists[entry.ArtistName], entry.TimePlayed.Milliseconds())
	}

	artistLeaderboard := make([]ArtistLeaderboardEntry, 0, len(leaderboard)/2)
	for artist, trackDurations := range artists {
		numberOfTracks := len(trackDurations)

		var durationSum int64 = 0
		for _, duration := range trackDurations {
			durationSum += duration
		}

		entry := ArtistLeaderboardEntry{
			ArtistName:     artist,
			NumberOfTracks: numberOfTracks,
			TimePlayed:     time.Duration(durationSum) * time.Millisecond,
		}

		artistLeaderboard = append(artistLeaderboard, entry)
	}

	sort.Slice(artistLeaderboard, func(i, j int) bool {
		return artistLeaderboard[i].TimePlayed > artistLeaderboard[j].TimePlayed
	})

	return artistLeaderboard
}
