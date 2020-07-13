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
