/*
 Reconciler app SQL
 bank_transaction_lis_insert.sql
 Insert a bank transaction line item

 Note @param comments declare a template value for middleware replacement.
 Note do _not_ use colons in sql or comments as it breaks the sqlx parser.
*/

WITH variables AS (
    SELECT
        'sf-opp-003' AS ID
        ,'Anonymous Donor' AS Name
        ,'21.20' AS Amount
        ,datetime('2025-04-14') AS CloseDate
        ,'JG-PAYOUT-2025-04-15' AS PayoutReference
        ,datetime('2025-04-01') AS CreatedDate
        ,'User1' AS CreatedBy
        ,datetime('2025-04-01') AS LastModifiedDate
        ,'User1' AS LastModifiedBy
        ,'' AS AdditionalFieldsJSON
)
INSERT INTO donations (
    id
    ,name
    ,amount
    ,close_date
    ,payout_reference_dfk
    ,created_date
    ,created_by_name
    ,last_modified_date
    ,last_modified_by_name
    ,additional_fields_json
)
SELECT
    v.ID
    ,v.Name
    ,v.Amount
    ,v.CloseDate
    ,v.PayoutReference
    ,v.CreatedDate
    ,v.CreatedBy
    ,v.LastModifiedDate
    ,v.LastModifiedBy
    ,v.AdditionalFieldsJSON
FROM
    variables v
-- https://sqlite.org/lang_upsert.html PARSING AMBIGUITY
WHERE
    true
ON CONFLICT (id) DO UPDATE SET
    name                    = excluded.name
    ,amount                 = excluded.amount
    ,close_date             = excluded.close_date
    ,payout_reference_dfk   = excluded.payout_reference_dfk
    ,created_date           = excluded.created_date
    ,created_by_name        = excluded.created_by_name
    ,last_modified_date     = excluded.last_modified_date
    ,last_modified_by_name  = excluded.last_modified_by_name
    ,additional_fields_json = excluded.additional_fields_json
;
