/*
 Reconciler app SQL
 invoices.sql
 list of invoices with reconciliation status
 started 01 January 2026

 Note @param comments declare a template value for middleware replacement.
 Note do _not_ use colons in sql or comments as it breaks the sqlx parser.
*/

WITH variables AS (
    SELECT
        date('2025-04-01') AS DateFrom   /* @param */
        ,date('2026-03-31') AS DateTo    /* @param */
        ,'^(53|55|57).*' AS AccountCodes /* @param */
        -- All | Reconciled | NotReconciled
        ,'AllGenerous Individual' AS ReconciliationStatus   /* @param */ 
)

,invoice_donation_totals AS (
    SELECT
        li.invoice_id
        ,SUM(li.line_amount) AS total_donation_amount
    FROM invoice_line_items li
    JOIN invoices i ON (i.id = li.invoice_id)
    ,variables
    WHERE
        account_code REGEXP variables.AccountCodes
        AND
        i.status NOT IN ('DRAFT', 'DELETED', 'VOIDED')
        AND
        date >= variables.DateFrom AND date <= variables.DateTo
    GROUP BY
        li.invoice_id
), 

crms_donation_totals AS (
    SELECT
        payout_reference_dfk
        ,SUM(amount) AS total_crms_amount
    FROM salesforce_opportunities
    JOIN variables
    WHERE
        payout_reference_dfk IS NOT NULL
        AND
        close_date BETWEEN date(variables.DateFrom,'-60 day') AND date(variables.DateTo, '+60 day')
    GROUP BY
        payout_reference_dfk
)

,reconciliation_data AS (
    SELECT
        i.id
        ,i.invoice_number
        ,i.date
        ,i.contact_name
        ,i.total
        ,COALESCE(idt.total_donation_amount, 0) AS donation_total
        ,COALESCE(cdt.total_crms_amount, 0) AS crms_total
    FROM invoices i
    JOIN variables c ON i.date BETWEEN c.DateFrom AND c.DateTo
    LEFT JOIN invoice_donation_totals idt ON i.id = idt.invoice_id
    LEFT JOIN crms_donation_totals cdt ON i.invoice_number = cdt.payout_reference_dfk
    WHERE
        i.status NOT IN ('DRAFT', 'DELETED', 'VOIDED')
        AND idt.invoice_id IS NOT NULL 
)
SELECT 
    r.*,
    (donation_total = crms_total) AS is_reconciled
FROM reconciliation_data r
     ,variables v
WHERE
    (v.ReconciliationStatus = 'All')
    OR
    (v.ReconciliationStatus = 'Reconciled' AND r.donation_total = r.crms_total)
    OR
    (v.ReconciliationStatus = 'NotReconciled' AND r.donation_total <> r.crms_total);
;
