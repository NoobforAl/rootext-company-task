CREATE TABLE IF NOT EXISTS user_info (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(255) UNIQUE NOT NULL,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'user' CHECK (role IN ('user', 'admin')),
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS post (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES user_info (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS tag (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS post_tag (
    post_id INTEGER NOT NULL,
    tag_id INTEGER NOT NULL,
    PRIMARY KEY (post_id, tag_id),
    FOREIGN KEY (post_id) REFERENCES post (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS rating (
    id SERIAL PRIMARY KEY,
    post_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    rating_value INTEGER NOT NULL CHECK (rating_value IN (-1, 1)),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (post_id) REFERENCES post (id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES user_info (id) ON DELETE CASCADE,
    UNIQUE (post_id, user_id)
);

CREATE TABLE IF NOT EXISTS post_rating_summary (
    post_id INTEGER PRIMARY KEY,
    total_rating INTEGER DEFAULT 0,
    FOREIGN KEY (post_id) REFERENCES post (id) ON DELETE CASCADE
);


-- Install pgcrypto extension for hash functions
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-----------------------------------
------------ TRIGGERS -------------
-----------------------------------

-- Create a function to update the updated_at column
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to update post_rating_summary when a new rating is added or updated
CREATE OR REPLACE FUNCTION update_post_rating_summary() RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE post_rating_summary
        SET total_rating = total_rating + NEW.rating_value
        WHERE post_id = NEW.post_id;
    ELSIF TG_OP = 'UPDATE' THEN
        UPDATE post_rating_summary
        SET total_rating = total_rating - OLD.rating_value + NEW.rating_value
        WHERE post_id = NEW.post_id;
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE post_rating_summary
        SET total_rating = total_rating - OLD.rating_value
        WHERE post_id = OLD.post_id;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create post_rating_summary trigger if not exists
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_trigger
        WHERE tgname = 'update_user_info_updated_at'
    ) THEN
        CREATE TRIGGER update_user_info_updated_at
        BEFORE UPDATE ON user_info
        FOR EACH ROW
        EXECUTE FUNCTION update_updated_at_column();
    END IF;
END
$$;

-- Initialize post_rating_summary for existing posts
INSERT INTO
    post_rating_summary (post_id)
SELECT id
FROM post ON CONFLICT (post_id) DO NOTHING;

-- Create a function to update the updated_at column
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create triggers for each table that needs the updated_at column to be auto-updated
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_trigger
        WHERE tgname = 'update_user_info_updated_at'
    ) THEN
        CREATE TRIGGER update_user_info_updated_at
        BEFORE UPDATE ON user_info
        FOR EACH ROW
        EXECUTE FUNCTION update_updated_at_column();
    END IF;
END
$$;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_trigger
        WHERE tgname = 'update_post_updated_at'
    ) THEN
        CREATE TRIGGER update_post_updated_at
        BEFORE UPDATE ON post
        FOR EACH ROW
        EXECUTE FUNCTION update_updated_at_column();
    END IF;
END
$$;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_trigger
        WHERE tgname = 'update_tag_updated_at'
    ) THEN
        CREATE TRIGGER update_tag_updated_at
        BEFORE UPDATE ON tag
        FOR EACH ROW
        EXECUTE FUNCTION update_updated_at_column();
    END IF;
END
$$;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_trigger
        WHERE tgname = 'update_post_tag_updated_at'
    ) THEN
        CREATE TRIGGER update_post_tag_updated_at
        BEFORE UPDATE ON post_tag
        FOR EACH ROW
        EXECUTE FUNCTION update_updated_at_column();
    END IF;
END
$$;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_trigger
        WHERE tgname = 'update_rating_updated_at'
    ) THEN
        CREATE TRIGGER update_rating_updated_at
        BEFORE UPDATE ON rating
        FOR EACH ROW
        EXECUTE FUNCTION update_updated_at_column();
    END IF;
END
$$;

-- Add this after the update_post_rating_summary function definition

-- Create trigger for rating changes
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_trigger
        WHERE tgname = 'rating_post_summary_trigger'
    ) THEN
        CREATE TRIGGER rating_post_summary_trigger
        AFTER INSERT OR UPDATE OR DELETE ON rating
        FOR EACH ROW
        EXECUTE FUNCTION update_post_rating_summary();
    END IF;
END
$$;

-- Also, make sure post_rating_summary is created for new posts
CREATE OR REPLACE FUNCTION init_post_rating_summary() RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO post_rating_summary (post_id, total_rating)
    VALUES (NEW.id, 0);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_trigger
        WHERE tgname = 'init_post_rating_summary_trigger'
    ) THEN
        CREATE TRIGGER init_post_rating_summary_trigger
        AFTER INSERT ON post
        FOR EACH ROW
        EXECUTE FUNCTION init_post_rating_summary();
    END IF;
END
$$;