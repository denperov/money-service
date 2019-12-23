package main

import (
	"context"
	"flag"
	"log"

	"github.com/denperov/money-service/internal/accounts/endpoints"
	"github.com/denperov/money-service/internal/accounts/handlers"
	"github.com/denperov/money-service/internal/accounts/repository/postgres_repository"
	"github.com/denperov/money-service/internal/accounts/service"
	"github.com/denperov/money-service/internal/pkg/http_server"
	"github.com/denperov/money-service/internal/pkg/stop_signal"

	"github.com/gorilla/mux"
)

var (
	cfgListenAddress = flag.String("listen-address", "0.0.0.0:8080", "API listen address")

	cfgDatabaseAddress  = flag.String("database-address", "", "Database address")
	cfgDatabaseName     = flag.String("database-name", "", "Database name")
	cfgDatabaseUser     = flag.String("database-user", "", "Database user")
	cfgDatabasePassword = flag.String("database-password", "", "Database password")
)

func main() {
	flag.Parse()

	log.Println("accounts service: start")
	defer log.Println("accounts service: stop")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// STOP SIGNAL

	stopSignal := stop_signal.New(
		cancel,
	)

	// REPOSITORY

	rep := postgres_repository.New(
		*cfgDatabaseAddress,
		*cfgDatabaseName,
		*cfgDatabaseUser,
		*cfgDatabasePassword,
	)

	// SERVICE

	svc := service.New(rep)

	// ENDPOINTS AND HANDLERS

	r := mux.NewRouter()
	r.Methods("GET").Path("/accounts").Handler(
		handlers.MakeGetAccountsHandler(endpoints.MakeGetAccountsEndpoint(svc)),
	)
	r.Methods("GET").Path("/payments").Handler(
		handlers.MakeGetPaymentsHandler(endpoints.MakeGetPaymentsEndpoint(svc)),
	)
	r.Methods("POST").Path("/transfers").Handler(
		handlers.MakeCreateTransferHandler(endpoints.MakeCreateTransferEndpoint(svc)),
	)

	// SERVER

	server := http_server.New(
		*cfgListenAddress,
		r,
	)

	// RUN

	stopSignal.Start(ctx)
	defer stopSignal.Stop()

	err := rep.Start(ctx)
	if err != nil {
		log.Fatalf("start repository: %v", err)
	}
	defer rep.Stop()

	server.Start(ctx)
	defer server.Stop()

	<-ctx.Done()
}
