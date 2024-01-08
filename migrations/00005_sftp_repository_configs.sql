-- +goose Up
-- +goose StatementBegin
CREATE TABLE sftp_repository_configs (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    "user" VARCHAR(64) NOT NULL,
    "host" VARCHAR(255) NOT NULL,
    "pem_enc" TEXT NOT NULL,
    "repository_id" INT NOT NULL,
    CONSTRAINT fk_repositories
        FOREIGN KEY (repository_id) REFERENCES repositories(id)
        ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE sftp_repository_configs;
-- +goose StatementEnd
