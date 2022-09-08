-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.roles(
    id   SMALLSERIAL PRIMARY KEY,
    name VARCHAR(255)
);
INSERT INTO roles (id, name) VALUES (1,'Admin'),(2,'User');
CREATE TABLE IF NOT EXISTS public.users (
    id        SERIAL PRIMARY KEY,
    email     VARCHAR(255) NOT NULL UNIQUE,
    full_name VARCHAR(255) NOT NULL,
    role      SMALLINT REFERENCES roles (id) NOT NULL,
    password  VARCHAR(255) NOT NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.users;
DROP TABLE IF EXISTS public.roles;

-- +goose StatementEnd
