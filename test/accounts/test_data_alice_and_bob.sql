INSERT INTO accounts (public_id, currency, balance) VALUES
('bob123', 'USD', 100.00),
('alice456', 'USD', 0.01);

INSERT INTO transfers (account_from, account_to, amount) VALUES
(
    (select id from accounts where public_id = 'alice456'),
    (select id from accounts where public_id = 'bob123'),
    100.00
);
