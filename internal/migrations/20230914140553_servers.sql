-- +goose Up
-- +goose StatementBegin
CREATE TABLE servers (
  id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL,
  image_url TEXT NOT NULL,
  invite_code TEXT NOT NULL,
  profile_id UUID NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX servers_profile_id_idx ON servers("profile_id");

CREATE FUNCTION update_updated_at_servers()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = now();
  RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_servers_updated_at
  BEFORE UPDATE
  ON
    servers
  FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_servers();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER update_servers_updated_at ON servers;

DROP FUNCTION update_updated_at_servers();

DROP INDEX servers_profile_id_idx;

DROP TABLE servers;
-- +goose StatementEnd
