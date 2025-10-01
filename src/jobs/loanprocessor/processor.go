package loanprocessor

import (
	"context"
	"los/src/constants"
	"los/src/models"
	"math/rand"
	"time"

	"los/src/domain/notification"
)

type Repo interface {
	ClaimNextApplied(ctx context.Context) (*models.Loan, error)
	UpdateStatusWithVersion(ctx context.Context, loanID string, fromVersion int64, newStatus constants.Status) error
}

type AgentService interface {
	PickAvailable(ctx context.Context) (string, error)
	AssignLoanToAgent(ctx context.Context, agentID, loanID string) error
	ManagerOf(ctx context.Context, agentID string) (string, error)
}

type Processor struct {
	repo     Repo
	agents   AgentService
	notifier notification.NotificationService
	workers  int
	quit     chan struct{}
}

func New(repo Repo, agents AgentService, notifier notification.NotificationService, workers int) *Processor {
	return &Processor{repo: repo, agents: agents, notifier: notifier, workers: workers, quit: make(chan struct{})}
}

func (p *Processor) Start() {
	for i := 0; i < p.workers; i++ {
		go p.loop()
	}
}

func (p *Processor) Stop() { close(p.quit) }

func (p *Processor) loop() {
	ctx := context.Background()
	for {
		select {
		case <-p.quit:
			return
		default:
		}

		l, err := p.repo.ClaimNextApplied(ctx)
		if err != nil || l == nil {
			time.Sleep(500 * time.Millisecond)
			continue
		}
		delay := time.Duration(5+rand.Intn(21)) * time.Second
		time.Sleep(delay)

		// simple rule: even amounts approved by system, odd -> under review
		var newStatus constants.Status
		if int(l.Amount)%2 == 0 {
			newStatus = constants.StatusApprovedBySystem
		} else {
			newStatus = constants.StatusUnderReview
		}
		if err := p.repo.UpdateStatusWithVersion(ctx, l.LoanID, l.Version, newStatus); err != nil {
			continue
		}
		if newStatus == constants.StatusUnderReview {
			// assign agent and notify
			agentID, err := p.agents.PickAvailable(ctx)
			if err == nil && agentID != "" {
				_ = p.agents.AssignLoanToAgent(ctx, agentID, l.LoanID)
				_ = p.notifier.SendPushNotification(ctx, agentID, "New loan assigned: "+l.LoanID)
				if mid, err := p.agents.ManagerOf(ctx, agentID); err == nil && mid != "" {
					_ = p.notifier.SendPushNotification(ctx, mid, "Agent "+agentID+" received loan "+l.LoanID)
				}
			}
		}
	}
}
