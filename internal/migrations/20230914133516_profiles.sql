-- +goose Up
-- +goose StatementBegin
CREATE TABLE profiles (
  id UUID PRIMARY KEY,
  user_id TEXT NOT NULL UNIQUE, -- clerk id
  name TEXT NOT NULL,
  image_url TEXT NOT NULL,
  email TEXT NOT NULL,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE FUNCTION update_updated_at_profiles()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_profiles_updated_at
  BEFORE UPDATE
  ON
    profiles
  FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_profiles();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER update_profiles_updated_at ON profiles;

DROP FUNCTION update_updated_at_profiles();

DROP TABLE profiles;
-- +goose StatementEnd
