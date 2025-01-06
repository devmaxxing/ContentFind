package models

type JobRequest struct {
	PlatformID int    `json:"platform_id" binding:"required"`
	ChannelID  string `json:"channel_id,omitempty"`
	ContentID  string `json:"content_id,omitempty"`
}

type ClipRequest struct {
	Platform  string  `json:"platform" binding:"required"`
	ChannelID string  `json:"channel_id"`
	ContentID string  `json:"content_id" binding:"required"`
	StartTime float64 `json:"start_time" binding:"required"`
	Duration  float64 `json:"duration" binding:"required"`
	Title     string  `json:"title"`
}

type User struct {
	UUID        string `json:"uuid"`
	LastRequest int64  `json:"last_request"`
	IsPremium   bool   `json:"is_premium"`
	Credits     int    `json:"credits"`
	Identities  string `json:"identities"`
}

type Job struct {
	PlatformID    string `json:"platform_id"`
	ChannelID     string `json:"channel_id"`
	ContentID     string `json:"content_id,omitempty"`
	JobState      int    `json:"job_state"`
	Queued        int64  `json:"queued"`
	LastCompleted int64  `json:"last_completed,omitempty"`
}

type JobStatus struct {
	State int `json:"state"`
	Pos   int `json:"pos,omitempty"`
}
