package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	fossildb "github.com/atbu/fossil/db"
	"github.com/atbu/fossil/records"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Incorrect usage.") // TODO make this more user friendly
		os.Exit(1)
	}

	db, err := fossildb.InitDB("spotify_data.db")
	if err != nil {
		log.Fatalf("Database initialisation error: %v", err)
	}
	defer db.Close()

	switch os.Args[1] {
	case "import":
		importCmd := flag.NewFlagSet("import", flag.ExitOnError)
		importCmd.Parse(os.Args[2:])

		if importCmd.NArg() < 1 {
			fmt.Println("Error: missing JSON file path")
			fmt.Println("Usage: fossil import <path-to-json>")
			os.Exit(1)
		}

		filePath := importCmd.Arg(0)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			fmt.Printf("Error: File %s does not exist.\n", filePath)
			os.Exit(1)
		}

		audioRecords, err := parseAudioStreamingHistoryFile(filePath)
		if err != nil {
			panic(err)
		}

		totalInserted, skipped, err := fossildb.BulkInsertAudioRecords(db, audioRecords)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Inserted %d rows and skipped %d rows.\n", totalInserted, skipped)
	}

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
