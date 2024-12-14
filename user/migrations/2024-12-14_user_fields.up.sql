BEGIN;

ALTER TABLE user_storage ADD column role VARCHAR(255) not NULL;

ALTER TABLE user_storage ADD column blocked BOOLEAN DEFAULT FALSE;

END;