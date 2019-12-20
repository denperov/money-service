package service

import (
	"context"
)

type Repository interface {
	GetAccounts(context.Context) ([]Account, error)
	GetTransfers(context.Context) ([]Transfer, error)
	AddTransfer(context.Context, Transfer) error
}

type AccountsService interface {
	GetAccounts(context.Context) ([]Account, error)
	GetPayments(context.Context) ([]Payment, error)
	SendPayment(context.Context, Transfer) error
}

func New(
	repository Repository,
) AccountsService {
	return &accountsService{
		repository: repository,
	}
}

type accountsService struct {
	repository Repository
}

func (s *accountsService) GetAccounts(ctx context.Context) ([]Account, error) {
	return s.repository.GetAccounts(ctx)
}

func (s *accountsService) GetPayments(ctx context.Context) ([]Payment, error) {
	transfers, err := s.repository.GetTransfers(ctx)
	if err != nil {
		return nil, err
	}

	var payments []Payment
	for _, transfer := range transfers {
		outgoing := Payment{
			Direction: OutgoingDirection,
			Account:   transfer.FromAccount,
			ToAccount: transfer.ToAccount,
			Amount:    transfer.Amount,
		}
		incoming := Payment{
			Direction:   IncomingDirection,
			Account:     transfer.ToAccount,
			FromAccount: transfer.FromAccount,
			Amount:      transfer.Amount,
		}
		payments = append(payments, outgoing, incoming)
	}
	return payments, nil
}

func (s *accountsService) SendPayment(ctx context.Context, transfer Transfer) error {
	return s.repository.AddTransfer(ctx, transfer)
}
