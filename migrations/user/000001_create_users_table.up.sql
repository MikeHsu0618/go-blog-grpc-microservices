DROP TABLE IF EXISTS users CASCADE;

create table users
(
    id         BIGSERIAL PRIMARY KEY,
    uuid       VARCHAR(36)  NOT NULL UNIQUE,
    username   VARCHAR(255) NULL,
    email      VARCHAR(255) NULL,
    avatar     VARCHAR(255) DEFAULT NULL,
    password   VARCHAR(255) NOT NULL,
    created_at TIMESTAMP    NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP    NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP    NULL
);

CREATE UNIQUE INDEX idx_unq_users_email ON users (email) WHERE trim(email) != '';
CREATE UNIQUE INDEX idx_unq_users_username ON users (username) WHERE trim(username) != '';

COMMENT
ON COLUMN users.username IS '使用者名稱';
COMMENT
ON COLUMN users.email IS 'Email';
COMMENT
ON COLUMN users.avatar IS '大頭貼';
COMMENT
ON COLUMN users.password IS '密碼';
