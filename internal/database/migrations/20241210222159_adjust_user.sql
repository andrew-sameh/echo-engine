-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN first_name VARCHAR(255) NOT NULL;
ALTER TABLE users ADD COLUMN last_name VARCHAR(255) NOT NULL;

CREATE TYPE user_role AS ENUM ('admin', 'user', 'owner', 'manager', 'read_only');
ALTER TABLE users ADD COLUMN role user_role;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN role;
DROP TYPE user_role;

ALTER TABLE users DROP COLUMN last_name;
ALTER TABLE users DROP COLUMN first_name;
-- +goose StatementEnd