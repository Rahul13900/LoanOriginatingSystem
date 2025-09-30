package models

import (
	"los/src/constants"
	"time"
)

type CreateLoanRequest struct {
	CustomerName string  `json:"customer_name"`
	PhoneNumber  string  `json:"phone_number"`
	Amount       float64 `json:"amount"`
	Type         string  `json:"type"`
}

type CreateLoanResponse struct {
	LoanID string `json:"loan_id"`
	Status string `json:"status"`
}

type AgentDecision struct {
	Decision string `json:"decision"`
}

type Loan struct {
	LoanID       string           `json:"loan_id"`
	CustomerName string           `json:"customer_name"`
	PhoneNumber  string           `json:"phone_number"`
	Amount       float64          `json:"amount"`
	Type         string           `json:"type"`
	Status       constants.Status `json:"status"`
	Version      int64            `json:"version"`
	CreatedAt    time.Time        `json:"created_at"`
}

type Agent struct {
	AgentID   string `json:"agent_id"`
	Name      string `json:"name"`
	ManagerID string `json:"manager_id"`
}
