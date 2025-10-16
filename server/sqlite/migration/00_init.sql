CREATE TABLE user
(
    login    TEXT PRIMARY KEY,
    password TEXT NOT NULL
);

CREATE TABLE login
(
    id       INTEGER PRIMARY KEY AUTOINCREMENT,
    name     TEXT NOT NULL,
    login    TEXT NOT NULL,
    password TEXT,
    website  TEXT,
    notes    TEXT,
    user     TEXT NOT NULL,
    FOREIGN KEY (user) REFERENCES user (login)
);

CREATE TABLE note
(
    id   INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    text TEXT,
    user TEXT NOT NULL,
    FOREIGN KEY (user) REFERENCES user (login)
);

CREATE TABLE binary
(
    id       INTEGER PRIMARY KEY AUTOINCREMENT,
    name     TEXT NOT NULL,
    filename TEXT    NOT NULL,
    size     INTEGER NOT NULL,
    notes    TEXT,
    user     TEXT    NOT NULL,
    FOREIGN KEY (user) REFERENCES user (login)
);

CREATE TABLE card
(
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    name       TEXT NOT NULL,
    number     TEXT NOT NULL,
    exp_date   TEXT NOT NULL,
    cvv        TEXT NOT NULL,
    cardholder TEXT NOT NULL,
    notes      TEXT,
    user       TEXT NOT NULL,
    FOREIGN KEY (user) REFERENCES user (login)
);