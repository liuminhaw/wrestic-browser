-- +goose Up
-- +goose StatementBegin
CREATE TABLE s3_repository_configs (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    "access_key_id_enc" VARCHAR(255) NOT NULL,
    "secret_access_key_enc" VARCHAR(255) NOT NULL,
    "region" VARCHAR(40) NOT NULL,
    "repository_id" INT NOT NULL,
    CONSTRAINT fk_repositories
        FOREIGN KEY (repository_id) REFERENCES repositories(id)
        ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE s3_repository_configs;
-- +goose StatementEnd
