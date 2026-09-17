ALTER TABLE users ADD COLUMN avatar_data bytea;
ALTER TABLE users ADD COLUMN avatar_type text NOT NULL DEFAULT '';
