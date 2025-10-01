package postgres

import (
	"context"
	"los/src/models"
	"los/src/constants"

	"github.com/jackc/pgx/v5/pgxpool"
)

type LoanStore struct{ db *pgxpool.Pool }

func NewLoanStore(db *pgxpool.Pool) *LoanStore { return &LoanStore{db: db} }

func (r *LoanStore) CreateLoan(ctx context.Context, l *models.Loan) (string, error) {
	row := r.db.QueryRow(ctx, `
		INSERT INTO loans (customer_name, phone, amount, type, status)
		VALUES ($1,$2,$3,$4,$5)
		RETURNING loan_id
	`, l.CustomerName, l.PhoneNumber, l.Amount, l.Type, l.Status)
	var id string
	if err := row.Scan(&id); err != nil {
		return "", err
	}
	return id, nil
}

func (r *LoanStore) ListByStatus(ctx context.Context, status constants.Status, limit, offset int) ([]models.Loan, error) {
	rows, err := r.db.Query(ctx, `
		SELECT loan_id, customer_name, phone, amount, type, status, version, created_at
		FROM loans
		WHERE ($1 = '' OR status = $1)
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`, string(status), limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.Loan
	for rows.Next() {
		var l models.Loan
		if err := rows.Scan(&l.LoanID, &l.CustomerName, &l.PhoneNumber, &l.Amount, &l.Type, &l.Status, &l.Version, &l.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, l)
	}
	return out, rows.Err()
}

func (r *LoanStore) GetLoanByID(ctx context.Context, loanID string) (*models.Loan, error) {
	row := r.db.QueryRow(ctx, `
		SELECT loan_id, customer_name, phone, amount, type, status, version, created_at
		FROM loans WHERE loan_id = $1
	`, loanID)
	var l models.Loan
	if err := row.Scan(&l.LoanID, &l.CustomerName, &l.PhoneNumber, &l.Amount, &l.Type, &l.Status, &l.Version, &l.CreatedAt); err != nil {
		return nil, err
	}
	return &l, nil
}
