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
  amount INTEGER NOT NULL,
  account_id INTEGER NOT NULL,
  type entry_type NOT NULL,
  created_at TEXT NOT NULL,
  CONSTRAINT fk_account FOREIGN KEY (account_id) REFERENCES accounts (id)
);

CREATE TABLE transaction_entrys (
  id SERIAL PRIMARY KEY,
  entry_id INTEGER NOT NULL,
  transaction_id INTEGER NOT NULL,
  created_at TEXT NOT NULL,
  CONSTRAINT fk_entry FOREIGN KEY (entry_id) REFERENCES entrys (id),
  CONSTRAINT fk_transaction FOREIGN KEY (transaction_id) REFERENCES transactions (id),
  CONSTRAINT unique_entry_transaction UNIQUE (entry_id, transaction_id)
);

CREATE INDEX idx_tx_entry_tx_id ON transaction_entrys (transaction_id);
