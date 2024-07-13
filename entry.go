package ledger

import "time"

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
