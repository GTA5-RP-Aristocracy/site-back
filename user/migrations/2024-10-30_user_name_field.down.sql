BEGIN;

ALTER TABLE user_storage DROP column name;
ALTER TABLE user_storage DROP column created;
ALTER TABLE user_storage DROP column updated;

ALTER TABLE user_storage ADD column created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE user_storage ADD column updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP;
END;