package ledger

import "context"

type Ledger struct {
	ctx                context.Context
	AccountService     AccountService
	TransactionService TransactionService
}
