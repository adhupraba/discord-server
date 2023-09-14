-- +goose Up
-- +goose StatementBegin
CREATE TABLE channels (
  id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL,
  type channel_type NOT NULL DEFAULT E'TEXT',
  profile_id UUID NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
  server_id UUID NOT NULL REFERENCES servers(id) ON DELETE CASCADE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX channels_profile_id_idx ON channels("profile_id");

CREATE INDEX channels_server_id_idx ON channels("server_id");

CREATE FUNCTION update_updated_at_channels()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = now();
  RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_channels_updated_at
  BEFORE UPDATE
  ON
    channels
  FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_channels();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER update_channels_updated_at ON channels;

DROP FUNCTION update_updated_at_channels();

DROP INDEX channels_server_id_idx;

DROP INDEX channels_profile_id_idx;

DROP TABLE channels;
-- +goose StatementEnd
