-- +goose Up
CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE voxsphere_rules (
    id SERIAL PRIMARY KEY,
    short_name TEXT NOT NULL,
    description TEXT NOT NULL
);

CREATE TABLE topics (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE flair_emojis (
    id SERIAL PRIMARY KEY,
    text VARCHAR(255) NOT NULL,
    url TEXT NOT NULL
);

CREATE TABLE flair_texts (
    id SERIAL PRIMARY KEY,
    text VARCHAR(255) NOT NULL
);

CREATE TABLE post_flairs(
    id SERIAL PRIMARY KEY,
    full_text VARCHAR(255),
    background_color VARCHAR(7)  -- hex color value
);

CREATE TABLE user_flairs(
    id SERIAL PRIMARY KEY,
    full_text VARCHAR(255),
    background_color VARCHAR(7)  -- hex color value
);

CREATE TABLE post_flairs_rich_text (
    id SERIAL,
    flair_id INTEGER NOT NULL,
    flair_emoji_id INTEGER,
    flair_text_id INTEGER,
    PRIMARY KEY (id,flair_id, flair_emoji_id,flair_text_id),
    FOREIGN KEY (flair_id) REFERENCES post_flairs(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (flair_emoji_id) REFERENCES flair_emojis(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (flair_text_id) REFERENCES flair_texts(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CHECK((flair_emoji_id IS NULL AND flair_text_id IS NOT NULL) OR (flair_emoji_id IS NOT NULL AND flair_text_id IS NULL))
);

CREATE TABLE user_flairs_rich_text (
    id SERIAL,
    flair_id INTEGER NOT NULL,
    flair_emoji_id INTEGER,
    flair_text_id INTEGER,
    PRIMARY KEY (id,flair_id, flair_emoji_id,flair_text_id),
    FOREIGN KEY (flair_id) REFERENCES user_flairs(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (flair_emoji_id) REFERENCES flair_emojis(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (flair_text_id) REFERENCES flair_texts(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CHECK((flair_emoji_id IS NULL AND flair_text_id IS NOT NULL) OR (flair_emoji_id IS NOT NULL AND flair_text_id IS NULL))
);


CREATE TABLE voxspheres (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    public_description TEXT,
    community_icon TEXT,
    banner_background_image TEXT,
    key_color VARCHAR(7),  -- hex color value
    primary_color VARCHAR(7),  -- hex color value
    banner_background_color VARCHAR(7),  -- hex color value
    updated_at TIMESTAMP(6) NOT NULL,
    created_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at_unix INTEGER NOT NULL,
    created_at_human TEXT NOT NULL,
    subscribers INTEGER NOT NULL,
    subscribers_human TEXT NOT NULL,
    over18 BOOLEAN NOT NULL DEFAULT false,
    spoilers_enabled BOOLEAN NOT NULL DEFAULT false
);

CREATE TABLE trophies (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    image_link TEXT NOT NULL
);

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    updated_at TIMESTAMP(6) NOT NULL,
    created_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    cake_day_unix INTEGER NOT NULL,
    cake_day_human TEXT NOT NULL,
    age TEXT NOT NULL, -- 17yrs || 6months || 1day
    avatar_img TEXT,
    banner_img TEXT,
    public_description TEXT,
    keycolor VARCHAR(7), -- hex color value
    primarycolor VARCHAR(7), -- hex color value
    iconcolor VARCHAR(7), -- hex color value
    over18 BOOLEAN NOT NULL DEFAULT false,
    suspended BOOLEAN NOT NULL DEFAULT false
);

-- Create many-to-many relationships between users and trophies
CREATE TABLE user_trophies_mapping (
    user_id INTEGER NOT NULL,
    trophy_id INTEGER NOT NULL,
    PRIMARY KEY (user_id, trophy_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (trophy_id) REFERENCES trophies(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE awards (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    image_link TEXT NOT NULL
);

CREATE TABLE comments (
    id SERIAL PRIMARY KEY,
    author_id INTEGER NOT NULL,
    body TEXT NOT NULL,
    body_html TEXT NOT NULL,
    ups INTEGER NOT NULL DEFAULT 0,
    score INTEGER NOT NULL DEFAULT 0,
    parent_comment_id INTEGER,
    FOREIGN KEY (author_id) REFERENCES comments(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (parent_comment_id) REFERENCES users(id)
);

CREATE TABLE media_content (
    id SERIAL PRIMARY KEY
);

CREATE TABLE galleries (
    id SERIAL PRIMARY KEY,
    media_type TEXT NOT NULL,
    media_content_id INTEGER,
    FOREIGN KEY (media_content_id) REFERENCES media_content(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE gallery_image_resolutions (
    id SERIAL PRIMARY KEY,
    height INTEGER NOT NULL,
    width INTEGER NOT NULL,
    url TEXT NOT NULL,
    gallery_id INTEGER NOT NULL,
    FOREIGN KEY (gallery_id) REFERENCES galleries(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE links (
    id SERIAL PRIMARY KEY,
    link TEXT NOT NULL,
    media_type TEXT NOT NULL,
    media_content_id INTEGER,
    FOREIGN KEY (media_content_id) REFERENCES media_content(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE images (
    id SERIAL PRIMARY KEY,
    media_type TEXT NOT NULL,
    media_content_id INTEGER,
    FOREIGN KEY (media_content_id) REFERENCES media_content(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE image_resolutions (
    id SERIAL PRIMARY KEY,
    height INTEGER NOT NULL,
    width INTEGER NOT NULL,
    url TEXT NOT NULL,
    image_id INTEGER NOT NULL,
    FOREIGN KEY (image_id) REFERENCES images(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE videos (
    id SERIAL PRIMARY KEY,
    url TEXT NOT NULL,
    height INTEGER NOT NULL,
    width INTEGER NOT NULL,
    media_type TEXT NOT NULL,
    media_content_id INTEGER,
    FOREIGN KEY (media_content_id) REFERENCES media_content(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE gifs (
    id SERIAL PRIMARY KEY,
    media_type TEXT NOT NULL,
    media_content_id INTEGER,
    FOREIGN KEY (media_content_id) REFERENCES media_content(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE gif_image_resolutions (
    id SERIAL PRIMARY KEY,
    height INTEGER NOT NULL,
    width INTEGER NOT NULL,
    url TEXT NOT NULL,
    gif_id INTEGER NOT NULL,
    FOREIGN KEY (gif_id) REFERENCES gifs(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    updated_at TIMESTAMP(6) NOT NULL,
    created_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at_unix INTEGER NOT NULL,
    created_at_human TEXT NOT NULL,
    num_comments INTEGER NOT NULL,
    ups INTEGER NOT NULL,
    text TEXT,
    text_html TEXT,
    over18 BOOLEAN NOT NULL DEFAULT false,
    spoiler BOOLEAN NOT NULL DEFAULT false
);

-- Create many-to-many relationships between posts and media content
CREATE TABLE posts_media_content_mapping (
    post_id INTEGER NOT NULL,
    media_content_id INTEGER NOT NULL,
    PRIMARY KEY (post_id, media_content_id),
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (media_content_id) REFERENCES media_content(id) ON DELETE CASCADE ON UPDATE CASCADE
);

-- Create many-to-many relationships between posts and awards
CREATE TABLE post_awards_mapping (
    post_id INTEGER NOT NULL,
    award_id INTEGER NOT NULL,
    PRIMARY KEY (post_id, award_id),
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (award_id) REFERENCES awards(id) ON DELETE CASCADE ON UPDATE CASCADE
);

-- Create many-to-many relationships between posts and comments
CREATE TABLE posts_comments (
    post_id INTEGER NOT NULL,
    comment_id INTEGER NOT NULL,
    PRIMARY KEY (post_id, comment_id),
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (comment_id) REFERENCES comments(id) ON DELETE CASCADE ON UPDATE CASCADE
);

-- Create many-to-many relationships between posts and users
CREATE TABLE post_users_mapping (
    post_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    PRIMARY KEY (post_id, user_id),
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE
);

-- Create many-to-many relationships between posts and voxspheres
CREATE TABLE post_voxspheres_mapping (
    post_id INTEGER NOT NULL,
    voxsphere_id INTEGER NOT NULL,
    PRIMARY KEY (post_id, voxsphere_id),
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (voxsphere_id) REFERENCES voxspheres(id) ON DELETE CASCADE ON UPDATE CASCADE
);


-- Create many-to-many relationships between posts and voxspheres
CREATE TABLE voxsphere_rules_mapping (
    voxsphere_id INTEGER NOT NULL,
    rule_id INTEGER NOT NULL,
    PRIMARY KEY (voxsphere_id, rule_id),
    FOREIGN KEY (voxsphere_id) REFERENCES voxspheres(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (rule_id) REFERENCES voxsphere_rules(id) ON DELETE CASCADE ON UPDATE CASCADE
);

-- Create many-to-many relationships between posts and voxspheres
CREATE TABLE topics_mapping (
    voxsphere_id INTEGER NOT NULL,
    topic_id INTEGER NOT NULL,
    PRIMARY KEY (voxsphere_id, topic_id),
    FOREIGN KEY (voxsphere_id) REFERENCES voxspheres(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (topic_id) REFERENCES topics(id) ON DELETE CASCADE ON UPDATE CASCADE
);

-- Create many-to-many relationships between voxsphere,post,post_flair
CREATE TABLE voxsphere_post_post_flairs_mapping (
    voxsphere_id INTEGER NOT NULL,
    post_flair_id INTEGER NOT NULL,
    post_id INTEGER NOT NULL,
    PRIMARY KEY (voxsphere_id, post_flair_id,post_id),
    FOREIGN KEY (voxsphere_id) REFERENCES voxspheres(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (post_flair_id) REFERENCES post_flairs(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE ON UPDATE CASCADE
);

-- Create many-to-many relationships between voxsphere,user,user_flair
CREATE TABLE voxsphere_user_user_flairs_mapping (
    voxsphere_id INTEGER NOT NULL,
    user_flair_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    PRIMARY KEY (voxsphere_id, user_flair_id,user_id),
    FOREIGN KEY (voxsphere_id) REFERENCES voxspheres(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (user_flair_id) REFERENCES user_flairs(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON voxspheres
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON posts
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

-- Create Unique Index
CREATE UNIQUE INDEX "topics_name" ON "topics"("name");
CREATE UNIQUE INDEX "trophies_title" ON "trophies"("title");
CREATE UNIQUE INDEX "awards_title" ON "awards"("title");
CREATE UNIQUE INDEX "voxsphere_title" ON "awards"("title");
CREATE UNIQUE INDEX "user_name" ON "users"("name");

-- Indexes for frequently used columns in WHERE and JOIN clauses
CREATE INDEX idx_voxspheres_title ON voxspheres (title);
CREATE INDEX idx_voxspheres_updated_at ON voxspheres (updated_at);
CREATE INDEX idx_voxspheres_created_at ON voxspheres (created_at);

CREATE INDEX idx_posts_title ON posts (title);
CREATE INDEX idx_posts_updated_at ON posts (updated_at);
CREATE INDEX idx_posts_created_at ON posts (created_at);

CREATE INDEX idx_users_name ON users (name);
CREATE INDEX idx_users_updated_at ON users (updated_at);
CREATE INDEX idx_users_created_at ON users (created_at);

-- Indexes for foreign key columns
CREATE INDEX idx_posts_media_content_mapping_post_id ON posts_media_content_mapping (post_id);
CREATE INDEX idx_posts_media_content_mapping_media_content_id ON posts_media_content_mapping (media_content_id);

CREATE INDEX idx_post_awards_mapping_post_id ON post_awards_mapping (post_id);
CREATE INDEX idx_post_awards_mapping_award_id ON post_awards_mapping (award_id);

CREATE INDEX idx_posts_comments_post_id ON posts_comments (post_id);
CREATE INDEX idx_posts_comments_comment_id ON posts_comments (comment_id);

CREATE INDEX idx_post_users_mapping_post_id ON post_users_mapping (post_id);
CREATE INDEX idx_post_users_mapping_user_id ON post_users_mapping (user_id);

CREATE INDEX idx_post_voxspheres_mapping_post_id ON post_voxspheres_mapping (post_id);
CREATE INDEX idx_post_voxspheres_mapping_voxsphere_id ON post_voxspheres_mapping (voxsphere_id);

CREATE INDEX idx_voxsphere_rules_mapping_voxsphere_id ON voxsphere_rules_mapping (voxsphere_id);
CREATE INDEX idx_voxsphere_rules_mapping_rule_id ON voxsphere_rules_mapping (rule_id);

CREATE INDEX idx_topics_mapping_voxsphere_id ON topics_mapping (voxsphere_id);
CREATE INDEX idx_topics_mapping_topic_id ON topics_mapping (topic_id);

CREATE INDEX idx_voxsphere_post_post_flairs_mapping_voxsphere_id ON voxsphere_post_post_flairs_mapping (voxsphere_id);
CREATE INDEX idx_voxsphere_post_post_flairs_mapping_post_flair_id ON voxsphere_post_post_flairs_mapping (post_flair_id);
CREATE INDEX idx_voxsphere_post_post_flairs_mapping_post_id ON voxsphere_post_post_flairs_mapping (post_id);

CREATE INDEX idx_voxsphere_user_user_flairs_mapping_voxsphere_id ON voxsphere_user_user_flairs_mapping (voxsphere_id);
CREATE INDEX idx_voxsphere_user_user_flairs_mapping_user_flair_id ON voxsphere_user_user_flairs_mapping (user_flair_id);
CREATE INDEX idx_voxsphere_user_user_flairs_mapping_user_id ON voxsphere_user_user_flairs_mapping (user_id);

-- Indexes for columns used in ORDER BY and LIMIT clauses
CREATE INDEX idx_posts_num_comments ON posts (num_comments);
CREATE INDEX idx_posts_ups ON posts (ups);


-- +goose Down

-- Drop indexes
DROP INDEX idx_voxsphere_user_user_flairs_mapping_user_id;
DROP INDEX idx_voxsphere_user_user_flairs_mapping_user_flair_id;
DROP INDEX idx_voxsphere_user_user_flairs_mapping_voxsphere_id;
DROP INDEX idx_voxsphere_post_post_flairs_mapping_post_id;
DROP INDEX idx_voxsphere_post_post_flairs_mapping_post_flair_id;
DROP INDEX idx_voxsphere_post_post_flairs_mapping_voxsphere_id;
DROP INDEX idx_topics_mapping_topic_id;
DROP INDEX idx_topics_mapping_voxsphere_id;
DROP INDEX idx_voxsphere_rules_mapping_rule_id;
DROP INDEX idx_voxsphere_rules_mapping_voxsphere_id;
DROP INDEX idx_post_voxspheres_mapping_voxsphere_id;
DROP INDEX idx_post_voxspheres_mapping_post_id;
DROP INDEX idx_post_users_mapping_user_id;
DROP INDEX idx_post_users_mapping_post_id;
DROP INDEX idx_posts_comments_comment_id;
DROP INDEX idx_posts_comments_post_id;
DROP INDEX idx_post_awards_mapping_award_id;
DROP INDEX idx_post_awards_mapping_post_id;
DROP INDEX idx_posts_media_content_mapping_media_content_id;
DROP INDEX idx_posts_media_content_mapping_post_id;
DROP INDEX idx_users_created_at;
DROP INDEX idx_users_updated_at;
DROP INDEX idx_users_name;
DROP INDEX idx_posts_created_at;
DROP INDEX idx_posts_updated_at;
DROP INDEX idx_posts_title;
DROP INDEX idx_voxspheres_created_at;
DROP INDEX idx_voxspheres_updated_at;
DROP INDEX idx_voxspheres_title;

-- Drop unique indexes
DROP INDEX "user_name";
DROP INDEX "voxsphere_title";
DROP INDEX "awards_title";
DROP INDEX "trophies_title";
DROP INDEX "topics_name";

-- Drop tables
DROP TABLE voxsphere_user_user_flairs_mapping;
DROP TABLE voxsphere_post_post_flairs_mapping;
DROP TABLE topics_mapping;
DROP TABLE voxsphere_rules_mapping;
DROP TABLE post_voxspheres_mapping;
DROP TABLE post_users_mapping;
DROP TABLE posts_comments;
DROP TABLE post_awards_mapping;
DROP TABLE posts_media_content_mapping;
DROP TABLE gif_image_resolutions;
DROP TABLE gifs;
DROP TABLE videos;
DROP TABLE image_resolutions;
DROP TABLE images;
DROP TABLE links;
DROP TABLE galleries;
DROP TABLE gallery_image_resolutions;
DROP TABLE comments;
DROP TABLE awards;
DROP TABLE user_trophies_mapping;
DROP TABLE trophies;
DROP TABLE users;
DROP TABLE posts;
DROP TABLE media_content;
DROP TABLE voxspheres;
DROP TABLE user_flairs_rich_text;
DROP TABLE post_flairs_rich_text;
DROP TABLE user_flairs;
DROP TABLE post_flairs;
DROP TABLE flair_texts;
DROP TABLE flair_emojis;
DROP TABLE topics;
DROP TABLE voxsphere_rules;

-- Drop triggers
DROP TRIGGER set_timestamp ON users;
DROP TRIGGER set_timestamp ON posts;
DROP TRIGGER set_timestamp ON voxspheres;

-- Drop function
DROP FUNCTION trigger_set_timestamp;
