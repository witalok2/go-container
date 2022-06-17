BEGIN;

CREATE TABLE IF NOT EXISTS users
(
	id varchar(255) primary key,
	username varchar(255)
);

CREATE TABLE IF NOT EXISTS migrations
(
    name character varying(100) NOT NULL PRIMARY KEY,
    applied_at timestamp without time zone NOT NULL
);

INSERT INTO migrations VALUES ('00001_create_table_user.up.sql', NOW());

COMMIT;
