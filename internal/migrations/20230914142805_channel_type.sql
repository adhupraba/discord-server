-- +goose Up
CREATE TYPE channel_type AS ENUM ('TEXT', 'AUDIO', 'VIDEO');

-- +goose Down
DROP TYPE channel_type;
