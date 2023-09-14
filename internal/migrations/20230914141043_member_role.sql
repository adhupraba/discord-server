-- +goose Up
CREATE TYPE member_role AS ENUM ('ADMIN', 'MODERATOR', 'GUEST');

-- +goose Down
DROP TYPE member_role;
