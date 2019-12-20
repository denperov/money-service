package postgres_repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/denperov/money-service/internal/accounts/service"
	"github.com/denperov/money-service/internal/pkg/user_errors"
)

type postgresRepository struct {
	url  string
	pool *pgxpool.Pool
}

func New(
	address string,
	name string,
	user string,
	password string,
) *postgresRepository {
	return &postgresRepository{
		url: fmt.Sprintf("postgresql://%s:%s@%s/%s?pool_max_conns=3", user, password, address, name),
	}
}

func (r *postgresRepository) Start(ctx context.Context) error {
	if r.pool != nil {
		return nil
	}
	pool, err := connectWithRetries(ctx, r.url)
	if err != nil {
		return err
	}
	r.pool = pool

	return nil
}

func (r *postgresRepository) Stop() {
	if r.pool == nil {
		return
	}
	r.pool.Close()
}

func connectWithRetries(ctx context.Context, url string) (*pgxpool.Pool, error) {
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, err
	}

	for {
		pool, err := pgxpool.ConnectConfig(ctx, cfg)
		if err != nil {
			// wait
			select {
			case <-ctx.Done(): // cancellation
				return nil, ctx.Err()
			case <-ticker.C:
				continue
			}
		}
		return pool, nil
	}
}

const queryGetAccounts = `
select
	public_id,
	currency,
	balance::text
from accounts;
`

func (r *postgresRepository) GetAccounts(ctx context.Context) ([]service.Account, error) {
	rows, err := r.pool.Query(ctx, queryGetAccounts)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []service.Account
	for rows.Next() {
		select {
		case <-ctx.Done(): // cancellation
			return nil, ctx.Err()
		default:
		}

		var account service.Account
		err = rows.Scan(&account.ID, &account.Currency, &account.Balance)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
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

func (r *postgresRepository) GetTransfers(ctx context.Context) ([]service.Transfer, error) {
	rows, err := r.pool.Query(ctx, queryGetPayments)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transfers []service.Transfer
	for rows.Next() {
		select {
		case <-ctx.Done(): // cancellation
			return nil, ctx.Err()
		default:
		}

		var transfer service.Transfer
		err = rows.Scan(&transfer.FromAccount, &transfer.ToAccount, &transfer.Amount)
		if err != nil {
			return nil, err
		}

		transfers = append(transfers, transfer)
	}
	return transfers, nil
}

const queryAddTransferCheckAccounts = `
select (select currency from accounts where public_id = $1), (select currency from accounts where public_id = $2);
`
const queryAddTransferInsert = `
insert into transfers (account_from, account_to, amount)
select
	(select id from accounts where public_id = $1),
	(select id from accounts where public_id = $2),
	$3::numeric(13,2);
`
const queryAddTransferUpdateFrom = `
update accounts
set balance = balance - $2::numeric(13,2)
where public_id = $1;
`
const queryAddTransferUpdateTo = `
update accounts
set balance = balance + $2::numeric(13,2)
where public_id = $1;
`

func (r *postgresRepository) AddTransfer(ctx context.Context, transfer service.Transfer) (err error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted, AccessMode: pgx.ReadWrite})
	if err != nil {
		return fmt.Errorf("begin database transaction: %w", err)
	}
	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback(ctx)
			if rollbackErr != nil {
				log.Printf("transaction rollback: %v", rollbackErr)
			}
		}
	}()

	row := r.pool.QueryRow(ctx, queryAddTransferCheckAccounts, transfer.FromAccount, transfer.ToAccount)

	var currencyFrom sql.NullString
	var currencyTo sql.NullString
	err = row.Scan(&currencyFrom, &currencyTo)
	if err != nil {
		return err
	}
	if !currencyFrom.Valid {
		return user_errors.ErrorWrongSourceAccount
	}
	if !currencyTo.Valid {
		return user_errors.ErrorWrongDestinationAccount
	}
	if currencyFrom != currencyTo {
		return user_errors.ErrorDifferentCurrencies
	}

	_, err = tx.Exec(ctx, queryAddTransferInsert, transfer.FromAccount, transfer.ToAccount, transfer.Amount)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.ConstraintName == "positive_amount" {
			return user_errors.ErrorWrongAmount
		}
		return err
	}
	_, err = tx.Exec(ctx, queryAddTransferUpdateFrom, transfer.FromAccount, transfer.Amount)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.ConstraintName == "positive_balance" {
			return user_errors.ErrorNotEnoughMoney
		}
		return err
	}
	_, err = tx.Exec(ctx, queryAddTransferUpdateTo, transfer.ToAccount, transfer.Amount)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
