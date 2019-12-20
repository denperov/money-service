package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"syscall"

	"github.com/denperov/money-service/internal/accounts/endpoints"
	"github.com/denperov/money-service/internal/accounts/handlers"
	"github.com/denperov/money-service/internal/accounts/repository/postgres_repository"
	"github.com/denperov/money-service/internal/accounts/service"
	"github.com/denperov/money-service/internal/pkg/configs"
	"github.com/denperov/money-service/internal/pkg/http_server"
	"github.com/denperov/money-service/internal/pkg/signals_waiter"
)

func main() {
	log.Println("accounts service: start")
	defer log.Println("accounts service: stop")

	var cfg Config

	// STOP SIGNAL

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sig := signals_waiter.Wait(ctx, []os.Signal{syscall.SIGTERM, syscall.SIGINT})
		log.Printf("received signal %s", sig)
		cancel()
	}()

	// REPOSITORY

	configs.MustReadConfig(&cfg.Database)
	rep := postgres_repository.New(
		cfg.Database.Address,
		cfg.Database.Name,
		cfg.Database.User,
		cfg.Database.Password,
	)

	// SERVICE

	svc := service.New(rep)

	// ENDPOINTS AND HANDLERS

	mux := http.NewServeMux()
	mux.Handle("/get_accounts", handlers.MakeGetAccountsHandler(endpoints.MakeGetAccountsEndpoint(svc)))
	mux.Handle("/get_payments", handlers.MakeGetPaymentsHandler(endpoints.MakeGetPaymentsEndpoint(svc)))
	mux.Handle("/send_payment", handlers.MakeSendPaymentHandler(endpoints.MakeSendPaymentEndpoint(svc)))

	// SERVER

	configs.MustReadConfig(&cfg.API)
	server := http_server.New(
		cfg.API.ListenAddress,
		mux,
	)

	// RUN

	err := rep.Start(ctx)
	if err != nil {
		log.Fatalf("start repository: %v", err)
	}
	defer rep.Stop()

	server.Start(ctx)
	defer server.Stop()

	<-ctx.Done()
}
