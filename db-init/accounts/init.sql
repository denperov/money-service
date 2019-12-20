CREATE TABLE accounts (
    id         serial not null primary key,
    public_id  varchar(255) not null unique,
    currency   varchar(3) not null,
    balance    numeric(13,2) CONSTRAINT positive_balance check (balance >= 0),
    created_at timestamptz default now()
);

CREATE TABLE transfers (
    id           serial not null primary key,
    account_from integer references accounts (id),
    account_to   integer references accounts (id),
    amount       numeric(13,2) CONSTRAINT positive_amount check (amount > 0),
    created_at   timestamptz default now()
);
