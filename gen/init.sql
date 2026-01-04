CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username TEXT UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL,
    email TEXT UNIQUE NOT NULL,
    salt TEXT NOT NULL,
    password TEXT,
    avatar_id TEXT,
    last_avatar_change TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS topics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT UNIQUE NOT NULL,
    description TEXT NOT NULL,
    lead_moderator_id UUID,
    logo_id TEXT,
    created_at TIMESTAMP NOT NULL,
    CONSTRAINT FK_topics_users FOREIGN KEY (lead_moderator_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS topic_moderators (
    user_id UUID NOT NULL,
    topic_id UUID NOT NULL,
    CONSTRAINT PK_topic_moderators PRIMARY KEY (user_id, topic_id),
    CONSTRAINT FK_topic_moderators_users FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT FK_topic_moderators_topics FOREIGN KEY (topic_id) REFERENCES topics(id)
);


CREATE TABLE IF NOT EXISTS posts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    upvotes INTEGER NOT NULL DEFAULT 0,
    downvotes INTEGER NOT NULL DEFAULT 0,
    deleted BOOLEAN NOT NULL DEFAULT FALSE,
    topic_id UUID NOT NULL,
    author_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL,
    edited_at TIMESTAMP NOT NULL,
    media_ids TEXT[] NOT NULL,
    CONSTRAINT FK_posts_users FOREIGN KEY (author_id) REFERENCES users(id),
    CONSTRAINT FK_posts_topics FOREIGN KEY (topic_id) REFERENCES topics(id)
);

CREATE TABLE IF NOT EXISTS comments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    description TEXT NOT NULL,
    upvotes INTEGER NOT NULL DEFAULT 0,
    downvotes INTEGER NOT NULL DEFAULT 0,
    deleted BOOLEAN NOT NULL DEFAULT FALSE,
    topic_id UUID NOT NULL,
    post_id UUID NOT NULL,
    replied_comment_id UUID,
    author_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL,
    edited_at TIMESTAMP NOT NULL,
    media_ids TEXT[] NOT NULL,
    CONSTRAINT FK_comments_users FOREIGN KEY (author_id) REFERENCES users(id),
    CONSTRAINT FK_comments_comments FOREIGN KEY (replied_comment_id) REFERENCES comments(id),
    CONSTRAINT FK_comments_posts FOREIGN KEY (post_id) REFERENCES posts(id),
    CONSTRAINT FK_comments_topics FOREIGN KEY (topic_id) REFERENCES topics(id)
);

CREATE TABLE IF NOT EXISTS banned_users (
    user_id UUID NOT NULL,
    topic_id UUID NOT NULL,
    CONSTRAINT PK_banned_users PRIMARY KEY (user_id, topic_id),
    CONSTRAINT FK_banned_users_users FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT FK_banned_users_topics FOREIGN KEY (topic_id) REFERENCES topics(id)
);

CREATE TABLE IF NOT EXISTS user_feeds (
    user_id UUID NOT NULL,
    topic_id UUID NOT NULL,
    CONSTRAINT PK_user_feeds PRIMARY KEY (user_id, topic_id),
    CONSTRAINT FK_user_feeds_users FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT FK_user_feeds_topics FOREIGN KEY (topic_id) REFERENCES topics(id)
);

CREATE TABLE IF NOT EXISTS post_votes (
    user_id UUID NOT NULL,
    post_id UUID NOT NULL,
    vote SMALLINT,
    CONSTRAINT PK_post_votes PRIMARY KEY (user_id, post_id),
    CONSTRAINT FK_post_votes_users FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT FK_post_votes_posts FOREIGN KEY (post_id) REFERENCES posts(id)
);

CREATE TABLE IF NOT EXISTS comment_votes (
    user_id UUID NOT NULL,
    comment_id UUID NOT NULL,
    vote SMALLINT,
    CONSTRAINT PK_comment_votes PRIMARY KEY (user_id, comment_id),
    CONSTRAINT FK_comment_votes_users FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT FK_comment_votes_posts FOREIGN KEY (comment_id) REFERENCES comments(id)
);

CREATE INDEX idx_posts_created_at ON posts(created_at);
CREATE INDEX idx_comments_created_at ON comments(created_at);




