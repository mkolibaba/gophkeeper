INSERT INTO user (login, password)
VALUES ('demo', 'demo');
-- TODO: не хранить пароли в открытом доступе

INSERT INTO login (name, login, password, metadata, user)
VALUES ('Google', 'iivanov', '123456', '{"Website": "google.com"}', 'demo'),
       ('Ozon', '79031002030', 'qwe123', '{"Tag": "shop"}', 'demo'),
       ('Wildberries', '79031002030', 'qwerty', null, 'demo'),
       ('Госуслуги', 'iivanov@gmail.com', 'qweasd', null, 'demo'),
       ('Mail', 'ivanivanov', 'qqwwee', null, 'demo'),
       ('VK', 'ivanivanov@mail.ru', 'qweqwe', null, 'demo');

INSERT INTO note (name, text, metadata, user)
VALUES ('Записки о природе',
        'Кто никогда не видал, как растет клюква, тот может очень долго идти по болоту и не замечать, что он по клюкве идет.',
        null, 'demo'),
       ('Мысль',
        'Живешь ты, может быть, сам триста лет, и кто породил тебя, тот в яичке своем пересказал все, что он тоже узнал за свои триста лет жизни.',
        null, 'demo');

INSERT INTO binary (name, filename, metadata, user)
VALUES ('Squirtle',  'squirtle_pokemon.png', '{"Tag": "pokemon"}', 'demo');

INSERT INTO card (name, number, exp_date, cvv, cardholder, metadata, user)
VALUES ('Сбербанк', '2200123456789019', '03/28', '541', 'IVAN IVANOV', null, 'demo'),
       ('Т-Банк', '2201987654321000', '01/29', '192', 'IVAN IVANOV', null, 'demo');
