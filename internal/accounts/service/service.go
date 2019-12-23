package service

import (
	"context"

	"github.com/denperov/money-service/internal/accounts/models"
)

type Repository interface {
	GetAccounts(context.Context) ([]models.Account, error)
	GetTransfers(context.Context) ([]models.Transfer, error)
	AddTransfer(context.Context, models.Transfer) error
}

type AccountsService interface {
	GetAccounts(context.Context) ([]models.Account, error)
	GetPayments(context.Context) ([]models.Payment, error)
	CreateTransfer(context.Context, models.Transfer) error
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

func (s *accountsService) GetAccounts(ctx context.Context) ([]models.Account, error) {
	return s.repository.GetAccounts(ctx)
}

func (s *accountsService) GetPayments(ctx context.Context) ([]models.Payment, error) {
	transfers, err := s.repository.GetTransfers(ctx)
	if err != nil {
		return nil, err
	}

	var payments []models.Payment
	for _, transfer := range transfers {
		outgoing := models.Payment{
			Direction: models.OutgoingDirection,
			Account:   transfer.FromAccount,
			ToAccount: transfer.ToAccount,
			Amount:    transfer.Amount,
		}
		incoming := models.Payment{
			Direction:   models.IncomingDirection,
			Account:     transfer.ToAccount,
			FromAccount: transfer.FromAccount,
			Amount:      transfer.Amount,
		}
		payments = append(payments, outgoing, incoming)
	}
	return payments, nil
}

func (s *accountsService) CreateTransfer(ctx context.Context, transfer models.Transfer) error {
	return s.repository.AddTransfer(ctx, transfer)
}
