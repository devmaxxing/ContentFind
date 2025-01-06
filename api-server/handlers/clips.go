package handlers

import (
	"api-server/db"
	"api-server/models"
	"context"
	"strings"

	"github.com/gin-gonic/gin"
)

func CreateClip(c *gin.Context) {
	user := c.MustGet("user").(*models.SupabaseJWTPayload)
	var req models.ClipRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	req.ChannelID = strings.ToLower(req.ChannelID)

	clipsDB, err := db.GetClipsDB(req.Platform)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	_, err = clipsDB.Exec(context.Background(), `
        INSERT INTO clips (channel_id, content_id, start_time, duration, title, user_id)
        VALUES ($1, $2, $3, $4, $5, $6)
    `, req.ChannelID, req.ContentID, req.StartTime, req.Duration, req.Title, user.Sub)

	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create clip"})
		return
	}

	c.JSON(201, gin.H{"message": "Clip created successfully"})
}

func GetClips(c *gin.Context) {
	videoID := c.Param("videoId")
	platform := c.Query("platform")

	if platform == "" {
		c.JSON(400, gin.H{"error": "Platform query parameter is required"})
		return
	}

	clipsDB, err := db.GetClipsDB(platform)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	rows, err := clipsDB.Query(context.Background(), `
        SELECT channel_id, content_id, start_time, duration, title, user_id 
        FROM clips WHERE content_id = $1
    `, videoID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch clips"})
		return
	}
	defer rows.Close()

	var clips []map[string]interface{}
	for rows.Next() {
		var clip struct {
			ChannelID string
			ContentID string
			StartTime float64
			Duration  float64
			Title     string
			UserID    string
		}
		if err := rows.Scan(&clip.ChannelID, &clip.ContentID, &clip.StartTime,
			&clip.Duration, &clip.Title, &clip.UserID); err != nil {
			continue
		}
		clips = append(clips, map[string]interface{}{
			"channel_id": clip.ChannelID,
			"content_id": clip.ContentID,
			"start_time": clip.StartTime,
			"duration":   clip.Duration,
			"title":      clip.Title,
			"user_id":    clip.UserID,
		})
	}

	c.JSON(200, clips)
}
