-- +goose Up
-- +goose StatementBegin
CREATE TABLE messages (
  id UUID PRIMARY KEY,
  content TEXT NOT NULL,
  file_url TEXT,
  member_id UUID NOT NULL REFERENCES members(id) ON DELETE CASCADE,
  channel_id UUID NOT NULL REFERENCES channels(id) ON DELETE CASCADE,
  deleted BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX messages_member_id_idx ON messages("member_id");

CREATE INDEX messages_channel_id_idx ON messages("channel_id");

CREATE FUNCTION update_updated_at_messages()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_messages_updated_at
  BEFORE UPDATE
  ON
    messages
  FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_messages();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER update_messages_updated_at ON messages;

DROP FUNCTION update_updated_at_messages();

DROP INDEX messages_channel_id_idx;

DROP INDEX messages_member_id_idx;

DROP TABLE messages;
-- +goose StatementEnd
