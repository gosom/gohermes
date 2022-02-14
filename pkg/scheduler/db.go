package scheduler

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/rs/zerolog"
)

const dbschema = `
CREATE SCHEMA IF NOT EXISTS scheduler;
SET search_path TO scheduler;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS jobs(
	id UUID NOT NULL DEFAULT uuid_generate_v4(),
	data JSONB NOT NULL,
	scheduled_at TIMESTAMP WITH TIME ZONE NOT NULL,
	created_at TIMESTAMP WITH TIME ZONE NOT NULL,
	PRIMARY KEY(id)
);
CREATE INDEX IF NOT EXISTS jobs_scheduled_at_idx ON jobs(scheduled_at);

CREATE TABLE IF NOT EXISTS job_executions(
	id SERIAL NOT NULL,
	job_id UUID NOT NULL,
	created_at TIMESTAMP WITH TIME ZONE NOT NULL,
	success BOOLEAN NOT NULL,
	msg VARCHAR(100) NOT NULL,
	PRIMARY KEY(id),
	FOREIGN KEY(job_id)
		REFERENCES jobs(id) ON DELETE CASCADE
);
CREATE UNIQUE INDEX IF NOT EXISTS job_executions_success ON job_executions (job_id, success)
WHERE success;
`

type schedulerRepository struct {
	log zerolog.Logger
	db  *sql.DB
}

func newSchedulerRepository(dsn string, logger zerolog.Logger) (*schedulerRepository, error) {
	ans := schedulerRepository{
		log: logger,
	}
	var err error
	ans.db, err = sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	return &ans, ans.db.Ping()
}

func (o *schedulerRepository) migrate(ctx context.Context) error {
	tx, err := o.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.ExecContext(ctx, dbschema)
	if err != nil {
		return err
	}
	return tx.Commit()
}
