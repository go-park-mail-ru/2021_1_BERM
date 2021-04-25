DROP SCHEMA IF EXISTS user CASCADE;
CREATE SCHEMA user;


CREATE TABLE user.users
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

CREATE TABLE user.specializes
(
    id              SERIAL PRIMARY KEY NOT NULL,
    specialize_name VARCHAR UNIQUE     NOT NULL
);

CREATE TABLE user.user_specializes
(
    user_id       INTEGER NOT NULL,
    specialize_id INTEGER NOT NULL,
    FOREIGN KEY (user_id)
        REFERENCES user.users (id),
    FOREIGN KEY (specialize_id)
        REFERENCES user.specializes (id)
);

REATE TABLE user.user_reviews
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
