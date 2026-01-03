/*
 Reconciler app SQL
 donations.sql
 list of donations with linkage status
 linkage status in this case only relates to whether the distributed foreign key
 (DFK) field payout_reference_dfk has a value or not.
 started 02 January 2026

 Note @param comments declare a template value for middleware replacement.
 Note do _not_ use colons in sql or comments as it breaks the sqlx parser.
*/

WITH variables AS (
    SELECT
        date('2025-04-01') AS DateFrom             /* @param */
        ,date('2026-03-31') AS DateTo              /* @param */
        -- All | Linked | NotLinked
        ,'Linked' AS LinkageStatus                 /* @param */
        ,'JG-PAYOUT-2025-04-15' AS PayoutReference /* @param */
)

SELECT
    s.id  
    ,s.name 
    ,s.amount
    ,date(substring(s.close_date, 1, 10)) as close_date
    ,s.payout_reference_dfk
    ,date(substring(s.created_date, 1, 10)) as created_date
    ,s.created_by_name
    ,date(substring(s.last_modified_date, 1, 10)) as last_modified_date
    ,s.last_modified_by_name 
    /* see www.sqlitetutorial.net/sqlite-json-functions/sqlite-json_extract-function/ */
    -- s.additional_fields_json  TEXT -- A JSON blob for all other fields
FROM salesforce_opportunities s
JOIN variables v ON s.close_date BETWEEN v.DateFrom AND v.DateTo
WHERE
    (v.LinkageStatus = 'All')
    OR
    (v.LinkageStatus = 'Linked'
        AND s.payout_reference_dfk IS NOT NULL
        AND (
            (v.PayoutReference is null)
            OR
            (v.PayoutReference = s.payout_reference_dfk)
        )
    )
    OR
    (v.LinkageStatus = 'NotLinked' AND s.payout_reference_dfk IS NULL)
;
