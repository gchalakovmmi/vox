-- TODO: RLS NOT ACTIVE!! App needs to connect as a non privileged user
BEGIN;

CREATE OR REPLACE FUNCTION current_app_user() RETURNS int AS $$
BEGIN
	RETURN NULLIF(current_setting('app.current_user_id', TRUE), '')::int;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE conversations (
	id		BIGSERIAL	PRIMARY KEY,
	user_id	INT	NOT NULL	REFERENCES users(id) ON DELETE CASCADE,
	title	TEXT	NOT NULL 	DEFAULT '',
	created_at	TIMESTAMPTZ	DEFAULT now(),
	updated_at	TIMESTAMPTZ	DEFAULT now()
);

CREATE TABLE conversation_messages (
	id		BIGSERIAL	PRIMARY KEY,
	conversation_id	BIGINT 		NOT NULL	REFERENCES conversations(id) ON DELETE CASCADE,
	role		VARCHAR(10)	CHECK (role IN ('user','assistant')),
	content		TEXT		NOT NULL,
	created_at	TIMESTAMPTZ	DEFAULT now()
);

CREATE INDEX idx_messages_conv_id ON conversation_messages(conversation_id, created_at);

-- Row Level Security
ALTER TABLE conversations ENABLE ROW LEVEL SECURITY;

CREATE POLICY conversations_owner_policy ON conversations
FOR ALL
USING (user_id = current_app_user());

COMMIT;
