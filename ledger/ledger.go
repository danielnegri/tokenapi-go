package ledger

import "context"

const (
	Description = "AdhereTech Ledger Service"
)

type Checker interface {
	Check(ctx context.Context) error
}

type Ledger interface {
	Insert(ctx context.Context, token Token) error
}

type Token string
