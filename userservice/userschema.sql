ALTER USER postgres WITH ENCRYPTED PASSWORD 'admin';
DROP SCHEMA IF EXISTS userservice CASCADE;
CREATE EXTENSION IF NOT EXISTS citext;
CREATE SCHEMA userservice;

CREATE TABLE userservice.users
(
    id           SERIAL PRIMARY KEY NOT NULL,
    email        VARCHAR UNIQUE     NOT NULL,
    password     bytea              NOT NULL,
    login        VARCHAR            NOT NULL,
    name_surname citext            NOT NULL,
    about        VARCHAR DEFAULT NULL,
    executor     boolean            NOT NULL,
    img          VARCHAR DEFAULT ''
);

CREATE TABLE userservice.specializes
(
    id              SERIAL PRIMARY KEY NOT NULL,
    specialize_name VARCHAR UNIQUE     NOT NULL
);

CREATE TABLE userservice.user_specializes
(
    user_id       INTEGER NOT NULL,
    specialize_id INTEGER NOT NULL,
    FOREIGN KEY (user_id)
        REFERENCES userservice.users (id),
    FOREIGN KEY (specialize_id)
        REFERENCES userservice.specializes (id),
    PRIMARY KEY (user_id, specialize_id)
);

CREATE TABLE userservice.reviews
(
    id          SERIAL PRIMARY KEY NOT NULL,
    user_id     INTEGER            NOT NULL,
    to_user_id  INTEGER            NOT NULL,
    order_id    INTEGER            NOT NULL,
    description VARCHAR            NOT NULL,
    score       INTEGER            NOT NULL,
    UNIQUE (user_id, to_user_id, order_id),
    FOREIGN KEY (user_id)
        REFERENCES userservice.users (id),
    FOREIGN KEY (to_user_id)
        REFERENCES userservice.users (id)

);
