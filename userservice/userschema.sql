ALTER USER postgres WITH ENCRYPTED PASSWORD 'admin';
DROP SCHEMA IF EXISTS userservice CASCADE;
CREATE SCHEMA userservice;

CREATE TABLE userservice.users
(
    id           SERIAL PRIMARY KEY NOT NULL,
    email        VARCHAR UNIQUE     NOT NULL,
    password     bytea              NOT NULL,
    login        VARCHAR            NOT NULL,
    name_surname VARCHAR            NOT NULL,
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



SELECT id, email, password, login, name_surname, about, executor, img, AVG(score) AS rating
FROM userservice.users AS users
INNER JOIN userservice.reviews
ON users.id = reviews.to_user_id
WHERE CASE WHEN $1 != 0 THEN (SELECT AVG(score) FROM userservice.reviews WHERE to_user_id = users.id) >= $1 ELSE true END
  AND CASE WHEN $2 != 0 THEN (SELECT AVG(score) FROM userservice.reviews WHERE to_user_id = users.id) <= $2 ELSE true END
  AND CASE WHEN $3 != '~' THEN to_tsvector(name_surname) @@ to_tsquery($3)) ELSE true END
ORDER BY rating LIMIT $5 OFFSET $6

SELECT id, email, password, login, name_surname, about, executor, img, AVG(score) AS rating
FROM userservice.users AS users
         INNER JOIN userservice.reviews
                    ON users.id = reviews.to_user_id
WHERE CASE WHEN $1 != 0 THEN (SELECT AVG(score) FROM userservice.reviews WHERE to_user_id = users.id) >= $1 ELSE true END
  AND CASE WHEN $2 != 0 THEN (SELECT AVG(score) FROM userservice.reviews WHERE to_user_id = users.id) <= $2 ELSE true END
  AND CASE WHEN $3 != '~' THEN to_tsvector(name_surname) @@ to_tsquery($3)) ELSE true END
ORDER BY rating DESC LIMIT $5 OFFSET $6

SELECT id, email, password, login, name_surname, about, executor, img, AVG(score) AS rating
FROM userservice.users AS users
INNER JOIN userservice.reviews
 ON users.id = reviews.to_user_id
WHERE CASE WHEN $1 != 0 THEN (SELECT AVG(score) FROM userservice.reviews WHERE to_user_id = users.id) >= $1 ELSE true END
  AND CASE WHEN $2 != 0 THEN (SELECT AVG(score) FROM userservice.reviews WHERE to_user_id= users.id) <= $2 ELSE true END
  AND CASE WHEN $3 != '~' THEN to_tsvector(name_surname) @@ to_tsquery($3)) ELSE true END
ORDER BY name_surname DESC LIMIT $5 OFFSET $6

SELECT id, email, password, login, name_surname, about, executor, img, AVG(score) AS rating
FROM userservice.users AS users
         INNER JOIN userservice.reviews
                    ON users.id = reviews.to_user_id
WHERE CASE WHEN $1 != 0 THEN (SELECT AVG(score) FROM userservice.reviews WHERE to_user_id = users.id) >= $1 ELSE true END
  AND CASE WHEN $2 != 0 THEN (SELECT AVG(score) FROM userservice.reviews WHERE to_user_id= users.id) <= $2 ELSE true END
  AND CASE WHEN $3 != '~' THEN to_tsvector(name_surname) @@ to_tsquery($3)) ELSE true END
ORDER BY name_surname LIMIT $5 OFFSET $6

    SELECT users.id as id , email, password, login, name_surname, about, executor, img, AVG(score) AS rating
    FROM userservice.users AS users
    INNER JOIN userservice.reviews
     ON users.id = reviews.to_user_id
    WHERE CASE WHEN 4 != 0 THEN (SELECT AVG(score) FROM userservice.reviews WHERE to_user_id = users.id) >= 4 ELSE true END
    AND CASE WHEN 4 != 0 THEN (SELECT AVG(score) FROM userservice.reviews WHERE to_user_id = users.id) <= 4 ELSE true END
    AND CASE WHEN '~' != '~' THEN to_tsvector(name_surname) @@ to_tsquery('~') ELSE true END
    GROUP BY users.id
    ORDER BY rating LIMIT 10 OFFSET 0;