-- name: PingDB :one
SELECT NOW();

-- name: CreateUser :execresult
INSERT INTO user (first_name, last_name, role, email, password, phone, created_at, updated_at) VALUES (?,?,?,?,?,?,?,?);

-- name: GetUserByEmail :one
SELECT * FROM user WHERE email = ?;

-- name: GetUserByID :one
SELECT * FROM user WHERE id = ?;

-- name: GetCustomers :many
SELECT 
    u.id, 
    u.first_name, 
    u.last_name, 
    u.status, 
    u.email, 
    u.phone,
    u.created_at, 
    COALESCE(inv.currency, '') AS currency,
    COUNT(inv.id) AS number_of_invoices, 
    CAST(COALESCE(SUM(inv.total_amount), 0) AS SIGNED) AS total_billed, 
    CAST(COALESCE(SUM(inv.amount_paid), 0) AS SIGNED) AS total_collected
FROM user u
LEFT JOIN (
    SELECT 
        i.id,
        i.currency, 
        i.user_id, 
        (i.amount + (i.amount * i.vat * 0.01)) AS total_amount, 
        COALESCE(SUM(t.amount), 0) AS amount_paid
    FROM invoice i
    LEFT JOIN transaction t 
        ON t.invoice_id = i.id
    GROUP BY i.id
) inv ON inv.user_id = u.id
WHERE u.role = 'customer'
  AND (sqlc.narg('customer_id') IS NULL OR u.id = sqlc.narg('customer_id'))
  AND (sqlc.narg('first_name') IS NULL OR u.first_name LIKE CONCAT('%', sqlc.narg('first_name'), '%'))
  AND (sqlc.narg('last_name') IS NULL OR u.last_name LIKE CONCAT('%', sqlc.narg('last_name'), '%'))
  AND (sqlc.narg('email') IS NULL OR u.email LIKE CONCAT('%', sqlc.narg('email'), '%'))
  AND (sqlc.narg('phone') IS NULL OR u.phone LIKE CONCAT('%', sqlc.narg('phone'), '%'))
GROUP BY 
    u.id, u.first_name, u.last_name, u.status, u.email, u.phone, u.created_at, inv.currency
ORDER BY
    CASE WHEN sqlc.narg('sort_by') = 'id' AND sqlc.narg('sort_order') = 'asc' THEN u.id END ASC,
    CASE WHEN sqlc.narg('sort_by') = 'id' AND sqlc.narg('sort_order') = 'desc' THEN u.id END DESC,
    CASE WHEN sqlc.narg('sort_by') = 'first_name' AND sqlc.narg('sort_order') = 'asc' THEN u.first_name END ASC,
    CASE WHEN sqlc.narg('sort_by') = 'first_name' AND sqlc.narg('sort_order') = 'desc' THEN u.first_name END DESC,
    CASE WHEN sqlc.narg('sort_by') = 'last_name' AND sqlc.narg('sort_order') = 'asc' THEN u.last_name END ASC,
    CASE WHEN sqlc.narg('sort_by') = 'last_name' AND sqlc.narg('sort_order') = 'desc' THEN u.last_name END DESC,
    CASE WHEN sqlc.narg('sort_by') = 'status' AND sqlc.narg('sort_order') = 'asc' THEN u.status END ASC,
    CASE WHEN sqlc.narg('sort_by') = 'status' AND sqlc.narg('sort_order') = 'desc' THEN u.status END DESC,
    CASE WHEN sqlc.narg('sort_by') = 'email' AND sqlc.narg('sort_order') = 'asc' THEN u.email END ASC,
    CASE WHEN sqlc.narg('sort_by') = 'email' AND sqlc.narg('sort_order') = 'desc' THEN u.email END DESC,
    CASE WHEN sqlc.narg('sort_by') = 'phone' AND sqlc.narg('sort_order') = 'asc' THEN u.phone END ASC,
    CASE WHEN sqlc.narg('sort_by') = 'phone' AND sqlc.narg('sort_order') = 'desc' THEN u.phone END DESC,
    CASE WHEN sqlc.narg('sort_by') = 'updated_at' AND sqlc.narg('sort_order') = 'asc' THEN u.updated_at END ASC,
    CASE WHEN sqlc.narg('sort_by') = 'updated_at' AND sqlc.narg('sort_order') = 'desc' THEN u.updated_at END DESC,
    CASE WHEN sqlc.narg('sort_by') = 'created_at' AND sqlc.narg('sort_order') = 'asc' THEN u.created_at END ASC,
    CASE WHEN sqlc.narg('sort_by') = 'created_at' AND sqlc.narg('sort_order') = 'desc' THEN u.created_at END DESC;




-- name: GetDefaultCurrency :one
SELECT setting_value FROM setting WHERE setting_key = 'currency';

-- name: GetDefaultVAT :one
SELECT setting_value FROM setting WHERE setting_key = 'vat';

-- name: GetInvoiceSettings :many
SELECT setting_key, setting_value FROM setting WHERE setting_key IN ('currency', 'vat', 'deadline_days');

-- name: CreateInvoice :execresult
INSERT INTO invoice (user_id, amount, vat, type, issued_at, from_date, until_date, created_at, updated_at, currency, deadline, created_by) VALUES (?,?,?,?,?,?,?,?,?,?,?,?);

-- name: GetInvoices :many
SELECT 
  inv.id, 
  inv.amount, 
  inv.vat,
  CAST(inv.amount + (inv.amount * inv.vat * 0.01) AS signed) AS total_amount, 
  inv.type, 
  inv.issued_at, 
  inv.from_date, 
  inv.until_date, 
  inv.deadline, 
  inv.currency, 
  inv.invoice_file, 
  u.first_name, 
  u.last_name, 
  u.email, 
  u.photo,
  u.phone
FROM 
  invoice inv
LEFT JOIN 
  user u 
  ON u.id = inv.user_id
WHERE 
  1;

-- name: CreateInvoiceActivity :execresult
INSERT INTO invoice_activity (user_id, invoice_id, action_type, resource_id, created_at, updated_at, attachment) VALUES (?,?,?,?,?,?,?);
