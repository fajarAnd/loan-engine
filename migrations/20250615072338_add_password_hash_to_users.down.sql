ALTER TABLE borrowers DROP COLUMN IF EXISTS password_hash;
ALTER TABLE investors DROP COLUMN IF EXISTS password_hash;