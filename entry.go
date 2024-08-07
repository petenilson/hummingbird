package hummingbird

import (
	"context"
	"time"
)

type EntryType string

const (
	DEBIT  EntryType = "DEBIT"
	CREDIT EntryType = "CREDIT"
)

type Entry struct {
	ID        int
	AccountID int
	CreatedAt time.Time
	Amount    int
	Type      EntryType
}

type EntryFilter struct {
	AccountID *int
}

type EntryService interface {
	FindEntrys(context.Context, EntryFilter) ([]*Entry, int, error)
	CreateEntry(context.Context, *Entry) error
}
