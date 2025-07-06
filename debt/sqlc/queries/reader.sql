-- name: ClientUserDebts :many
-- description: Get all debts for a specific user client
SELECT 
    id, 
    description,
    total_value,
    due_date,
    installments_quantity,
    debt_date,
    status,
    user_client_id,
    product_ids,
    service_ids,
    finished_at
FROM public.debts
WHERE user_client_id = $1
ORDER BY created_at DESC;

-- name: DebtInstallments :many
-- description: Get all installments for a specific debt
SELECT 
    id,
    description,
    value,
    due_date,
    deb_date,
    status,
    payment_date,
    payment_method,
    number,
    debt_id
FROM public.installments
WHERE debt_id = $1
ORDER BY number;

-- name: DebtCancelInfo :one
-- description: Get cancel info for a specific debt
SELECT 
    id,
    reason,
    cancel_date,
    cancelled_by,
    debt_id
FROM public.cancel_info
WHERE debt_id = $1;

-- name: DebtReversalInfo :one
-- description: Get reversal info for a specific debt
SELECT 
    id,
    reason,
    reversal_date,
    reversed_by,
    reversed_installment_qtd,
    cancelled_installment_qtd,
    debt_id
FROM public.reversal_info
WHERE debt_id = $1;

-- name: DebtByFilters :many
-- description: Get debts by various filters
SELECT 
    id, 
    description,
    total_value,
    due_date,
    installments_quantity,
    debt_date,
    status,
    user_client_id,
    product_ids,
    service_ids,
    finished_at
FROM public.debts
WHERE 
    ($1::TEXT IS NULL OR user_client_id = $1) AND
    ($2::TEXT IS NULL OR status = $2) AND
    ($3::TIMESTAMP IS NULL OR due_date >= $3) AND
    ($4::TIMESTAMP IS NULL OR due_date <= $4)
ORDER BY created_at DESC LIMIT $5 OFFSET $6;