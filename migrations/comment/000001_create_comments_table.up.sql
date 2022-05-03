DROP TABLE IF EXISTS comments CASCADE;

create table comments
(
    id         BIGSERIAL PRIMARY KEY,
    uuid       VARCHAR(36)  NOT NULL UNIQUE,
    user_id    NUMERIC      NOT NULL,
    post_id    NUMERIC      NOT NULL,
    content    TEXT         NOT NULL,
    created_at     TIMESTAMP    NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMP    NOT NULL DEFAULT NOW(),
    deleted_at     TIMESTAMP    NULL
);

CREATE INDEX user_id_key ON comments (user_id);
CREATE INDEX post_id_key ON comments (post_id);

COMMENT ON COLUMN comments.content IS '內容';
