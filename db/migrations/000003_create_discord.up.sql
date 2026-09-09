CREATE TABLE IF NOT EXISTS discord_interaction (
    id TEXT PRIMARY KEY,
    
    application_id TEXT NOT NULL,
    guild_id TEXT NOT NULL,
    
    created_at DATETIME NOT NULL DEFAULT current_timestamp
);

ALTER TABLE discord_interaction ADD COLUMN message_id TEXT REFERENCES discord_message(id);
ALTER TABLE discord_interaction ADD COLUMN member_id INTEGER REFERENCES discord_member(id);
ALTER TABLE discord_interaction ADD COLUMN user_id TEXT REFERENCES discord_user(id);

CREATE TABLE IF NOT EXISTS discord_message (
    id TEXT PRIMARY KEY,
    
    guild_id TEXT NOT NULL,
    content TEXT NOT NULL,
    timestamp DATETIME NOT NULL DEFAULT current_timestamp,
    edited_timestamp DATETIME DEFAULT current_timestamp
);

CREATE TABLE IF NOT EXISTS discord_member (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    
    guild_id TEXT NOT NULL,
    joined_at DATETIME NOT NULL DEFAULT current_timestamp,
    nick TEXT NOT NULL,
    deaf BOOLEAN NOT NULL,
    mute BOOLEAN NOT NULL,

    created_at DATETIME NOT NULL DEFAULT current_timestamp,
    updated_at DATETIME NOT NULL DEFAULT current_timestamp,

    user_id TEXT UNIQUE,
    FOREIGN KEY (user_id) REFERENCES discord_user(id)
);

CREATE TABLE IF NOT EXISTS discord_user (
    id TEXT PRIMARY KEY,
    
    email TEXT NOT NULL,
    username TEXT NOT NULL,
    avatar TEXT NOT NULL,
    locale TEXT NOT NULL,
    discriminator TEXT NOT NULL,
    global_name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS discord_to_foundry_user (
    discord_user_id TEXT,
    foundry_user_id TEXT,
    is_verified BOOLEAN DEFAULT FALSE,
    
    PRIMARY KEY (discord_user_id, foundry_user_id),
    FOREIGN KEY (discord_user_id) REFERENCES discord_user(id) ON DELETE CASCADE,
    FOREIGN KEY (foundry_user_id) REFERENCES user(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS discord_command_history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    
    num INTEGER NOT NULL,
    command TEXT NOT NULL,

    msg_id TEXT,
    FOREIGN KEY (msg_id) REFERENCES discord_message(id)
);