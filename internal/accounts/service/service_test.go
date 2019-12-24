package service_test

import (
	"context"
	"testing"

	"github.com/denperov/money-service/internal/accounts/models"
	"github.com/denperov/money-service/internal/accounts/service"
)

type fakeRepository struct {
	fnGetAccountRepository  func() service.AccountRepository
	fnGetTransferRepository func() service.TransferRepository
	fnWithTransaction       func(context.Context, func(service.Repository) error) error
}

func (f fakeRepository) GetAccountRepository() service.AccountRepository {
	return f.fnGetAccountRepository()
}

func (f fakeRepository) GetTransferRepository() service.TransferRepository {
	return f.fnGetTransferRepository()
}

func (f fakeRepository) WithTransaction(ctx context.Context, fn func(service.Repository) error) error {
	return f.fnWithTransaction(ctx, fn)
}

type fakeTransferRepository struct {
	fnGetTransfers func(context.Context) ([]models.Transfer, error)
	fnAddTransfer  func(context.Context, models.Transfer) error
}

func (f fakeTransferRepository) GetTransfers(ctx context.Context) ([]models.Transfer, error) {
	return f.fnGetTransfers(ctx)
}

func (f fakeTransferRepository) AddTransfer(ctx context.Context, t models.Transfer) error {
	return f.fnAddTransfer(ctx, t)
}

type fakeAccountRepository struct {
	fnGetAccounts       func(context.Context) ([]models.Account, error)
	fnGetAccount        func(context.Context, models.AccountID) (models.Account, bool, error)
	fnSetAccountBalance func(context.Context, models.AccountID, models.Money) error
}

func (f fakeAccountRepository) GetAccounts(ctx context.Context) ([]models.Account, error) {
	return f.fnGetAccounts(ctx)
}

func (f fakeAccountRepository) GetAccount(ctx context.Context, id models.AccountID) (models.Account, bool, error) {
	return f.fnGetAccount(ctx, id)
}

func (f fakeAccountRepository) SetAccountBalance(ctx context.Context, id models.AccountID, b models.Money) error {
	return f.fnSetAccountBalance(ctx, id, b)
}

func TestAccountsService_GetAccounts(t *testing.T) {
	ctx := context.Background()

	balance := models.MaxMoney()
	expectedAccount := models.Account{
		ID:       "",
		Currency: "",
		Balance:  balance,
	}

	accRep := fakeAccountRepository{
		fnGetAccounts: func(ctx context.Context) (accounts []models.Account, err error) {
			return []models.Account{expectedAccount}, nil
		},
	}
	rep := fakeRepository{
		fnGetAccountRepository: func() service.AccountRepository {
			return accRep
		},
	}
	srv := service.New(rep)

	accounts, err := srv.GetAccounts(ctx)
	if err != nil {
		t.Fail()
	}
	if len(accounts) != 1 {
		t.Fail()
	}
	if accounts[0] != expectedAccount {
		t.Fail()
	}
}

func TestAccountsService_GetPayments(t *testing.T) {
	ctx := context.Background()

	amount := models.MaxMoney()
	transfer := models.Transfer{
		FromAccount: "a",
		ToAccount:   "b",
		Amount:      amount,
	}
	expectedOutgoingPayment := models.Payment{
		Direction: models.OutgoingDirection,
		Account:   "a",
		ToAccount: "b",
		Amount:    amount,
	}
	expectedIncomingPayment := models.Payment{
		Direction:   models.IncomingDirection,
		Account:     "b",
		FromAccount: "a",
		Amount:      amount,
	}

	trRep := fakeTransferRepository{
		fnGetTransfers: func(ctx context.Context) (transfers []models.Transfer, err error) {
			return []models.Transfer{transfer}, nil
		},
	}
	rep := fakeRepository{
		fnGetTransferRepository: func() service.TransferRepository {
			return trRep
		},
	}
	srv := service.New(rep)

	payments, err := srv.GetPayments(ctx)
	if err != nil {
		t.Fail()
	}
	if len(payments) != 2 {
		t.Fail()
	}
	if payments[0].Direction == models.OutgoingDirection {
		if payments[0] != expectedOutgoingPayment || payments[1] != expectedIncomingPayment {
			t.Fail()
		}
	} else {
		if payments[0] != expectedIncomingPayment || payments[1] != expectedOutgoingPayment {
			t.Fail()
		}
	}
}

func TestAccountsService_CreateTransfer(t *testing.T) {
	ctx := context.Background()

	amount := models.MaxMoney()
	balanceA := models.MaxMoney()
	balanceB := models.ZeroMoney()
	expectedTransfer := models.Transfer{
		FromAccount: "a",
		ToAccount:   "b",
		Amount:      amount,
	}
	accountA := models.Account{
		ID:       "a",
		Currency: "USD",
		Balance:  balanceA,
	}
	accountB := models.Account{
		ID:       "b",
		Currency: "USD",
		Balance:  balanceB,
	}
	accRep := fakeAccountRepository{
		fnGetAccount: func(ctx context.Context, id models.AccountID) (account models.Account, b bool, err error) {
			switch id {
			case "a":
				return accountA, true, nil
			case "b":
				return accountB, true, nil
			default:
				t.Fail()
				return models.Account{}, false, nil
			}
		},
		fnSetAccountBalance: func(ctx context.Context, id models.AccountID, money models.Money) error {
			switch id {
			case "a":
				if money != models.ZeroMoney() {
					t.Fail()
				}
			case "b":
				if money != models.MaxMoney() {
					t.Fail()
				}
			default:
				t.Fail()
			}
			return nil
		},
	}
	trRep := fakeTransferRepository{
		fnAddTransfer: func(ctx context.Context, transfer models.Transfer) error {
			if transfer != expectedTransfer {
				t.Fail()
			}
			return nil
		},
	}
	tr := fakeRepository{
		fnGetAccountRepository: func() service.AccountRepository {
			return accRep
		},
		fnGetTransferRepository: func() service.TransferRepository {
			return trRep
		},
	}
	rep := fakeRepository{
		fnWithTransaction: func(ctx context.Context, f func(service.Repository) error) error {
			return f(tr)
		},
	}
	srv := service.New(rep)

	err := srv.CreateTransfer(ctx, expectedTransfer)
	if err != nil {
		t.Fail()
	}
}
