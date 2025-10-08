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
    metadata JSON,
    user     TEXT NOT NULL,
    FOREIGN KEY (user) REFERENCES user (login)
);

CREATE TABLE note
(
    name     TEXT PRIMARY KEY,
    text     TEXT,
    metadata JSON,
    user     TEXT NOT NULL,
    FOREIGN KEY (user) REFERENCES user (login)
);

CREATE TABLE binary
(
    name     TEXT PRIMARY KEY,
    filename TEXT NOT NULL,
    metadata JSON,
    user     TEXT NOT NULL,
    FOREIGN KEY (user) REFERENCES user (login)
);

CREATE TABLE card
(
    name       TEXT PRIMARY KEY,
    number     TEXT NOT NULL,
    exp_date   TEXT NOT NULL,
    cvv        TEXT NOT NULL,
    cardholder TEXT NOT NULL,
    metadata   JSON,
    user       TEXT NOT NULL,
    FOREIGN KEY (user) REFERENCES user (login)
);