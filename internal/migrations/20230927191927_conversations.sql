-- +goose Up
-- +goose StatementBegin
CREATE TABLE conversations (
  id UUID PRIMARY KEY,
  member_one_id UUID NOT NULL REFERENCES members(id) ON DELETE CASCADE,
  member_two_id UUID NOT NULL REFERENCES members(id) ON DELETE CASCADE,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,

  UNIQUE(member_one_id, member_two_id)
);

CREATE INDEX conversations_member_one_id_idx ON conversations("member_one_id");

CREATE INDEX conversations_member_two_id_idx ON conversations("member_two_id");

CREATE FUNCTION update_updated_at_conversations()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_conversations_updated_at
  BEFORE UPDATE
  ON
    conversations
  FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_conversations();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER update_conversations_updated_at ON conversations;

DROP FUNCTION update_updated_at_conversations();

DROP INDEX conversations_member_two_id_idx;

DROP INDEX conversations_member_one_id_idx;

DROP TABLE conversations;
-- +goose StatementEnd
