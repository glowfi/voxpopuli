-- +goose Up

ALTER DATABASE voxpopuli SET TIMEZONE TO utc;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION fn_auto_update_updated_at_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd


CREATE TABLE topics (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE voxspheres (
    id UUID PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    public_description TEXT,
    topic_id UUID NOT NULL,
    community_icon TEXT,
    banner_background_image TEXT,
    banner_background_color VARCHAR(8),
    key_color VARCHAR(8),
    primary_color VARCHAR(8),
    over18 BOOLEAN NOT NULL DEFAULT false,
    spoilers_enabled BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at_unix BIGINT NOT NULL,
    updated_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_topic FOREIGN KEY(topic_id) REFERENCES topics(id) ON DELETE CASCADE ON UPDATE CASCADE
);


CREATE TABLE rules (
    id UUID PRIMARY KEY,
    voxsphere_id UUID NOT NULL,
    short_name TEXT NOT NULL,
    description TEXT NOT NULL,
    CONSTRAINT fk_voxsphere_id FOREIGN KEY(voxsphere_id) REFERENCES voxspheres(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE users (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    public_description TEXT,
    avatar_img TEXT,
    banner_img TEXT,
    iconcolor VARCHAR(8),
    keycolor VARCHAR(8),
    primarycolor VARCHAR(8),
    over18 BOOLEAN NOT NULL DEFAULT false,
    suspended BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at_unix BIGINT NOT NULL,
    updated_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE emojis (
    id UUID PRIMARY KEY,
    title VARCHAR(128) UNIQUE NOT NULL
);

CREATE TABLE custom_emojis(
    id UUID PRIMARY KEY,
    voxsphere_id UUID NOT NULL,
    url TEXT UNIQUE NOT NULL,
    title VARCHAR(128) NOT NULL,
    CONSTRAINT uk_voxsphere_id_link_text UNIQUE(voxsphere_id,title),
    CONSTRAINT fk_voxsphere_id FOREIGN KEY(voxsphere_id) REFERENCES voxspheres(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE user_flairs(
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    voxsphere_id UUID NOT NULL,
    full_text VARCHAR(255) NOT NULL,
    background_color VARCHAR(8),
    CONSTRAINT uk_user_id_voxsphere_id UNIQUE(user_id,voxsphere_id),
    CONSTRAINT fk_user_id FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_voxsphere_id FOREIGN KEY(voxsphere_id) REFERENCES voxspheres(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE user_flair_custom_emojis (
    custom_emoji_id UUID NOT NULL,
    user_flair_id UUID NOT NULL,
    order_index INTEGER NOT NULL,
    PRIMARY KEY(custom_emoji_id,user_flair_id,order_index),
    CONSTRAINT fk_custom_emoji_id FOREIGN KEY(custom_emoji_id) REFERENCES custom_emojis(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_user_flair_id FOREIGN KEY(user_flair_id) REFERENCES user_flairs(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE user_flair_emojis (
    emoji_id UUID NOT NULL,
    user_flair_id UUID NOT NULL,
    order_index INTEGER NOT NULL,
    PRIMARY KEY(emoji_id,user_flair_id,order_index),
    CONSTRAINT fk_emoji_id FOREIGN KEY(emoji_id) REFERENCES emojis(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_user_flair_id FOREIGN KEY(user_flair_id) REFERENCES user_flairs(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE user_flair_descriptions (
    user_flair_id UUID NOT NULL,
    order_index INTEGER NOT NULL,
    description VARCHAR(255) NOT NULL,
    PRIMARY KEY(user_flair_id,order_index),
    CONSTRAINT fk_user_flair_id FOREIGN KEY(user_flair_id) REFERENCES user_flairs(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE trophies (
    id UUID PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    image_link TEXT NOT NULL
);

CREATE TABLE user_trophies (
    user_id UUID NOT NULL,
    trophy_id UUID NOT NULL,
    PRIMARY KEY (user_id,trophy_id),
    CONSTRAINT fk_user_id FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT trophy_id FOREIGN KEY(trophy_id) REFERENCES trophies(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE voxsphere_members (
    voxsphere_id UUID NOT NULL,
    user_id UUID NOT NULL,
    PRIMARY KEY (voxsphere_id, user_id),
    CONSTRAINT fk_voxsphere_id FOREIGN KEY(voxsphere_id) REFERENCES voxspheres(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_user_id FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE voxsphere_moderators (
    voxsphere_id UUID NOT NULL,
    user_id UUID NOT NULL,
    PRIMARY KEY (voxsphere_id, user_id),
    CONSTRAINT fk_voxsphere_id FOREIGN KEY(voxsphere_id) REFERENCES voxspheres(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_user_id FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE posts (
    id UUID PRIMARY KEY,
    author_id UUID NOT NULL,
    voxsphere_id UUID NOT NULL,
    title VARCHAR(255) NOT NULL,
    text TEXT,
    text_html TEXT,
    ups INTEGER NOT NULL,
    over18 BOOLEAN NOT NULL DEFAULT false,
    spoiler BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at_unix BIGINT NOT NULL,
    updated_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_author_id FOREIGN KEY(author_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_voxsphere_id FOREIGN KEY(voxsphere_id) REFERENCES voxspheres(id) ON DELETE CASCADE ON UPDATE CASCADE
);

create type media_type as enum (
  'image',
  'gif',
  'video',
  'gallery',
  'link',
  'multi'
);


CREATE TABLE post_medias (
    id UUID PRIMARY KEY,
    post_id UUID NOT NULL,
    media_type media_type NOT NULL,
    CONSTRAINT fk_post_id FOREIGN KEY(post_id) REFERENCES posts(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE images (
    id UUID PRIMARY KEY,
    media_id UUID NOT NULL,
    CONSTRAINT fk_media_id FOREIGN KEY(media_id) REFERENCES post_medias(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE image_metadatas (
    id UUID PRIMARY KEY,
    image_id UUID NOT NULL,
    height INTEGER NOT NULL,
    width INTEGER NOT NULL,
    url TEXT UNIQUE NOT NULL,
    created_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at_unix BIGINT NOT NULL,
    updated_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_image_id FOREIGN KEY(image_id) REFERENCES images(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE gifs (
    id UUID PRIMARY KEY,
    media_id UUID NOT NULL,
    CONSTRAINT fk_media_id FOREIGN KEY(media_id) REFERENCES post_medias(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE gif_metadatas (
    id UUID PRIMARY KEY,
    gif_id UUID NOT NULL,
    height INTEGER NOT NULL,
    width INTEGER NOT NULL,
    url TEXT UNIQUE NOT NULL,
    created_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at_unix BIGINT NOT NULL,
    updated_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_gif_id FOREIGN KEY(gif_id) REFERENCES gifs(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE videos (
    id UUID PRIMARY KEY,
    media_id UUID NOT NULL,
    url TEXT UNIQUE NOT NULL,
    height INTEGER NOT NULL,
    width INTEGER NOT NULL,
    created_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at_unix BIGINT NOT NULL,
    updated_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_media_id FOREIGN KEY(media_id) REFERENCES post_medias(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE galleries (
    id UUID PRIMARY KEY,
    media_id UUID NOT NULL,
    CONSTRAINT fk_media_id FOREIGN KEY(media_id) REFERENCES post_medias(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE gallery_metadatas (
    id UUID PRIMARY KEY,
    gallery_id UUID NOT NULL,
    order_index INTEGER NOT NULL,
    height INTEGER NOT NULL,
    width INTEGER NOT NULL,
    url TEXT UNIQUE NOT NULL,
    created_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at_unix BIGINT NOT NULL,
    updated_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_gallery_id FOREIGN KEY(gallery_id) REFERENCES galleries(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE links (
    id UUID PRIMARY KEY,
    media_id UUID NOT NULL,
    link TEXT NOT NULL,
    created_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at_unix BIGINT NOT NULL,
    updated_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_media_id FOREIGN KEY(media_id) REFERENCES post_medias(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE comments (
    id UUID PRIMARY KEY,
    author_id UUID NOT NULL,
    parent_comment_id UUID,
    post_id UUID NOT NULL,
    body TEXT NOT NULL,
    body_html TEXT NOT NULL,
    ups INTEGER NOT NULL DEFAULT 0,
    score INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at_unix BIGINT NOT NULL,
    updated_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_author_id FOREIGN KEY(author_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_parent_comment_id FOREIGN KEY(parent_comment_id) REFERENCES comments(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_post_id FOREIGN KEY(post_id) REFERENCES posts(id) ON DELETE CASCADE ON UPDATE CASCADE
);


CREATE TABLE awards (
    id UUID PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    image_link TEXT NOT NULL
);

CREATE TABLE post_awards (
    post_id UUID NOT NULL,
    award_id UUID NOT NULL,
    PRIMARY KEY (post_id,award_id),
    CONSTRAINT fk_post_id FOREIGN KEY(post_id) REFERENCES posts(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_award_id FOREIGN KEY(award_id) REFERENCES awards(id) ON DELETE CASCADE ON UPDATE CASCADE
);


CREATE TABLE post_flairs(
    id UUID PRIMARY KEY,
    post_id UUID NOT NULL,
    voxsphere_id UUID NOT NULL,
    full_text VARCHAR(255) NOT NULL,
    background_color VARCHAR(8),
    CONSTRAINT post_id_voxsphere_id UNIQUE(post_id,voxsphere_id),
    CONSTRAINT fk_post_id FOREIGN KEY(post_id) REFERENCES posts(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_voxsphere_id FOREIGN KEY(voxsphere_id) REFERENCES voxspheres(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE post_flair_custom_emojis (
    custom_emoji_id UUID NOT NULL,
    post_flair_id UUID NOT NULL,
    order_index INTEGER NOT NULL,
    PRIMARY KEY(custom_emoji_id,post_flair_id,order_index),
    CONSTRAINT fk_custom_emoji_id FOREIGN KEY(custom_emoji_id) REFERENCES custom_emojis(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_post_flair_id FOREIGN KEY(post_flair_id) REFERENCES post_flairs(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE post_flair_emojis (
    emoji_id UUID NOT NULL,
    post_flair_id UUID NOT NULL,
    order_index INTEGER NOT NULL,
    PRIMARY KEY(emoji_id,post_flair_id,order_index),
    CONSTRAINT fk_emoji_id FOREIGN KEY(emoji_id) REFERENCES emojis(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_post_flair_id FOREIGN KEY(post_flair_id) REFERENCES post_flairs(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE post_flair_descriptions (
    post_flair_id UUID NOT NULL,
    order_index INTEGER NOT NULL,
    description VARCHAR(255) NOT NULL,
    PRIMARY KEY(post_flair_id,order_index),
    CONSTRAINT fk_post_flair_id FOREIGN KEY(post_flair_id) REFERENCES post_flairs(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON voxspheres
FOR EACH ROW
EXECUTE PROCEDURE fn_auto_update_updated_at_timestamp();

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON posts
FOR EACH ROW
EXECUTE PROCEDURE fn_auto_update_updated_at_timestamp();

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE PROCEDURE fn_auto_update_updated_at_timestamp();

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON comments
FOR EACH ROW
EXECUTE PROCEDURE fn_auto_update_updated_at_timestamp();

-- Create Unique Index
CREATE UNIQUE INDEX "topic_name" ON "topics"("name");
CREATE UNIQUE INDEX "trophy_title" ON "trophies"("title");
CREATE UNIQUE INDEX "awards_title" ON "awards"("title");
CREATE UNIQUE INDEX "voxsphere_title" ON "voxspheres"("title");
CREATE UNIQUE INDEX "user_name" ON "users"("name");

-- Create indexes for foreign key columns
CREATE INDEX idx_voxsphere_topic_id ON voxspheres (topic_id);
CREATE INDEX idx_rules_voxsphere_id ON rules (voxsphere_id);
CREATE INDEX idx_user_flairs_user_id ON user_flairs (user_id);
CREATE INDEX idx_user_flairs_voxsphere_id ON user_flairs (voxsphere_id);
CREATE INDEX idx_user_flair_emojis_user_flair_id ON user_flair_emojis (user_flair_id);
CREATE INDEX idx_user_flair_descriptions_user_flair_id ON user_flair_descriptions (user_flair_id);
CREATE INDEX idx_user_trophies_user_id ON user_trophies (user_id);
CREATE INDEX idx_user_trophies_trophy_id ON user_trophies (trophy_id);
CREATE INDEX idx_voxsphere_members_voxsphere_id ON voxsphere_members (voxsphere_id);
CREATE INDEX idx_voxsphere_members_user_id ON voxsphere_members (user_id);
CREATE INDEX idx_voxsphere_moderators_voxsphere_id ON voxsphere_moderators (voxsphere_id);
CREATE INDEX idx_voxsphere_moderators_user_id ON voxsphere_moderators (user_id);
CREATE INDEX idx_posts_author_id ON posts (author_id);
CREATE INDEX idx_posts_voxsphere_id ON posts (voxsphere_id);
CREATE INDEX idx_post_medias_post_id ON post_medias (post_id);
CREATE INDEX idx_images_media_id ON images (media_id);
CREATE INDEX idx_gifs_media_id ON gifs (media_id);
CREATE INDEX idx_videos_media_id ON videos (media_id);
CREATE INDEX idx_galleries_media_id ON galleries (media_id);
CREATE INDEX idx_links_media_id ON links (media_id);
CREATE INDEX idx_comments_author_id ON comments (author_id);
CREATE INDEX idx_comments_parent_comment_id ON comments (parent_comment_id);
CREATE INDEX idx_comments_post_id ON comments (post_id);
CREATE INDEX idx_post_awards_post_id ON post_awards (post_id);
CREATE INDEX idx_post_awards_award_id ON post_awards (award_id);
CREATE INDEX idx_post_flairs_post_id ON post_flairs (post_id);
CREATE INDEX idx_post_flairs_voxsphere_id ON post_flairs (voxsphere_id);
CREATE INDEX idx_post_flair_emojis_post_flair_id ON post_flair_emojis (post_flair_id);
CREATE INDEX idx_post_flair_descriptions_post_flair_id ON post_flair_descriptions (post_flair_id);

-- Create indexes for columns used in WHERE and JOIN clauses
CREATE INDEX idx_voxspheres_title ON voxspheres (title);
CREATE INDEX idx_voxspheres_public_description ON voxspheres (public_description);
CREATE INDEX idx_posts_title ON posts (title);
CREATE INDEX idx_posts_text ON posts (text);
CREATE INDEX idx_comments_body ON comments (body);
CREATE INDEX idx_users_name ON users (name);
CREATE INDEX idx_users_public_description ON users (public_description);

-- Create indexes for columns used in ORDER BY and LIMIT clauses
CREATE INDEX idx_posts_created_at ON posts (created_at);
CREATE INDEX idx_comments_created_at ON comments (created_at);
CREATE INDEX idx_voxspheres_created_at ON voxspheres (created_at);
CREATE INDEX idx_users_created_at ON users (created_at);
CREATE INDEX idx_image_metadatas_created_at ON image_metadatas (created_at);
CREATE INDEX idx_gif_metadatas_created_at ON gif_metadatas (created_at);
CREATE INDEX idx_gallery_metadatas_created_at ON gallery_metadatas (created_at);
CREATE INDEX idx_videos_created_at ON videos (created_at);
CREATE INDEX idx_link_created_at ON links (created_at);

-- +goose Down
DROP EXTENSION "uuid-ossp";

DROP INDEX idx_voxspheres_created_at;
DROP INDEX idx_users_created_at;
DROP INDEX idx_comments_created_at;
DROP INDEX idx_posts_created_at;
DROP INDEX idx_image_metadatas_created_at;
DROP INDEX idx_gif_metadatas_created_at;
DROP INDEX idx_gallery_metadatas_created_at;
DROP INDEX idx_videos_created_at;
DROP INDEX idx_link_created_at;
DROP INDEX idx_users_public_description;
DROP INDEX idx_users_name;
DROP INDEX idx_comments_body;
DROP INDEX idx_posts_text;
DROP INDEX idx_posts_title;
DROP INDEX idx_voxspheres_public_description;
DROP INDEX idx_voxspheres_title;
DROP INDEX idx_post_flair_descriptions_post_flair_id;
DROP INDEX idx_post_flair_emojis_post_flair_id;
DROP INDEX idx_post_flairs_voxsphere_id;
DROP INDEX idx_post_flairs_post_id;
DROP INDEX idx_post_awards_award_id;
DROP INDEX idx_post_awards_post_id;
DROP INDEX idx_comments_post_id;
DROP INDEX idx_comments_parent_comment_id;
DROP INDEX idx_comments_author_id;
DROP INDEX idx_links_media_id;
DROP INDEX idx_galleries_media_id;
DROP INDEX idx_videos_media_id;
DROP INDEX idx_gifs_media_id;
DROP INDEX idx_images_media_id;
DROP INDEX idx_post_medias_post_id;
DROP INDEX idx_posts_voxsphere_id;
DROP INDEX idx_posts_author_id;
DROP INDEX idx_voxsphere_moderators_user_id;
DROP INDEX idx_voxsphere_moderators_voxsphere_id;
DROP INDEX idx_voxsphere_members_user_id;
DROP INDEX idx_voxsphere_members_voxsphere_id;
DROP INDEX idx_user_trophies_trophy_id;
DROP INDEX idx_user_trophies_user_id;
DROP INDEX idx_user_flair_descriptions_user_flair_id;
DROP INDEX idx_user_flair_emojis_user_flair_id;
DROP INDEX idx_user_flairs_voxsphere_id;
DROP INDEX idx_user_flairs_user_id;
DROP INDEX idx_rules_voxsphere_id;
DROP INDEX idx_voxsphere_topic_id;

DROP INDEX "user_name";
DROP INDEX "voxsphere_title";
DROP INDEX "awards_title";
DROP INDEX "trophy_title";
DROP INDEX "topic_name";

DROP TRIGGER set_timestamp ON comments;
DROP TRIGGER set_timestamp ON posts;
DROP TRIGGER set_timestamp ON users;
DROP TRIGGER set_timestamp ON voxspheres;

DROP TABLE post_flairs CASCADE;
DROP TABLE post_flair_emojis CASCADE;
DROP TABLE post_flair_custom_emojis CASCADE;
DROP TABLE post_flair_descriptions CASCADE;
DROP TABLE post_awards CASCADE;
DROP TABLE awards CASCADE;
DROP TABLE comments CASCADE;
DROP TABLE links CASCADE;
DROP TABLE galleries CASCADE;
DROP TABLE gallery_metadatas CASCADE;
DROP TABLE videos CASCADE;
DROP TABLE gifs CASCADE;
DROP TABLE gif_metadatas CASCADE;
DROP TABLE image_metadatas CASCADE;
DROP TABLE images CASCADE;
DROP TABLE post_medias CASCADE;
DROP TABLE posts CASCADE;
DROP TABLE voxsphere_moderators CASCADE;
DROP TABLE voxsphere_members CASCADE;
DROP TABLE user_trophies CASCADE;
DROP TABLE trophies CASCADE;
DROP TABLE rules CASCADE;
DROP TABLE users CASCADE;
DROP TABLE voxspheres CASCADE;
DROP TABLE topics CASCADE;
DROP TABLE custom_emojis CASCADE;
DROP TABLE emojis CASCADE;
DROP TABLE user_flairs CASCADE;
DROP TABLE user_flair_emojis CASCADE;
DROP TABLE user_flair_custom_emojis CASCADE;
DROP TABLE user_flair_descriptions CASCADE;

DROP TYPE media_type;

DROP FUNCTION fn_auto_update_updated_at_timestamp();
