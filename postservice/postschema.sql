ALTER USER postgres WITH ENCRYPTED PASSWORD 'admin';
DROP SCHEMA IF EXISTS post CASCADE;
CREATE EXTENSION IF NOT EXISTS citext;
CREATE SCHEMA post;

CREATE TABLE post.orders
(
    id          SERIAL PRIMARY KEY    NOT NULL,
    customer_id INTEGER               NOT NULL,
    executor_id INTEGER,
    order_name  citext                NOT NULL,
    category    VARCHAR               NOT NULL,
    budget      BIGINT                NOT NULL,
    deadline    BIGINT                NOT NULL,
    description VARCHAR               NOT NULL,
    is_archived BOOLEAN DEFAULT FALSE NOT NULL
);

CREATE TABLE post.archive_orders
(
    id          INTEGER PRIMARY KEY   NOT NULL,
    customer_id INTEGER               NOT NULL,
    executor_id INTEGER,
    order_name  citext                NOT NULL,
    category    VARCHAR               NOT NULL,
    budget      BIGINT                NOT NULL,
    deadline    BIGINT                NOT NULL,
    description VARCHAR               NOT NULL,
    is_archived BOOLEAN DEFAULT FALSE NOT NULL
);

CREATE TABLE post.vacancy
(
    id           SERIAL PRIMARY KEY    NOT NULL,
    customer_id  INTEGER               NOT NULL,
    executor_id  INTEGER DEFAULT 0,
    category     VARCHAR               NOT NULL,
    vacancy_name citext                NOT NULL,
    description  VARCHAR               NOT NULL,
    salary       BIGINT                NOT NULL,
    is_archived  BOOLEAN DEFAULT FALSE NOT NULL

);

CREATE TABLE post.archive_vacancy
(
    id           INTEGER PRIMARY KEY   NOT NULL,
    customer_id  INTEGER               NOT NULL,
    executor_id  INTEGER DEFAULT 0,
    category     VARCHAR               NOT NULL,
    vacancy_name citext                NOT NULL,
    description  VARCHAR               NOT NULL,
    salary       BIGINT                NOT NULL,
    is_archived  BOOLEAN DEFAULT FALSE NOT NULL

);

CREATE TABLE post.responses
(
    id               SERIAL PRIMARY KEY    NOT NULL,
    post_id          INTEGER               NOT NULL,
    user_id          INTEGER               NOT NULL,
    rate             INTEGER DEFAULT 0,
    text             VARCHAR DEFAULT '',
    order_response   BOOLEAN DEFAULT FALSE NOT NULL,
    vacancy_response BOOLEAN DEFAULT FALSE NOT NULL,
    time             BIGINT                NOT NULL
);
