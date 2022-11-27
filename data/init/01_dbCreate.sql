DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS tokens CASCADE;

DROP TRIGGER IF EXISTS updated_at_users ON users CASCADE;
DROP TRIGGER IF EXISTS updated_at_tokens ON tokens CASCADE;

CREATE TABLE users (
    id          UUID NOT NULL PRIMARY KEY,
    email       TEXT NOT NULL UNIQUE,
    encrypted_password TEXT NOT NULL,
    nickname    TEXT,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    logined_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE OR REPLACE FUNCTION updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
   IF ROW(new.*) IS DISTINCT FROM ROW(old.*) THEN
      new.updated_at = now(); 
      RETURN new;
   ELSE
      RETURN old;
   END IF;
END;
$$ language 'plpgsql';

CREATE TRIGGER updated_at_users
BEFORE UPDATE ON users FOR EACH 
ROW EXECUTE PROCEDURE  updated_at_column();

CREATE TABLE tokens (
    id          BIGSERIAL NOT NULL PRIMARY KEY,
    domain      TEXT NOT NULL UNIQUE,
    token       TEXT NOT NULL,
    --name        TEXT,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    requested_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expired_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER updated_at_tokens
BEFORE UPDATE ON tokens FOR EACH 
ROW EXECUTE PROCEDURE  updated_at_column();