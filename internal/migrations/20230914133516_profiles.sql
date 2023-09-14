-- +goose Up
-- +goose StatementBegin
CREATE TABLE profiles (
  id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id TEXT NOT NULL UNIQUE,
  name TEXT NOT NULL,
  image_url TEXT NOT NULL,
  email TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE FUNCTION update_updated_at_profiles()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = now();
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
