ALTER TABLE memberships
DROP COLUMN role_id;
ALTER TABLE memberships
ADD COLUMN role TEXT NOT NULL DEFAULT 'member';