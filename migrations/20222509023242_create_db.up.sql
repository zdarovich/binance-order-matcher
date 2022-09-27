CREATE TABLE IF NOT EXISTS books(
    order_id TEXT NOT NULL,
    update_id INTEGER NOT NULL,
    symbol TEXT NOT NULL,
    bid_price REAL NOT NULL,
    bid_quantity REAL NOT NULL,
    ask_price REAL NOT NULL,
    ask_quantity REAL NOT NULL,
    created_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS orders(
     id TEXT NOT NULL,
     order_type INTEGER NOT NULL,
     symbol TEXT NOT NULL,
     price REAL NOT NULL,
     quantity REAL NOT NULL,
     created_at TEXT NOT NULL
);