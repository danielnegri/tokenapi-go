package postgres

import (
	"context"
	"runtime"
	"strings"
	"time"

	"github.com/danielnegri/adheretech/errors"
	"github.com/danielnegri/adheretech/ledger"
	"github.com/danielnegri/adheretech/storage"
	"github.com/danielnegri/adheretech/valid"
	"github.com/go-pg/pg/v10"
)

const (
	DefaultURL = "postgres://localhost:5432/ledger"
)

type SecretToken struct {
	tableName struct{}     `pg:"secret_tokens,alias:tokens"`
	Data      ledger.Token `pg:"data,pk"`
}

type Postgres struct {
	db *pg.DB
}

var _ storage.Storage = (*Postgres)(nil)

// Connect parses a database URL into options that can be used to connect to PostgreSQL.
func Connect(opt *pg.Options) (*Postgres, error) {
	const op errors.Op = "storage/postgres.Connect"
	if opt == nil {
		return nil, errors.E(op, errors.Internal, "invalid database config")
	}

	if opt.MaxConnAge == 0 {
		opt.MaxConnAge = 10 * time.Minute
	}

	if opt.PoolSize == 0 {
		opt.PoolSize = runtime.NumCPU() * 2
	}

	db := pg.Connect(opt)
	db.AddQueryHook(DebugHook{})

	return &Postgres{db: db}, nil
}

func (p *Postgres) Insert(ctx context.Context, token ledger.Token) error {
	const op errors.Op = "storage/postgres.Insert"

	// This validation can be removed once it is enforce by database as well.
	// Although, it's much slower.
	if err := valid.Token(token); err != nil {
		return err
	}

	if err := p.db.Insert(&SecretToken{Data: token}); err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return errors.E(op, token, errors.Duplicate)
		}

		return errors.E(op, token, errors.Internal, err)
	}

	return nil

}

func (p *Postgres) Check(ctx context.Context) error {
	const op errors.Op = "storage/postgres.Check"

	if err := p.db.Ping(ctx); err != nil {
		return errors.E(op, errors.Internal, err)
	}

	return nil
}
