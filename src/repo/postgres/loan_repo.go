package postgres

import (
	"context"
	"errors"
	"los/src/constants"
	"los/src/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type LoanService struct{ db *pgxpool.Pool }

func NewLoanService(db *pgxpool.Pool) *LoanService { return &LoanService{db: db} }

func (r *LoanService) CreateLoan(ctx context.Context, l *models.Loan) (string, error) {
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

func (r *LoanService) ListByStatus(ctx context.Context, status constants.Status, limit, offset int) ([]models.Loan, error) {
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

func (r *LoanService) GetLoanByID(ctx context.Context, loanID string) (*models.Loan, error) {
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

func (r *LoanService) TopCustomers(ctx context.Context, limit int) ([]struct {
	Customer string
	Total    int64
}, error) {
	rows, err := r.db.Query(ctx, `
		SELECT customer_name, COUNT(*) AS total
		FROM loans
		WHERE status IN ('APPROVED_BY_SYSTEM','APPROVED_BY_AGENT')
		GROUP BY customer_name
		ORDER BY total DESC
		LIMIT $1
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []struct {
		Customer string
		Total    int64
	}
	for rows.Next() {
		var c string
		var t int64
		if err := rows.Scan(&c, &t); err != nil {
			return nil, err
		}
		out = append(out, struct {
			Customer string
			Total    int64
		}{Customer: c, Total: t})
	}
	return out, rows.Err()
}

func (r *LoanService) ClaimNextApplied(ctx context.Context) (*models.Loan, error) {
	row := r.db.QueryRow(ctx, `
		UPDATE loans SET version = version + 1
		WHERE loan_id = (
			SELECT loan_id FROM loans WHERE status = 'APPLIED' ORDER BY created_at ASC LIMIT 1
		)
		RETURNING loan_id, customer_name, phone, amount, type, status, version, created_at
	`)
	var l models.Loan
	if err := row.Scan(&l.LoanID, &l.CustomerName, &l.PhoneNumber, &l.Amount, &l.Type, &l.Status, &l.Version, &l.CreatedAt); err != nil {
		return nil, err
	}
	return &l, nil
}

func (r *LoanService) UpdateStatusWithVersion(ctx context.Context, loanID string, fromVersion int64, newStatus constants.Status) error {
	ct, err := r.db.Exec(ctx, `
		UPDATE loans SET status = $1, version = version + 1
		WHERE loan_id = $2 AND version = $3
	`, string(newStatus), loanID, fromVersion)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return errors.New("version conflict")
	}
	return nil
}

func (r *LoanService) StatusCounts(ctx context.Context) (map[constants.Status]int64, error) {
	rows, err := r.db.Query(ctx, `SELECT status, COUNT(*) FROM loans GROUP BY status`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := make(map[constants.Status]int64)
	for rows.Next() {
		var s string
		var c int64
		if err := rows.Scan(&s, &c); err != nil {
			return nil, err
		}
		res[constants.Status(s)] = c
	}
	return res, rows.Err()
}
