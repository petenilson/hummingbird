CREATE TABLE transfers (
  id SERIAL PRIMARY KEY,
  reason VARCHAR(255) NOT NULL,
  from_account_id int NOT NULL,
  to_account_id int NOT NULL,
  created_at TEXT NOT NULL,
  amount int NOT NULL,
  transaction_id int NOT NULL
);

CREATE TABLE transactions (
  id SERIAL PRIMARY KEY,
  created_at TEXT NOT NULL,
  description TEXT NOT NULL
);

CREATE TYPE entry_type AS ENUM ('DEBIT', 'CREDIT');

CREATE TABLE entrys (
  id SERIAL PRIMARY KEY,
  account_id INTEGER NOT NULL,
  transaction_id INTEGER NOT NULL,
  created_at TEXT NOT NULL,
  amount INTEGER NOT NULL,
  type entry_type NOT NULL
);

CREATE TABLE accounts (
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL,
  id SERIAL PRIMARY KEY,
  balance INTEGER NOT NULL,
  name VARCHAR(255) NOT NULL
);
