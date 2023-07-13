-- +migrate Up
-- add test user with password = test
INSERT INTO users(email, passhash)
VALUES ('test@test.test', '$2a$10$wb7vaRhEHkI9fX7M./5G..E7u/XoqfM5lxMRqzT6TjAI2RgPiFh7u');

-- +migrate Down
DELETE
FROM users
WHERE email = 'test@test.test';
