-- +goose Up
-- +goose StatementBegin
CREATE TABLE direct_messages (
  id UUID PRIMARY KEY,
  content TEXT NOT NULL,
  file_url TEXT,
  member_id UUID NOT NULL REFERENCES members(id) ON DELETE CASCADE,
  conversation_id UUID NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
  deleted BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX direct_messages_member_id_idx ON direct_messages("member_id");

CREATE INDEX direct_messages_conversation_id_idx ON direct_messages("conversation_id");

CREATE FUNCTION update_updated_at_direct_messages()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = now();
  RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_direct_messages_updated_at
  BEFORE UPDATE
  ON
    direct_messages
  FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_direct_messages();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER update_direct_messages_updated_at ON direct_messages;

DROP FUNCTION update_updated_at_direct_messages();

DROP INDEX direct_messages_conversation_id_idx;

DROP INDEX direct_messages_member_id_idx;

DROP TABLE direct_messages;
-- +goose StatementEnd
