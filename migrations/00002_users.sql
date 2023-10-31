-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    username VARCHAR(40) UNIQUE NOT NULL,
    email VARCHAR(256),
    password_hash TEXT NOT NULL,
    role_id INT NOT NULL,
    CONSTRAINT fk_roles
        FOREIGN KEY (role_id) REFERENCES roles(id)
        ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
