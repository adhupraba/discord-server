-- +goose Up
-- +goose StatementBegin
CREATE TABLE members (
  id UUID PRIMARY KEY,
  role member_role NOT NULL DEFAULT E'GUEST',
  profile_id UUID NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
  server_id UUID NOT NULL REFERENCES servers(id) ON DELETE CASCADE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX members_profile_id_idx ON members("profile_id");

CREATE INDEX members_server_id_idx ON members("server_id");

CREATE FUNCTION update_updated_at_members()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = now();
  RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_members_updated_at
  BEFORE UPDATE
  ON
    members
  FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_members();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER update_members_updated_at ON members;

DROP FUNCTION update_updated_at_members();

DROP INDEX members_server_id_idx;

DROP INDEX members_profile_id_idx;

DROP TABLE members;
-- +goose StatementEnd
