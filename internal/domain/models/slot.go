package models

import "time"

type StatusState int

const (
	StatusAvailable StatusState = iota
	StatusBooked
	StatusCompleted
	StatusCancelled
)

var status = map[StatusState]string{
	StatusAvailable: "available",
	StatusBooked:    "booked",
	StatusCompleted: "completed",
	StatusCancelled: "cancelled",
}

type Slot struct {
	ID            string
	InterviewerID string
	IntervieweeID string
	StartTime     *time.Time
	EndTime       *time.Time
	Status        string
}

func GetStatus(s StatusState) string {
	return status[s]
}
