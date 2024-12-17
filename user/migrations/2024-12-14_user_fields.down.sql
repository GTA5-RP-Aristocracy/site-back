BEGIN;

ALTER TABLE user_storage DROP column role;

ALTER TABLE user_storage DROP column blocked;

END;