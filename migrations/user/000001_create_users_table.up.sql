DROP TABLE IF EXISTS users CASCADE;

create table users
(
    id         BIGSERIAL PRIMARY KEY,
    uuid       VARCHAR(36)  NOT NULL UNIQUE,
    username   VARCHAR(255) NOT NULL UNIQUE,
    email      VARCHAR(255) NOT NULL UNIQUE,
    avatar     VARCHAR(255) DEFAULT NULL,
    password   VARCHAR(255) NOT NULL,
    created_at TIMESTAMP    NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP    NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP    NULL
);

COMMENT
ON COLUMN users.username IS '使用者名稱';
COMMENT
ON COLUMN users.email IS 'Email';
COMMENT
ON COLUMN users.avatar IS '大頭貼';
COMMENT
ON COLUMN users.password IS '密碼';
