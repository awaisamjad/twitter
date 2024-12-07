CREATE TABLE
  users (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    username VARCHAR(50) NOT NULL UNIQUE,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    bio TEXT CHECK (LENGTH (bio) <= 500),
    phone_number VARCHAR(15) CHECK (LENGTH (phone_number) <= 15)
  );

CREATE TABLE
  posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    content TEXT NOT NULL,
    username VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    like_num INTEGER DEFAULT 0,
    dislike_num INTEGER DEFAULT 0,
    FOREIGN KEY (username) REFERENCES users (username)
  );

CREATE TABLE
  user_relationships (
    user_id INTEGER NOT NULL,
    following_id INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (following_id) REFERENCES users (id),
    PRIMARY KEY (user_id, following_id)
  );