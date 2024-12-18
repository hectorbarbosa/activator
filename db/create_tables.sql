CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email varchar NOT NULL,
    user_name varchar NOT NULL,
    nick_name varchar NOT NULL,
    activated bool NOT NULL DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS tokens (
    hash bytea PRIMARY KEY NOT NULL,
    user_id int REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    expiry timestamp(0) with time zone NOT NULL
);

CREATE UNIQUE INDEX email_idx on public.users(email);