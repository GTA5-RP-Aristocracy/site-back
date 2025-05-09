BEGIN;

ALTER TABLE user_storage ADD column role VARCHAR(255) DEFAULT 'user';

ALTER TABLE user_storage ADD column blocked BOOLEAN DEFAULT FALSE;

END;