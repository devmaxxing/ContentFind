package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	UsersDB *pgxpool.Pool
	//IndexJobsDB         *pgxpool.Pool
	TranscriptionJobsDB *pgxpool.Pool
	YouTubeClipsDB      *pgxpool.Pool
	TwitchClipsDB       *pgxpool.Pool
)

func Init() error {
	connStr := os.Getenv("DATABASE_URL")
	ctx := context.Background()

	var err error

	// Initialize all DB connection pools
	UsersDB, err = pgxpool.New(ctx, connStr)
	if err != nil {
		return fmt.Errorf("failed to create users database pool: %v", err)
	}

	// IndexJobsDB, err = pgxpool.New(ctx, connStr)
	// if err != nil {
	// 	return fmt.Errorf("failed to create index jobs database pool: %v", err)
	// }

	TranscriptionJobsDB, err = pgxpool.New(ctx, connStr)
	if err != nil {
		return fmt.Errorf("failed to create transcription jobs database pool: %v", err)
	}

	YouTubeClipsDB, err = pgxpool.New(ctx, connStr)
	if err != nil {
		return fmt.Errorf("failed to create youtube clips database pool: %v", err)
	}

	TwitchClipsDB, err = pgxpool.New(ctx, connStr)
	if err != nil {
		return fmt.Errorf("failed to create twitch clips database pool: %v", err)
	}

	return nil
}

// Add cleanup function
func Cleanup() {
	if UsersDB != nil {
		UsersDB.Close()
	}
	// if IndexJobsDB != nil {
	// 	IndexJobsDB.Close()
	// }
	if TranscriptionJobsDB != nil {
		TranscriptionJobsDB.Close()
	}
	if YouTubeClipsDB != nil {
		YouTubeClipsDB.Close()
	}
	if TwitchClipsDB != nil {
		TwitchClipsDB.Close()
	}
}

func GetClipsDB(platform string) (*pgxpool.Pool, error) {
	switch platform {
	case "youtube":
		return YouTubeClipsDB, nil
	case "twitch":
		return TwitchClipsDB, nil
	default:
		return nil, fmt.Errorf("unsupported platform: %s", platform)
	}
}
