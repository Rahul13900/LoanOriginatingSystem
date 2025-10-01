package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AgentService struct{ db *pgxpool.Pool }

func NewAgentService(db *pgxpool.Pool) *AgentService { return &AgentService{db: db} }

func (r *AgentService) AssignLoanToAgent(ctx context.Context, agentID, loanID string) error {
	_, err := r.db.Exec(ctx, `INSERT INTO agent_loans(agent_id, loan_id) VALUES ($1,$2)`, agentID, loanID)
	return err
}

func (r *AgentService) IsAssigned(ctx context.Context, agentID, loanID string) (bool, error) {
	row := r.db.QueryRow(ctx, `SELECT 1 FROM agent_loans WHERE agent_id=$1 AND loan_id=$2 LIMIT 1`, agentID, loanID)
	var one int
	if err := row.Scan(&one); err != nil {
		return false, err
	}
	return true, nil
}

func (r *AgentService) PickAvailable(ctx context.Context) (string, error) {
	row := r.db.QueryRow(ctx, `SELECT agent_id FROM agents ORDER BY random() LIMIT 1`)
	var id string
	return id, row.Scan(&id)
}

func (r *AgentService) ManagerOf(ctx context.Context, agentID string) (string, error) {
	row := r.db.QueryRow(ctx, `SELECT manager_id FROM agents WHERE agent_id = $1`, agentID)
	var mid *string
	if err := row.Scan(&mid); err != nil {
		return "", err
	}
	if mid == nil {
		return "", nil
	}
	return *mid, nil
}
