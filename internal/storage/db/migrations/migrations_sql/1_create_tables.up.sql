CREATE TABLE Users (
    user_id  UUID PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL
);

CREATE TABLE Items (
    item_id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    price INTEGER NOT NULL
);

CREATE TABLE Balance (
    user_id UUID PRIMARY KEY,
    coins_number INTEGER CONSTRAINT positive_coins_number CHECK (coins_number >= 0)
);