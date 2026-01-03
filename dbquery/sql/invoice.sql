/*
 Reconciler app SQL
 invoice.sql
 detail of an invoice with line items and donation total
 started 01 January 2026

 Note @param comments declare a template value for middleware replacement.
 Note do _not_ use colons in sql or comments as it breaks the sqlx parser.
*/

SELECT
    *
    ,CASE WHEN donation_total = crms_total THEN
        1
     ELSE
        0
     END AS reconciled
FROM (
    WITH variables AS (
        SELECT
             'INV-2025-101' AS InvoiceNumber /* @param */
            ,'^(53|55|57).*' AS AccountCodes /* @param */
    )
    
    ,reconciled_donations_summed AS (
        SELECT
            payout_reference_dfk
            ,sum(amount) AS donation_sum
        FROM
            salesforce_opportunities
            ,variables
        WHERE
            payout_reference_dfk = variables.InvoiceNumber
        GROUP BY
            payout_reference_dfk
    )
    SELECT
        i.id
        ,i.invoice_number
        ,date(substring(i.date, 1, 10)) AS date
        ,i.type
        ,i.status
        ,i.reference
        ,i.contact_name
        ,i.total
        ,sum(li.line_amount) 
            FILTER (WHERE li.account_code REGEXP variables.AccountCodes)
            OVER (PARTITION BY i.id) AS donation_total
        ,rds.donation_sum AS crms_total
        ,a.name AS account_name
        ,li.description AS li_description
        ,li.tax_amount AS li_tax_amount
        ,li.line_amount AS li_line_amount
        ,CASE WHEN
            li.account_code REGEXP variables.AccountCodes
        THEN
            li.line_amount
         ELSE
            0
         END AS li_donation_amount
    FROM
        invoices i
        ,variables
        JOIN invoice_line_items li ON (li.invoice_id = i.id)
        LEFT OUTER JOIN accounts a ON (li.account_code = a.code)
        LEFT OUTER JOIN reconciled_donations_summed rds ON (rds.payout_reference_dfk = i.invoice_number)
    WHERE
        variables.InvoiceNumber = i.invoice_number
        /*
        AND
        i.status NOT IN ('DRAFT', 'DELETED', 'VOIDED')
        */
) x
;
