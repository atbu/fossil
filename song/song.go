package song

import "time"

type Song struct {
	Ts                            time.Time `json:"ts"`
	Platform                      string    `json:"platform"`
	MsPlayed                      int64     `json:"ms_played"`
	ConnCountry                   string    `json:"conn_country"`
	IPAddr                        string    `json:"ip_addr"`
	MasterMetadataTrackName       string    `json:"master_metadata_track_name"`
	MasterMetadataAlbumArtistName string    `json:"master_metadata_album_artist_name"`
	MasterMetadataAlbumAlbumName  string    `json:"master_metadata_album_album_name"`
	SpotifyTrackURI               string    `json:"spotify_track_uri"`
	EpisodeName                   string    `json:"episode_name"`
	EpisodeShowName               string    `json:"episode_show_name"`
	SpotifyEpisodeURI             string    `json:"spotify_episode_uri"`
	AudiobookTitle                string    `json:"audiobook_title"`
	AudiobookURI                  string    `json:"audiobook_uri"`
	AudiobookChapterURI           string    `json:"audiobook_chapter_uri"`
	AudiobookChapterTitle         string    `json:"audiobook_chapter_title"`
	ReasonStart                   string    `json:"reason_start"`
	ReasonEnd                     string    `json:"reason_end"`
	Shuffle                       bool      `json:"shuffle"`
	Skipped                       bool      `json:"skipped"`
	Offline                       bool      `json:"offline"`
	OfflineTimestamp              int64     `json:"offline_timestamp"`
	IncognitoMode                 bool      `json:"incognito_mode"`
}
