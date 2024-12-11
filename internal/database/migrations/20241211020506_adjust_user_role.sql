-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ALTER COLUMN role SET NOT NULL;

ALTER TABLE users ADD COLUMN updated_at TIMESTAMP;
ALTER TABLE users ADD COLUMN deleted_at TIMESTAMP;

ALTER TABLE users ADD COLUMN is_superuser BOOLEAN NOT NULL DEFAULT false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN is_superuser;

ALTER TABLE users DROP COLUMN deleted_at;
ALTER TABLE users DROP COLUMN updated_at;

ALTER TABLE users ALTER COLUMN role DROP NOT NULL;
-- +goose StatementEnd