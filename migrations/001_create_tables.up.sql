-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS Users (
    user_id  UUID PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS Items (
    item_id SERIAL PRIMARY KEY,
    type TEXT UNIQUE NOT NULL,
    price INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS Balance (
    balance_id SERIAL PRIMARY KEY,
    user_id UUID REFERENCES Users(user_id) UNIQUE NOT NULL,
    coins_number INTEGER CONSTRAINT positive_coins_number CHECK (coins_number >= 0)
);

CREATE TABLE IF NOT EXISTS Inventory (
    buying_id SERIAL PRIMARY KEY,
    user_id UUID REFERENCES Users(user_id) NOT NULL,
    item_id INTEGER REFERENCES Items(item_id) NOT NULL,
    quantity INTEGER NOT NULL,
    create_date TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS Sent (
    sent_id SERIAL PRIMARY KEY,
    user_id UUID REFERENCES Users(user_id) NOT NULL,
    to_user_id UUID REFERENCES Users(user_id) NOT NULL,
    amount INTEGER NOT NULL,
    transaction_date TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS Received (
    received_id SERIAL PRIMARY KEY,
    user_id UUID REFERENCES Users(user_id) NOT NULL,
    from_user_id UUID REFERENCES Users(user_id) NOT NULL,
    amount INTEGER NOT NULL,
    transaction_date TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS name_idx ON Users (name);
CREATE INDEX IF NOT EXISTS item_id_idx ON Items (item_id);
CREATE INDEX IF NOT EXISTS type_idx ON Items (type);
CREATE INDEX IF NOT EXISTS balance_user_id_idx ON Balance (user_id);
CREATE INDEX IF NOT EXISTS inventory_user_id_idx ON Inventory (user_id);
CREATE UNIQUE INDEX IF NOT EXISTS unique_inventory_idx ON Inventory (user_id, item_id);
CREATE INDEX IF NOT EXISTS sent_user_id_idx ON Sent (user_id);
CREATE INDEX IF NOT EXISTS received_user_id_idx ON Received (user_id);

-- +goose StatementEnd

-- +goose StatementBegin

INSERT INTO items (type, price) VALUES
    ('t-shirt', 80),
    ('cup', 20),
    ('book', 50),
    ('pen', 10),
    ('powerbank', 200),
    ('hoody', 300),
    ('umbrella', 200),
    ('socks', 10),
    ('wallet', 50),
    ('pink-hoody', 500);

-- +goose StatementEnd