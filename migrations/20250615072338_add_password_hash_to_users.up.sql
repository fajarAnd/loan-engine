ALTER TABLE borrowers ADD COLUMN password_hash VARCHAR(255);
ALTER TABLE investors ADD COLUMN password_hash VARCHAR(255);