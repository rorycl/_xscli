package db

// schema defines the SQL statements to create the application's database schema for SQLite.
// It is designed to be idempotent using `CREATE TABLE IF NOT EXISTS`.
const schema = `
-- This schema is designed to be appended to by the xerocli tool.
-- The CREATE TABLE statements are idempotent.
CREATE TABLE IF NOT EXISTS salesforce_opportunities (
    id                      TEXT PRIMARY KEY,
    name                    TEXT,
    amount                  REAL,
    close_date              TEXT, -- Using TEXT for dates is simplest for this app
    stage_name              TEXT,
    record_type_name        TEXT,
    payout_reference_dfk    TEXT,
    last_modified_date      DATETIME
);
`

// oppsUpsertSQL is the SQL statement for inserting or updating a Salesforce Opportunity in SQLite.
const oppsUpsertSQL = `
INSERT INTO salesforce_opportunities (id, name, amount, close_date, stage_name, record_type_name, payout_reference_dfk, last_modified_date)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT (id) DO UPDATE SET
    name = excluded.name,
    amount = excluded.amount,
    close_date = excluded.close_date,
    stage_name = excluded.stage_name,
    record_type_name = excluded.record_type_name,
    payout_reference_dfk = excluded.payout_reference_dfk,
    last_modified_date = excluded.last_modified_date;
`
