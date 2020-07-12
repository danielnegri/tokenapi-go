package storage

import (
	"context"

	"github.com/danielnegri/adheretech/ledger"
)

type Storage interface {
	Insert(ctx context.Context, token ledger.Token) error
	Check(ctx context.Context) error
}
