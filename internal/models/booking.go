// internal/models/booking.go
package models

import (
	"time"

	"gorm.io/gorm"
)

type PriorityLevel int

const (
	UnbaptizedContact    PriorityLevel = 1
	PersecutedMember     PriorityLevel = 2
	UnbaptizedZoom       PriorityLevel = 3
	PersecutedZoomMember PriorityLevel = 4
	BaptizedZoom         PriorityLevel = 5
	GroupActivities      PriorityLevel = 6
	TeamActivities       PriorityLevel = 7
)

type Booking struct {
	gorm.Model // Adds fields `ID`, `CreatedAt`, `UpdatedAt`, `DeletedAt`
	SpaceID    int64
	StartTime  time.Time
	EndTime    time.Time
	User       string
	Notes      string
	Status     string
	Priority   PriorityLevel
}
