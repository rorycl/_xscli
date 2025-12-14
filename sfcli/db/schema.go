package db

// schema defines the SQL statements to create the application's database schema.
// It is designed to be idempotent using `CREATE TABLE IF NOT EXISTS`.
const schema = `
-- This schema is designed to be appended to by the xerocli tool.
-- The CREATE TABLE statements are idempotent.
CREATE TABLE IF NOT EXISTS salesforce_opportunities (
    id                      VARCHAR PRIMARY KEY,
    name                    VARCHAR,
    amount                  DECIMAL(18, 2),
    close_date              DATE,
    stage_name              VARCHAR,
    record_type_name        VARCHAR,
    payout_reference_dfk    VARCHAR,
    last_modified_date      TIMESTAMP
);
`

// oppsUpsertSQL is the SQL statement for inserting or updating a Salesforce Opportunity.
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
