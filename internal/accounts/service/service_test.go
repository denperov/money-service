package service

import (
	"context"
	"testing"
)

type fakeRepository struct {
}

func (f fakeRepository) GetAccounts(context.Context) ([]Account, error) {
	panic("implement me")
}

func (f fakeRepository) GetTransfers(context.Context) ([]Transfer, error) {
	panic("implement me")
}

func (f fakeRepository) AddTransfer(context.Context, Transfer) error {
	panic("implement me")
}

func TestAccountsService_GetAccounts(t *testing.T) {
	rep := fakeRepository{}
	srv := New(rep)

	ctx := context.Background()
	accounts, err := srv.GetAccounts(ctx)
	if err != nil {

	}

}

func TestAccountsService_GetPayments(t *testing.T) {

}

func TestAccountsService_SendPayment(t *testing.T) {

}
