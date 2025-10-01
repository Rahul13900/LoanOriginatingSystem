package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AgentStore struct{ db *pgxpool.Pool }

func NewAgentStore(db *pgxpool.Pool) *AgentStore { return &AgentStore{db: db} }

func (r *AgentStore) AssignLoanToAgent(ctx context.Context, agentID, loanID string) error {
	_, err := r.db.Exec(ctx, `INSERT INTO agent_loans(agent_id, loan_id) VALUES ($1,$2)`, agentID, loanID)
	return err
}

func (r *AgentStore) IsAssigned(ctx context.Context, agentID, loanID string) (bool, error) {
	row := r.db.QueryRow(ctx, `SELECT 1 FROM agent_loans WHERE agent_id=$1 AND loan_id=$2 LIMIT 1`, agentID, loanID)
	var one int
	if err := row.Scan(&one); err != nil {
		return false, err
	}
	return true, nil
}
