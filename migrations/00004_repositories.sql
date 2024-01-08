-- +goose Up
-- +goose StatementBegin
CREATE TABLE repositories (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    "name" VARCHAR(255) UNIQUE NOT NULL,
    "type" VARCHAR(40) UNIQUE NOT NULL CHECK ("type" IN ('local','sftp','s3')),
    "destination" TEXT NOT NULL,
    "password_enc" VARCHAR(255) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE repositories;
-- +goose StatementEnd
