package models

import "time"

type Incident struct {
	FeedID int
	IncidentKey string
	Title string
	Link string
	Status string
	Severity string
	StartedAt time.Time
	UpdatedAt time.Time
	ResolvedAt *time.Time
	LastEventHash string
	LastEventGUID string
}

type Feed struct {
	FeedID   int
	IncidentKey string
	EventGUID string
	ContentHash string
	Title string
	Link string
	DerivedStatus string
	DerivedSeverity string
	PublishedAt *time.Time
	IngestedAt time.Time
	RawContent string
}