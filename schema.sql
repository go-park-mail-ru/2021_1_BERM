drop table if exists users cascade;
drop table if exists specializes cascade;
drop table if exists user_specializes cascade;
drop table if exists orders cascade;
drop table if exists order_specializes cascade;
drop table if exists user_reviews cascade;
drop table if exists vacancy cascade;
drop table if exists responses cascade;

CREATE TABLE users
(
    id           SERIAL PRIMARY KEY NOT NULL,
    email        VARCHAR UNIQUE     NOT NULL,
    password     VARCHAR            NOT NULL,
    login        VARCHAR            NOT NULL,
    name_surname VARCHAR            NOT NULL,
    about        VARCHAR DEFAULT NULL,
    executor     boolean            NOT NULL,
    img          VARCHAR DEFAULT '',
    rating       INTEGER DEFAULT 0
);

CREATE TABLE specializes
(
    id              SERIAL PRIMARY KEY NOT NULL,
    specialize_name VARCHAR UNIQUE     NOT NULL
);

CREATE TABLE user_specializes
(
    user_id       INTEGER NOT NULL,
    specialize_id INTEGER NOT NULL,
    FOREIGN KEY (user_id)
        REFERENCES users (id),
    FOREIGN KEY (specialize_id)
        REFERENCES specializes (id)
);

CREATE TABLE orders
(
    id          SERIAL PRIMARY KEY NOT NULL,
    customer_id INTEGER            NOT NULL,
    executor_id INTEGER,
    order_name  VARCHAR            NOT NULL,
    category    VARCHAR            NOT NULL,
    budget      INTEGER            NOT NULL,
    deadline    BIGINT             NOT NULL,
    description VARCHAR            NOT NULL
);

CREATE TABLE vacancy
(
    id           SERIAL PRIMARY KEY NOT NULL,
    category     VARCHAR            NOT NULL,
    vacancy_name VARCHAR            NOT NULL,
    description  VARCHAR            NOT NULL,
    salary       INTEGER            NOT NULL
);



CREATE TABLE user_reviews
(
    id          SERIAL PRIMARY KEY NOT NULL,
    user_id     INTEGER            NOT NULL,
    order_id    INTEGER            NOT NULL,
    description VARCHAR            NOT NULL,
    executor    boolean            NOT NULL,
    FOREIGN KEY (user_id)
        REFERENCES users (id),
    FOREIGN KEY (order_id)
        REFERENCES orders (id)
);

CREATE TABLE responses
(
    id            SERIAL PRIMARY KEY NOT NULL,
    order_id      INTEGER            NOT NULL,
    user_id       INTEGER            NOT NULL,
    rate          INTEGER            NOT NULL,
    user_login VARCHAR            NOT NULL,
    user_img      VARCHAR DEFAULT '',
    FOREIGN KEY (user_id)
        REFERENCES users (id),
    FOREIGN KEY (order_id)
        REFERENCES orders (id),
    FOREIGN KEY (user_login)
        REFERENCES users (login),
    FOREIGN KEY (user_img)
        REFERENCES users (img)
);


-- SELECT users.*, array_agg(specialize_name) AS specializes from users
-- INNER JOIN user_specializes ON users.id = user_specializes.user_id
-- INNER JOIN specializes ON user_specializes.specialize_id = specializes.id
-- WHERE users.email = 'kek@mem.ru'
-- GROUP BY users.id

-- SELECT * from orders
-- WHERE id = 1
