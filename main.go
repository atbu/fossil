package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/atbu/fossil/records"
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
}

func parseAudioStreamingHistoryFile(filePath string) ([]records.AudioRecord, error) {
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

	songs := make([]records.AudioRecord, 0, 20000)

	for decoder.More() {
		var song records.AudioRecord
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
