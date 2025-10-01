package server

import (
	"context"
	"log"
	"los/src/api/handlers"
	"los/src/config"
	"los/src/domain/loan"
	"los/src/domain/notification"
	"los/src/jobs/loanprocessor"
	"los/src/repo/postgres"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func New(cfg config.Config) (*http.Server, func()) {
	r := gin.New()
	r.Use(gin.Recovery())
	r.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })

	// Initialize DB pool
	dbpool, err := pgxpool.New(context.Background(), cfg.DBURL)
	if err != nil {
		log.Fatalf("db pool init failed: %v", err)
	}

	loanRepo := postgres.NewLoanService(dbpool)
	agentRepo := postgres.NewAgentService(dbpool)
	notifier := notification.MockService{}
	svc := loan.NewService(loanRepo, agentRepo)
	loansHandler := handlers.NewLoansHandler(svc)

	handlers.RegisterLoanRoutes(r, loansHandler)

	// Start background processor
	proc := loanprocessor.New(loanRepo, agentRepo, notifier, cfg.WorkerCount)
	proc.Start()

	srv := &http.Server{Addr: ":" + cfg.Port, Handler: r}
	cleanup := func() {
		// stop workers and close db
		proc.Stop()
		dbpool.Close()
	}
	return srv, cleanup
}
