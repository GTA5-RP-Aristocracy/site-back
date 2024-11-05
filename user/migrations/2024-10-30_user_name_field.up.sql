BEGIN;

ALTER TABLE user_storage ADD column name VARCHAR(255) not NULL;
ALTER TABLE user_storage ADD column created TIMESTAMP DEFAULT NOW();
ALTER TABLE user_storage ADD column updated TIMESTAMP DEFAULT NOW();

ALTER TABLE user_storage DROP column created_at;
ALTER TABLE user_storage DROP column updated_at;
END;
