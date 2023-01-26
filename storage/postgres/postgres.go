// Copyright 2020 The Ledger Authors
//
// Licensed under the AGPL, Version 3.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.gnu.org/licenses/agpl-3.0.en.html
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package postgres

import (
	"context"
	"runtime"
	"strings"
	"time"

	"github.com/danielnegri/tokenapi-go/errors"
	"github.com/danielnegri/tokenapi-go/ledger"
	"github.com/danielnegri/tokenapi-go/storage"
	"github.com/danielnegri/tokenapi-go/valid"
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

		if strings.Contains(err.Error(), "violates check constraint") {
			return errors.E(op, token, errors.Invalid)
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
