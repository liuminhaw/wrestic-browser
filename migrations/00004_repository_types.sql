-- +goose Up
-- +goose StatementBegin
CREATE TABLE repository_types (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    "name" VARCHAR(40) UNIQUE NOT NULL
);

INSERT INTO repository_types ("name")
VALUES ('local'), ('s3'), ('sftp');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE repository_types;
-- +goose StatementEnd
