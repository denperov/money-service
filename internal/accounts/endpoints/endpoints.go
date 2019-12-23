package endpoints

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	"github.com/denperov/money-service/internal/accounts/models"
	"github.com/denperov/money-service/internal/accounts/service"
)

type GetAccountsRequest struct{}

type GetAccountsResponse struct {
	Accounts []models.Account `json:"accounts"`
}

func MakeGetAccountsEndpoint(s service.AccountsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		accounts, err := s.GetAccounts(ctx)
		if err != nil {
			return nil, err
		}
		return GetAccountsResponse{
			Accounts: accounts,
		}, nil
	}
}

type GetPaymentsRequest struct{}

type GetPaymentsResponse struct {
	Payments []models.Payment `json:"payments"`
}

func MakeGetPaymentsEndpoint(s service.AccountsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		payments, err := s.GetPayments(ctx)
		if err != nil {
			return nil, err
		}
		return GetPaymentsResponse{
			Payments: payments,
		}, nil
	}
}

type CreateTransferRequest struct {
	models.Transfer `json:"transfer"`
}

type CreateTransferResponse struct{}

func MakeCreateTransferEndpoint(s service.AccountsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateTransferRequest)
		err := s.CreateTransfer(ctx, req.Transfer)
		if err != nil {
			return nil, err
		}
		return CreateTransferResponse{}, nil
	}
}
