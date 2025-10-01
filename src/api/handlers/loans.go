package handlers

import (
	"los/src/constants"
	"los/src/domain/loan"
	"los/src/models"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type LoansHandler struct{ svc *loan.Service }

func NewLoansHandler(s *loan.Service) *LoansHandler { return &LoansHandler{svc: s} }

// RegisterLoanRoutes mounts loan-related endpoints (gin)
func RegisterLoanRoutes(r *gin.Engine, h *LoansHandler) {
	v1 := r.Group("/api/v1")
	v1.POST("/loans", h.createLoan)
	v1.GET("/loans", h.listLoans)
}

func (h *LoansHandler) createLoan(c *gin.Context) {
	var req models.CreateLoanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, "invalid json")
		return
	}
	if req.CustomerName == "" || req.PhoneNumber == "" || req.Amount <= 0 || req.Type == "" || math.IsNaN(req.Amount) || math.IsInf(req.Amount, 0) {
		c.String(http.StatusBadRequest, "invalid payload")
		return
	}
	id, err := h.svc.CreateLoan(c.Request.Context(), &models.Loan{
		CustomerName: req.CustomerName,
		PhoneNumber:  req.PhoneNumber,
		Amount:       req.Amount,
		Type:         req.Type,
	})
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, models.CreateLoanResponse{LoanID: id, Status: string(constants.StatusApplied)})
}

func (h *LoansHandler) listLoans(c *gin.Context) {
	status := constants.Status(c.Query("status"))
	page := 1
	size := 10
	if v := c.Query("page"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			page = n
		}
	}
	if v := c.Query("size"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 100 {
			size = n
		}
	}
	offset := (page - 1) * size
	items, err := h.svc.ListByStatus(c.Request.Context(), status, size, offset)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, items)
}
