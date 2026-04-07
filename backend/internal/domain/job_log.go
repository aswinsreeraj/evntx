package domain

import "time"

type JobStatus string

const (
	JobStatusSuccess JobStatus = "SUCCESS"
	JobStatusFailed  JobStatus = "FAILED"
)

type JobLog struct {
	ID           string    `json:"id"`
	JobName      string    `json:"job_name"`
	Status       JobStatus `json:"status"`
	Attempts     int       `json:"attempts"`
	ErrorMessage string    `json:"error_message"`
	StartedAt    time.Time `json:"started_at"`
	EndedAt      time.Time `json:"ended_at"`
}
