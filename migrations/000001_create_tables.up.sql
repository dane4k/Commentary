CREATE TABLE users
(
    id       SERIAL PRIMARY KEY,
    username VARCHAR(20) UNIQUE NOT NULL
);

CREATE TABLE posts
(
    id          SERIAL PRIMARY KEY,
    author_id   INTEGER       NOT NULL REFERENCES users (id),
    title       VARCHAR(150)  NOT NULL,
    content     VARCHAR(5000) NOT NULL,
    created     TIMESTAMP     NOT NULL,
    commentable BOOLEAN       NOT NULL
);

CREATE TABLE comments
(
    id        SERIAL PRIMARY KEY,
    post_id   INTEGER       NOT NULL REFERENCES posts (id),
    author_id INTEGER       NOT NULL REFERENCES users (id),
    content   VARCHAR(2000) NOT NULL,
    created   TIMESTAMP     NOT NULL,
    parent_id INTEGER
);

CREATE INDEX idx_comments_post_id ON comments (post_id);
CREATE INDEX idx_comments_parent_id ON comments (parent_id);
