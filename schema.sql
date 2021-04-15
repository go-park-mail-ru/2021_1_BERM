DROP SCHEMA IF EXISTS ff CASCADE;
CREATE SCHEMA ff;

CREATE TABLE ff.users
(
    id           SERIAL PRIMARY KEY NOT NULL,
    email        VARCHAR UNIQUE     NOT NULL,
    password     bytea              NOT NULL,
    login        VARCHAR            NOT NULL,
    name_surname VARCHAR            NOT NULL,
    about        VARCHAR DEFAULT NULL,
    executor     boolean            NOT NULL,
    img          VARCHAR DEFAULT '',
    rating       INTEGER DEFAULT 0
);

CREATE TABLE ff.specializes
(
    id              SERIAL PRIMARY KEY NOT NULL,
    specialize_name VARCHAR UNIQUE     NOT NULL
);

CREATE TABLE ff.user_specializes
(
    user_id       INTEGER NOT NULL,
    specialize_id INTEGER NOT NULL,
    FOREIGN KEY (user_id)
        REFERENCES ff.users (id),
    FOREIGN KEY (specialize_id)
        REFERENCES ff.specializes (id)
);

CREATE TABLE ff.orders
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

CREATE TABLE ff.vacancy
(
    id           SERIAL PRIMARY KEY NOT NULL,
    user_id      INTEGER            NOT NULL,
    category     VARCHAR            NOT NULL,
    vacancy_name VARCHAR            NOT NULL,
    description  VARCHAR            NOT NULL,
    salary       BIGINT             NOT NULL,
    FOREIGN KEY (user_id)
        REFERENCES ff.users (id)
);



CREATE TABLE ff.user_reviews
(
    id          SERIAL PRIMARY KEY NOT NULL,
    user_id     INTEGER            NOT NULL,
    order_id    INTEGER            NOT NULL,
    description VARCHAR            NOT NULL,
    executor    boolean            NOT NULL,
    FOREIGN KEY (user_id)
        REFERENCES ff.users (id),
    FOREIGN KEY (order_id)
        REFERENCES ff.orders (id)
);

CREATE TABLE ff.order_responses
(
    id         SERIAL PRIMARY KEY NOT NULL,
    order_id   INTEGER            NOT NULL,
    user_id    INTEGER            NOT NULL,
    rate       INTEGER            NOT NULL,
    user_login VARCHAR            NOT NULL,
    user_img   VARCHAR DEFAULT '',
    time       BIGINT             NOT NULL,
    FOREIGN KEY (user_id)
        REFERENCES ff.users (id),
    FOREIGN KEY (order_id)
        REFERENCES ff.orders (id)
);

CREATE TABLE ff.vacancy_responses
(
    id         SERIAL PRIMARY KEY NOT NULL,
    vacancy_id INTEGER            NOT NULL,
    user_id    INTEGER            NOT NULL,
    rate       INTEGER            NOT NULL,
    user_login VARCHAR            NOT NULL,
    user_img   VARCHAR DEFAULT '',
    time       BIGINT             NOT NULL,
    FOREIGN KEY (user_id)
        REFERENCES ff.users (id),
    FOREIGN KEY (vacancy_id)
        REFERENCES ff.vacancy (id)
);

SELECT array_agg(specialize_name) AS specializes FROM ff.specializes
		INNER JOIN ff.user_specializes us on specializes.id = us.specialize_id
		WHERE user_id = $1