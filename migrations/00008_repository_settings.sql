-- +goose Up
-- +goose StatementBegin
CREATE TABLE repository_settings (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    "snapshot_check_interval" INTERVAL NOT NULL,
    "backup_status" VARCHAR(20) DEFAULT 'unknown' NOT NULL,
    "recent_snapshot" TIMESTAMPTZ DEFAULT NULL,
    "repository_id" INT NOT NULL,
    CONSTRAINT fk_repositories
        FOREIGN KEY ("repository_id") REFERENCES repositories(id)
        ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE repository_settings;
-- +goose StatementEnd
