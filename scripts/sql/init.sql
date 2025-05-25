CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS friendship_statuses (
    id SERIAL PRIMARY KEY,
    name VARCHAR(20) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS friendships (
    id SERIAL PRIMARY KEY,
    user1_id BIGINT NOT NULL CHECK (user1_id < user2_id),
    user2_id BIGINT NOT NULL,
    status_id BIGINT NOT NULL,
    action_user_id BIGINT NOT NULL,
    
    UNIQUE (user1_id, user2_id),
    
    FOREIGN KEY (user1_id) REFERENCES users(id),
    FOREIGN KEY (user2_id) REFERENCES users(id),
    FOREIGN KEY (status_id) REFERENCES friendship_statuses(id),
    FOREIGN KEY (action_user_id) REFERENCES users(id)
);

INSERT INTO friendship_statuses (name) VALUES
    ('pending'),
    ('accepted'),
    ('declined')
ON CONFLICT (name) DO NOTHING;