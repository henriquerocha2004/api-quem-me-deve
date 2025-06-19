-- name: SaveDebt :exec
-- description: Save a new debt
INSERT INTO public.debts (id, description, total_value, due_date, installments_quantity, debt_date, status, user_client_id, product_ids, service_ids, finished_at, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13);

-- name: UpdateDebt :exec
-- description: Update an existing debt
UPDATE public.debts
SET description = $2,
    total_value = $3,
    due_date = $4,
    installments_quantity = $5,
    debt_date = $6,
    status = $7,
    user_client_id = $8,
    product_ids = $9,
    service_ids = $10,
    finished_at = $11,
    updated_at = $12
WHERE id = $1;

-- name: CreateInstallment :exec
INSERT INTO public.installments (id, description, value, due_date, deb_date, status, payment_date, payment_method, number, debt_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);

-- name: CreateCancelInfo :exec
INSERT INTO public.cancel_info (id, reason, cancel_date, cancelled_by, created_at, updated_at, debt_id)
VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: CreateReversalInfo :exec
INSERT INTO public.reversal_info (id, reason, reversal_date, reversed_by, reversed_installment_qtd, cancelled_installment_qtd, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);