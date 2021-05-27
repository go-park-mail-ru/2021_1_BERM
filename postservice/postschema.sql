ALTER USER postgres WITH ENCRYPTED PASSWORD 'admin';
DROP SCHEMA IF EXISTS post CASCADE;
CREATE SCHEMA post;

CREATE TABLE post.orders
(
    id          SERIAL PRIMARY KEY    NOT NULL,
    customer_id INTEGER               NOT NULL,
    executor_id INTEGER,
    order_name  VARCHAR               NOT NULL,
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
    order_name  VARCHAR               NOT NULL,
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
    vacancy_name VARCHAR               NOT NULL,
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
    vacancy_name VARCHAR               NOT NULL,
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


SELECT *
FROM post.orders
WHERE CASE budget != 0 THEN budget >= 300 AND budget =< 400 ELSE true END
AND CASE search_str != "~" THEM to_tsvector(description) @@ to_tsquery(search_str) ELSE true END
AND CASE category$$ != "~" THEN category = category$$ ELSE true END
ORDER BY budget
LIMIT 1 OFFSET 25

SELECT *
FROM post.orders
WHERE CASE budget != 0 THEN budget >= 300 AND budget =< 400 ELSE true END
AND CASE search_str != "~" THEM to_tsvector(description) @@ to_tsquery(search_str) ELSE true END
AND CASE category$$ != "~" THEN category = category$$ ELSE true END
ORDER BY budget DESC
LIMIT 1 OFFSET 25;


SELECT *
FROM userservice.users AS users
WHERE CASE WHEN 1 != 0 THEN (SELECT AVG(score) FROM userservice.reviews WHERE to_user_id = users.id) >= 1 ELSE true END
  AND CASE WHEN 4.5 != 0 THEN (SELECT AVG(score) FROM userservice.reviews WHERE to_user_id = users.id) <= 4.5 ELSE true END
  AND CASE WHEN '19r156' != '~' THEN to_tsvector(login) @@ to_tsquery('19r156') ELSE true END
ORDER BY (SELECT AVG(score) FROM userservice.reviews WHERE to_user_id = users.id) LIMIT 10 OFFSET 0;