package postgres_repository

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"

	"github.com/denperov/money-service/internal/accounts/models"
	"github.com/denperov/money-service/internal/accounts/service"
)

type queryInterface interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
}

type postgresRepository struct {
	qi queryInterface
}

func (r *postgresRepository) WithTransaction(ctx context.Context, fn func(service.Repository) error) error {
	tx, err := r.qi.Begin(ctx)
	if err != nil {
		return err
	}

	err = fn(&postgresRepository{qi: tx})
	if err != nil {
		rollbackErr := tx.Rollback(ctx)
		if rollbackErr != nil {
			log.Printf("transaction rollback: %v", rollbackErr)
		}
		return err
	}

	return tx.Commit(ctx)
}

func (r *postgresRepository) GetAccountRepository() service.AccountRepository {
	return r
}

func (r *postgresRepository) GetTransferRepository() service.TransferRepository {
	return r
}

const queryGetAccounts = `
select
	public_id,
	currency,
	balance::text
from accounts;
`

func (r *postgresRepository) GetAccounts(ctx context.Context) ([]models.Account, error) {
	rows, err := r.qi.Query(ctx, queryGetAccounts)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []models.Account
	for rows.Next() {
		select {
		case <-ctx.Done(): // cancellation
			return nil, ctx.Err()
		default:
		}

		var account models.Account
		var balanceString string
		err = rows.Scan(&account.ID, &account.Currency, &balanceString)
		if err != nil {
			return nil, err
		}

		account.Balance, err = models.MoneyFromString(balanceString)
		if err != nil {
			return nil, err
		}

		accounts = append(accounts, account)
	}
	return accounts, nil
}

const queryGetAccount = `
select
	public_id,
	currency,
	balance::text
from accounts
where public_id = $1;
`

func (r *postgresRepository) GetAccount(ctx context.Context, id models.AccountID) (models.Account, bool, error) {
	var account models.Account
	var balanceString string
	err := r.qi.QueryRow(ctx, queryGetAccount, id).Scan(&account.ID, &account.Currency, &balanceString)
	if err == pgx.ErrNoRows {
		return models.Account{}, false, nil
	}
	if err != nil {
		return models.Account{}, false, err
	}

	account.Balance, err = models.MoneyFromString(balanceString)
	if err != nil {
		return models.Account{}, false, err
	}

	return account, true, nil
}

const querySetAccountBalance = `
update accounts
set balance = $2::numeric(16,2)
where public_id = $1;
`

func (r *postgresRepository) SetAccountBalance(ctx context.Context, id models.AccountID, money models.Money) error {
	t, err := r.qi.Exec(ctx, querySetAccountBalance, id, money.String())
	if err != nil {
		return err
	}
	if t.RowsAffected() != 1 {
		return errors.New("record not found")
	}
	return nil
}

const queryGetPayments = `
select
	af.public_id,
	at.public_id,
	t.amount::text
from transfers t
join accounts af on af.id = t.account_from
join accounts at on at.id = t.account_to;
`

func (r *postgresRepository) GetTransfers(ctx context.Context) ([]models.Transfer, error) {
	rows, err := r.qi.Query(ctx, queryGetPayments)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transfers []models.Transfer
	for rows.Next() {
		select {
		case <-ctx.Done(): // cancellation
			return nil, ctx.Err()
		default:
		}

		var transfer models.Transfer
		var amountString string
		err = rows.Scan(&transfer.FromAccount, &transfer.ToAccount, &amountString)
		if err != nil {
			return nil, err
		}

		transfer.Amount, err = models.MoneyFromString(amountString)
		if err != nil {
			return nil, err
		}

		transfers = append(transfers, transfer)
	}
	return transfers, nil
}

const queryAddTransferInsert = `
insert into transfers (account_from, account_to, amount)
select
	(select id from accounts where public_id = $1),
	(select id from accounts where public_id = $2),
	$3::numeric(16,2);
`

func (r *postgresRepository) AddTransfer(ctx context.Context, transfer models.Transfer) (err error) {
	t, err := r.qi.Exec(ctx, queryAddTransferInsert, transfer.FromAccount, transfer.ToAccount, transfer.Amount.String())
	if err != nil {
		return err
	}
	if t.RowsAffected() != 1 {
		return errors.New("record not found")
	}
	return nil
}
