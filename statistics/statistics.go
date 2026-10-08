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
