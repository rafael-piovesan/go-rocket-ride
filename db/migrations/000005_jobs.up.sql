--
-- A relation that holds our transactionally-staged jobs
-- (see "Background jobs and job staging" above).
--
CREATE TABLE staged_jobs (
    id                 BIGSERIAL       PRIMARY KEY,
    job_name           TEXT            NOT NULL,
    job_args           JSONB           NOT NULL
);