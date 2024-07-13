package ledger

import "time"

type TransactionEntry struct {
	CreatedAt     time.Time
	ID            int
	EntryID       int
	TransactionID int
}

type TransactionEntryFilter struct {
	TransactionID int
}
