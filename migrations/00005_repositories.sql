-- +goose Up
-- +goose StatementBegin
CREATE TABLE repositories (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    "name" VARCHAR(255) UNIQUE NOT NULL,
    "destination" TEXT NOT NULL,
    "password_enc" VARCHAR(255) NOT NULL,
    "type_id" INT NOT NULL,
    CONSTRAINT fk_types
        FOREIGN KEY ("type_id") REFERENCES repository_types(id)
        ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE repositories;
-- +goose StatementEnd
