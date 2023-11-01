-- +goose Up
-- +goose StatementBegin
CREATE TABLE roles (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    "role" VARCHAR(40) UNIQUE NOT NULL CHECK ("role" IN ('admin','user'))
);

INSERT INTO roles (role)
VALUES ('admin');

INSERT INTO roles (role)
VALUES ('user');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE roles;
-- +goose StatementEnd
