-- +goose Up
-- +goose StatementBegin
CREATE TABLE news
(
    id         INT PRIMARY KEY,
    title      TEXT      NOT NULL,
    summary    TEXT      NOT NULL,
    image      TEXT      NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX created_at_idx ON news (created_at);

CREATE TABLE users_have_news_with_reactions
(
    user_id    TEXT REFERENCES users (email) NOT NULL,
    new_id     INT REFERENCES news (id)      NOT NULL,
    is_like    boolean                       NOT NULL,
    viewed_at  TIMESTAMP                     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP                     NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX user_id_new_id_url_unique_idx ON users_have_news_with_reactions (user_id, new_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS created_at_idx;
DROP TABLE IF EXISTS news;
-- +goose StatementEnd
