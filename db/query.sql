-- name: PingDB :one
SELECT NOW();

-- name: CreateUser :execresult
INSERT INTO user (first_name, last_name, role, email, password, created_at, updated_at) VALUES (?,?,?,?,?,?,?);

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
GROUP BY 
    u.id, u.first_name, u.last_name, u.status, u.email, u.created_at;


-- name: GetDefaultCurrency :one
SELECT setting_value FROM setting WHERE setting_key = 'currency'