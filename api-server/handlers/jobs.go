package handlers

import (
	"api-server/db"
	"api-server/models"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func CreateJob(c *gin.Context) {
	user := c.MustGet("user").(*models.SupabaseJWTPayload)
	var req models.JobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()

	// Get user from database
	var dbUser models.User
	err := db.UsersDB.QueryRow(ctx, `
        SELECT uuid, last_request, is_premium, credits, identities 
        FROM users WHERE uuid = $1
    `, user.Sub).Scan(&dbUser.UUID, &dbUser.LastRequest, &dbUser.IsPremium, &dbUser.Credits, &dbUser.Identities)

	if err == pgx.ErrNoRows {
		identities := fmt.Sprintf("[[%q,%q,%q,%q,%q]]",
			user.UserMetadata.Iss,
			user.Sub,
			user.UserMetadata.Picture,
			user.UserMetadata.Name,
			user.UserMetadata.Nickname)

		_, err = db.UsersDB.Exec(ctx, `
            INSERT INTO users (uuid, last_request, is_premium, credits, identities)
            VALUES ($1, 0, false, 60, $2)
        `, user.Sub, identities)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to create user"})
			return
		}
		dbUser = models.User{UUID: user.Sub, LastRequest: 0, IsPremium: false, Credits: 60}
	}

	// Rate limiting logic
	currentTime := time.Now().Unix()
	rateLimit := int64(86400) // 24 hours
	if dbUser.IsPremium {
		rateLimit = 60 // 1 minute
	}

	if currentTime-dbUser.LastRequest < rateLimit {
		timeRemaining := rateLimit - (currentTime - dbUser.LastRequest)
		hours := int(timeRemaining / 3600)
		message := "Free users can make 1 request per day. Please wait %d hour(s)."
		if dbUser.IsPremium {
			message = "Premium users can make 1 request per minute."
		}
		c.JSON(429, gin.H{"error": fmt.Sprintf(message, hours)})
		return
	}

	// Create or update job
	if req.ContentID != "" {
		_, err = db.TranscriptionJobsDB.Exec(ctx, `
            INSERT INTO transcription_jobs (platform_id, channel_id, content_id, job_state, queued)
            VALUES ($1, $2, $3, 0, $4)
            ON CONFLICT (platform_id, content_id)
            DO UPDATE SET job_state = 0, queued = $4
        `, req.PlatformID, req.ChannelID, req.ContentID, currentTime)
	}
	// else {
	// 	_, err = db.IndexJobsDB.Exec(ctx, `
	//         INSERT INTO indexer_jobs (platform_id, channel_id, job_state, queued)
	//         VALUES ($1, $2, 0, $3)
	//         ON CONFLICT (platform_id, channel_id)
	//         DO UPDATE SET job_state = 0, queued = $3
	//     `, req.PlatformID, req.ChannelID, currentTime)
	// }

	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create job"})
		return
	}

	// Update last request time
	_, err = db.UsersDB.Exec(ctx, "UPDATE users SET last_request = $1 WHERE uuid = $2", currentTime, user.Sub)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(201, gin.H{"message": "Job created successfully"})
}

func GetJobStatus(c *gin.Context) {
	platformID := c.Query("platform_id")
	contentID := c.Query("content_id")
	channelID := strings.ToLower(c.Query("channel_id"))

	if platformID == "" || (contentID == "" && channelID == "") {
		c.JSON(400, gin.H{"error": "Missing required parameters"})
		return
	}

	ctx := context.Background()
	var job models.Job
	var err error
	if contentID != "" {
		err = db.TranscriptionJobsDB.QueryRow(ctx, `
            SELECT * FROM transcription_jobs WHERE platform_id = $1 AND content_id = $2
        `, platformID, contentID).Scan(
			&job.PlatformID, &job.ChannelID, &job.ContentID, &job.JobState, &job.Queued, &job.LastCompleted)
	}
	// else {
	// 	err = db.IndexJobsDB.QueryRow(ctx, `
	//         SELECT * FROM indexer_jobs WHERE platform_id = $1 AND channel_id = $2
	//     `, platformID, channelID).Scan(
	// 		&job.PlatformID, &job.ChannelID, &job.JobState, &job.Queued, &job.LastCompleted)
	// }

	if err == pgx.ErrNoRows {
		c.JSON(200, []interface{}{})
		return
	}

	response := map[string]interface{}{"state": job.JobState}
	if job.JobState == 0 {
		var count int
		if contentID != "" {
			err = db.TranscriptionJobsDB.QueryRow(ctx, `
                SELECT COUNT(*) FROM transcription_jobs 
                WHERE platform_id = $1 AND job_state = 0 AND queued < $2
            `, platformID, job.Queued).Scan(&count)
		}
		// else {
		// 	err = db.IndexJobsDB.QueryRow(ctx, `
		//         SELECT COUNT(*) FROM indexer_jobs
		//         WHERE platform_id = $1 AND job_state = 0 AND queued < $2
		//     `, platformID, job.Queued).Scan(&count)
		// }
		if err == nil {
			response["pos"] = count + 1
		}
	}

	c.JSON(200, []interface{}{response})
}
