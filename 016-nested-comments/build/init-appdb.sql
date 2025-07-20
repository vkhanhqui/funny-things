CREATE TABLE IF NOT EXISTS comments (
    id UUID PRIMARY KEY,
    post_id UUID NOT NULL,
    user_id UUID NOT NULL,

    content TEXT NOT NULL,
    path VARCHAR(255),

    updated_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_comments_post_id ON comments(post_id);
CREATE INDEX IF NOT EXISTS idx_comments_user_id ON comments(user_id);