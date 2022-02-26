package scheduler

const dbschema = `
CREATE SCHEMA IF NOT EXISTS scheduler;
SET search_path TO scheduler;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS jobs(
	id UUID NOT NULL DEFAULT uuid_generate_v4(),
	endpoint VARCHAR(254) NOT NULL,
	data JSONB NOT NULL,
	scheduled_at TIMESTAMP WITH TIME ZONE NOT NULL,
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT (now() at time zone 'utc'),
	selected BOOLEAN NOT NULL DEFAULT 'f',
	PRIMARY KEY(id)
);
CREATE INDEX IF NOT EXISTS jobs_scheduled_at_idx ON jobs(scheduled_at);

CREATE TABLE IF NOT EXISTS job_executions(
	id SERIAL NOT NULL,
	job_id UUID NOT NULL,
	created_at TIMESTAMP WITH TIME ZONE NOT NULL,
	success BOOLEAN NOT NULL,
	data JSONB NOT NULL DEFAULT '{}',
	PRIMARY KEY(id),
	FOREIGN KEY(job_id)
		REFERENCES jobs(id) ON DELETE CASCADE
);
CREATE UNIQUE INDEX IF NOT EXISTS job_executions_success ON job_executions (job_id, success)
WHERE success;

CREATE OR REPLACE FUNCTION notify_event() RETURNS TRIGGER AS $$

    DECLARE 
        notification jsonb;
    BEGIN
        notification = json_build_object(
                          'table',TG_TABLE_NAME,
                          'action', TG_OP);
                        
        PERFORM pg_notify('events', notification::text);
        RETURN NULL; 
    END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS jobs_notify_event ON jobs;

CREATE TRIGGER jobs_notify_event
AFTER INSERT OR UPDATE OR DELETE ON jobs
    FOR EACH ROW EXECUTE PROCEDURE notify_event();
`

const insertQ = `INSERT INTO "scheduler".jobs(endpoint, data, scheduled_at) VALUES($1, $2, $3) RETURNING id`
const findQ = `SELECT J.* FROM "scheduler".jobs J 
WHERE 
J.scheduled_at <= (NOW() at time zone 'utc') 
AND J.selected = 'f' 
ORDER BY scheduled_at
LIMIT $1
FOR UPDATE`
const updateJobsQ = `UPDATE "scheduler".jobs SET selected = 't' WHERE id IN (%s)`
const nextJobQ = `SELECT * from "scheduler".jobs where selected = 'f' and scheduled_at > $1 LIMIT 1`
const insertExecResp = `INSERT INTO "scheduler".job_executions(job_id, created_at, success, data) VALUES($1, $2, $3, $4)`
