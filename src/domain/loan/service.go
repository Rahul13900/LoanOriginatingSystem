package loan

import (
	"context"
	"errors"
	"los/src/constants"
	"los/src/models"
)

type LoanService interface {
	CreateLoan(ctx context.Context, req *models.Loan) (string, error)
	ListByStatus(ctx context.Context, status constants.Status, limit, offset int) ([]models.Loan, error)
	GetLoanByID(ctx context.Context, loanID string) (*models.Loan, error)
	StatusCounts(ctx context.Context) (map[constants.Status]int64, error)
	TopCustomers(ctx context.Context, limit int) ([]struct {
		Customer string
		Total    int64
	}, error)
}

type AgentService interface {
	AssignLoanToAgent(ctx context.Context, loanID string, agentID string) error
	IsAssigned(ctx context.Context, agentID, loanID string) (bool, error)
	PickAvailable(ctx context.Context) (string, error)
	ManagerOf(ctx context.Context, agentID string) (string, error)
}

type Service struct {
	LoanService  LoanService
	AgentService AgentService
}

func NewService(loanService LoanService, agentService AgentService) *Service {
	return &Service{
		LoanService:  loanService,
		AgentService: agentService,
	}
}

func (s *Service) CreateLoan(ctx context.Context, req *models.Loan) (string, error) {
	req.Status = constants.StatusApplied
	return s.LoanService.CreateLoan(ctx, req)
}

func (s *Service) ListByStatus(ctx context.Context, status constants.Status, limit, offset int) ([]models.Loan, error) {
	return s.LoanService.ListByStatus(ctx, status, limit, offset)
}

func (s *Service) GetLoanByID(ctx context.Context, loanID string) (*models.Loan, error) {
	return s.LoanService.GetLoanByID(ctx, loanID)
}

func (s *Service) StatusCounts(ctx context.Context) (map[constants.Status]int64, error) {
	return s.LoanService.StatusCounts(ctx)
}

func (s *Service) TopCustomers(ctx context.Context, n int) ([]struct {
	Customer string
	Total    int64
}, error) {
	return s.LoanService.TopCustomers(ctx, n)
}

func (s *Service) AgentDecision(ctx context.Context, agentID, loanID string, approve bool) error {
	l, err := s.LoanService.GetLoanByID(ctx, loanID)
	if err != nil {
		return err
	}
	if l.Status != constants.StatusUnderReview {
		return errors.New("loan not under review")
	}
	assigned, err := s.AgentService.IsAssigned(ctx, agentID, loanID)
	if err != nil {
		return err
	}
	if !assigned {
		return errors.New("loan not assigned to agent")
	}
	newStatus := constants.StatusRejectedByAgent
	if approve {
		newStatus = constants.StatusApprovedByAgent
	}
	// optimistic update by version would be better here; simplified for brevity
	if up, ok := s.LoanService.(interface {
		UpdateStatusWithVersion(context.Context, string, int64, constants.Status) error
	}); ok {
		return up.UpdateStatusWithVersion(ctx, loanID, l.Version, newStatus)
	}
	return nil
}
