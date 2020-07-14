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
	"time"

	"github.com/danielnegri/adheretech/log"
	"github.com/go-pg/pg/v10"
)

type DebugHook struct{}

var _ pg.QueryHook = (*DebugHook)(nil)

func (DebugHook) BeforeQuery(ctx context.Context, event *pg.QueryEvent) (context.Context, error) {
	return ctx, nil
}

func (DebugHook) AfterQuery(ctx context.Context, event *pg.QueryEvent) error {
	logger := log.WithField("latency", time.Since(event.StartTime))

	query, err := event.FormattedQuery()
	if err != nil {
		logger.Errorf("failed to format query: %v", err)
		return err
	}

	if event.Err != nil {
		logger.Errorf("error %s executing query: %s", event.Err, query)
	} else {
		logger.Debugf("Query processed: %s", query)
	}

	return nil
}
