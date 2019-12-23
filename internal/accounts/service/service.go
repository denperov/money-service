package service

import (
	"context"

	"github.com/denperov/money-service/internal/accounts/models"
	"github.com/denperov/money-service/internal/pkg/user_errors"
)

type Repository interface {
	GetAccountRepository() AccountRepository
	GetTransferRepository() TransferRepository

	WithTransaction(context.Context, func(Repository) error) error
}

type AccountRepository interface {
	GetAccounts(context.Context) ([]models.Account, error)
	GetAccount(context.Context, models.AccountID) (models.Account, bool, error)
	SetAccountBalance(context.Context, models.AccountID, models.Money) error
}

type TransferRepository interface {
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
	return s.repository.GetAccountRepository().GetAccounts(ctx)
}

func (s *accountsService) GetPayments(ctx context.Context) ([]models.Payment, error) {
	transfers, err := s.repository.GetTransferRepository().GetTransfers(ctx)
	if err != nil {
		return nil, err
	}

	payments := make([]models.Payment, 0, len(transfers)*2)
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
	return s.repository.WithTransaction(ctx, func(r Repository) error {

		accountRepository := r.GetAccountRepository()
		transferRepository := r.GetTransferRepository()

		if transfer.FromAccount == transfer.ToAccount {
			return user_errors.ErrorSameAccount
		}

		if !models.MoneyIsPositive(transfer.Amount) {
			return user_errors.ErrorWrongAmount
		}

		accountFrom, ok, err := accountRepository.GetAccount(ctx, transfer.FromAccount)
		if err != nil {
			return err
		}
		if !ok {
			return user_errors.ErrorWrongSourceAccount
		}

		accountTo, ok, err := accountRepository.GetAccount(ctx, transfer.ToAccount)
		if err != nil {
			return err
		}
		if !ok {
			return user_errors.ErrorWrongDestinationAccount
		}

		if accountFrom.Currency != accountTo.Currency {
			return user_errors.ErrorDifferentCurrencies
		}

		if models.MoneyLess(accountFrom.Balance, transfer.Amount) {
			return user_errors.ErrorNotEnoughMoney
		}

		fromBalance, err := models.MoneyDiff(accountFrom.Balance, transfer.Amount)
		if err != nil {
			return err
		}

		toBalance, err := models.MoneySum(accountTo.Balance, transfer.Amount)
		if err != nil {
			return err
		}

		err = transferRepository.AddTransfer(ctx, transfer)
		if err != nil {
			return err
		}

		err = accountRepository.SetAccountBalance(ctx, accountFrom.ID, fromBalance)
		if err != nil {
			return err
		}

		err = accountRepository.SetAccountBalance(ctx, accountTo.ID, toBalance)
		if err != nil {
			return err
		}

		return nil
	})
}
