CREATE TABLE user
(
    login    TEXT PRIMARY KEY,
    password TEXT NOT NULL
);

CREATE TABLE login
(
    name     TEXT PRIMARY KEY,
    login    TEXT NOT NULL,
    password TEXT,
    website  TEXT,
    notes    TEXT,
    user     TEXT NOT NULL,
    FOREIGN KEY (user) REFERENCES user (login)
);

CREATE TABLE note
(
    name TEXT PRIMARY KEY,
    text TEXT,
    user TEXT NOT NULL,
    FOREIGN KEY (user) REFERENCES user (login)
);

CREATE TABLE binary
(
    name     TEXT PRIMARY KEY,
    filename TEXT    NOT NULL,
    size     INTEGER NOT NULL,
    notes    TEXT,
    user     TEXT    NOT NULL,
    FOREIGN KEY (user) REFERENCES user (login)
);

CREATE TABLE card
(
    name       TEXT PRIMARY KEY,
    number     TEXT NOT NULL,
    exp_date   TEXT NOT NULL,
    cvv        TEXT NOT NULL,
    cardholder TEXT NOT NULL,
    notes      TEXT,
    user       TEXT NOT NULL,
    FOREIGN KEY (user) REFERENCES user (login)
);