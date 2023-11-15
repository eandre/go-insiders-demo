CREATE TABLE people (
    name TEXT NOT NULL PRIMARY KEY,
    first_met TIMESTAMPTZ NOT NULL,
    meeting_count INT
);