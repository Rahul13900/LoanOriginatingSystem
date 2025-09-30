package constants

type Status string

const (
	StatusApplied          Status = "APPLIED"
	StatusApprovedBySystem Status = "APPROVED_BY_SYSTEM"
	StatusRejectedBySystem Status = "REJECTED_BY_SYSTEM"
	StatusApprovedByAgent  Status = "APPROVED_BY_AGENT"
	StatusRejectedByAgent  Status = "REJECTED_BY_AGENT"
	StatusUnderReview      Status = "UNDER_REVIEW"
)
