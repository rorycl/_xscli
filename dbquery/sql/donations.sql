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
        date('2025-04-01') AS DateFrom /* @param */
        ,date('2026-03-31') AS DateTo  /* @param */
        -- All | Linked | NotLinked
        ,'All' AS LinkageStatus        /* @param */
        ,'' AS PayoutReference         /* @param */
        ,'Jane Smith' AS TextSearch    /* @param */
)

SELECT
    s.id  
    ,s.name 
    ,s.amount
    ,date(s.close_date) AS close_date /* column is defined as text */
    ,s.payout_reference_dfk
    ,s.created_date
    ,s.created_by_name
    ,s.last_modified_date
    ,s.last_modified_by_name 
    /* see www.sqlitetutorial.net/sqlite-json-functions/sqlite-json_extract-function/ */
    -- s.additional_fields_json  TEXT -- A JSON blob for all other fields
FROM salesforce_opportunities s
JOIN variables v ON s.close_date BETWEEN v.DateFrom AND v.DateTo
WHERE
    CASE 
        -- Searching by v.PayoutReference doesn't make sense if
        -- v.LinkageStatus = 'NotLinked', but there is no easy way of
        -- catching that error in the query.
        -- If translating to plpgsql add this check in the preamble.
        WHEN v.PayoutReference IS NOT NULL AND v.PayoutReference <> '' THEN
            v.PayoutReference = s.payout_reference_dfk
        ELSE TRUE
    END
    AND
    (
        (v.LinkageStatus = 'All')
        OR
        (v.LinkageStatus = 'Linked' AND s.payout_reference_dfk IS NOT NULL)
        OR
        (v.LinkageStatus = 'NotLinked' AND s.payout_reference_dfk IS NULL)
    )
    AND CASE
        WHEN v.TextSearch = '' OR v.TextSearch IS NULL THEN
            TRUE
        ELSE
            CONCAT(s.name, ' ', s.payout_reference_dfk) REGEXP v.TextSearch
        END
;
