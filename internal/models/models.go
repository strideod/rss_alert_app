package models

import "time"

type Feed struct {
	ID int
	Name string
	URL string
}

type Incident struct {
	ID int
	FeedID int
	EventKey string
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

type Event struct {
	ID int
	FeedID   int
	EventKey string
	EventGUID string
	ContentHash string
	FeedName string
	Title string
	Link string
	DerivedStatus string
	DerivedSeverity string
	PublishedAt *time.Time
	IngestedAt time.Time
	RawContent string
}