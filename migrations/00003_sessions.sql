-- +goose Up
-- +goose StatementBegin
CREATE TABLE sessions (
id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    token_hash TEXT UNIQUE NOT NULL,
    user_id INT UNIQUE,
    CONSTRAINT fk_users
        FOREIGN KEY (user_id) REFERENCES users(id)
        ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE sessions;
-- +goose StatementEnd
