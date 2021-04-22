DROP SCHEMA IF EXISTS post CASCADE;
CREATE SCHEMA post;

CREATE TABLE post.orders
(
    id          SERIAL PRIMARY KEY NOT NULL,
    customer_id INTEGER            NOT NULL,
    executor_id INTEGER,
    order_name  VARCHAR            NOT NULL,
    category    VARCHAR            NOT NULL,
    budget      BIGINT             NOT NULL,
    deadline    BIGINT             NOT NULL,
    description VARCHAR            NOT NULL
);

CREATE TABLE post.vacancy
(
    id           SERIAL PRIMARY KEY NOT NULL,
    user_id      INTEGER            NOT NULL,
    category     VARCHAR            NOT NULL,
    vacancy_name VARCHAR            NOT NULL,
    description  VARCHAR            NOT NULL,
    salary       BIGINT             NOT NULL
);

CREATE TABLE post.responses
(
    id      SERIAL PRIMARY KEY NOT NULL,
    post_id INTEGER            NOT NULL,
    user_id INTEGER            NOT NULL,
    rate    INTEGER            NOT NULL,
    time    BIGINT             NOT NULL
);
