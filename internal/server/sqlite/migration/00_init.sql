CREATE TABLE login
(
    name     TEXT PRIMARY KEY,
    login    TEXT NOT NULL,
    password TEXT,
    metadata JSON,
    user     TEXT NOT NULL
);

CREATE TABLE note
(
    name     TEXT PRIMARY KEY,
    text     TEXT,
    metadata JSON,
    user     TEXT NOT NULL
);

CREATE TABLE binary
(
    name     TEXT PRIMARY KEY,
    data     BLOB,
    path     TEXT NOT NULL,
    metadata JSON,
    user     TEXT   NOT NULL
);

CREATE TABLE card
(
    name       TEXT PRIMARY KEY,
    number     TEXT NOT NULL,
    exp_date   TEXT NOT NULL,
    cvv        TEXT NOT NULL,
    cardholder TEXT NOT NULL,
    metadata   JSON,
    user       TEXT NOT NULL
);