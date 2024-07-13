CREATE TABLE accounts (
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL,
  id SERIAL PRIMARY KEY,
  balance INTEGER NOT NULL,
  name VARCHAR(255) NOT NULL
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
  created_at TEXT NOT NULL,
  amount INTEGER NOT NULL,
  type entry_type NOT NULL,
  CONSTRAINT fk_account FOREIGN KEY (account_id) REFERENCES accounts (id)
);

CREATE TABLE transfers (
  id SERIAL PRIMARY KEY,
  transaction_id int NOT NULL,
  amount int NOT NULL,
  from_account_id int NOT NULL,
  to_account_id int NOT NULL,
  reason VARCHAR(255) NOT NULL,
  created_at TEXT NOT NULL,
  CONSTRAINT fk_to_account FOREIGN KEY (to_account_id) REFERENCES accounts (id),
  CONSTRAINT fk_from_account FOREIGN KEY (from_account_id) REFERENCES accounts (id),
  CONSTRAINT fk_transaction FOREIGN KEY (transaction_id) REFERENCES transactions (id)
);

CREATE TABLE transaction_entrys (
  id SERIAL PRIMARY KEY,
  entry_id INTEGER NOT NULL,
  transaction_id INTEGER NOT NULL,
  created_at TEXT NOT NULL,
  CONSTRAINT fk_entry FOREIGN KEY (entry_id) REFERENCES entrys (id),
  CONSTRAINT fk_transaction FOREIGN KEY (transaction_id) REFERENCES transactions (id)
);

CREATE INDEX idx_tx_entry_tx_id ON transaction_entrys (transaction_id);
