CREATE TABLE IF NOT EXISTS game (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    
    demo_mode BOOLEAN NOT NULL DEFAULT FALSE,
    idle_logout BOOLEAN NOT NULL DEFAULT FALSE,
    paused BOOLEAN NOT NULL DEFAULT FALSE,
    user_id VARCHAR(128) NOT NULL,

    created_at DATETIME NOT NULL DEFAULT current_timestamp
);

CREATE TABLE IF NOT EXISTS addresses (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    local VARCHAR(128) NOT NULL,
    remote VARCHAR(128) NOT NULL,
    remote_is_accessible BOOLEAN NOT NULL DEFAULT FALSE,

    game_id INTEGER UNIQUE,
    FOREIGN KEY (game_id) REFERENCES game(id) ON DELETE CASCADE
);

ALTER TABLE files ADD COLUMN game_id INTEGER REFERENCES game(id) ON DELETE CASCADE;
CREATE UNIQUE INDEX idx_files_game_id ON files(game_id);

CREATE TABLE IF NOT EXISTS game_options (
    id INTEGER PRIMARY KEY,

    language VARCHAR(128) NOT NULL,
    update_channel VARCHAR(128) NOT NULL,
    port INTEGER NOT NULL,

    game_id INTEGER UNIQUE,
    FOREIGN KEY (game_id) REFERENCES game(id) ON DELETE CASCADE
);

ALTER TABLE release_ ADD COLUMN game_id INTEGER REFERENCES game(id) ON DELETE CASCADE;
CREATE UNIQUE INDEX idx_release_game_id ON release_(game_id);
ALTER TABLE world ADD COLUMN game_id INTEGER REFERENCES game(id) ON DELETE CASCADE;
CREATE UNIQUE INDEX idx_world_game_id ON world(game_id);

CREATE TABLE IF NOT EXISTS game_to_systems (
    game_id INTEGER,
    system_id TEXT,
    
    PRIMARY KEY (game_id, system_id),
    FOREIGN KEY (game_id) REFERENCES game(id) ON DELETE CASCADE,
    FOREIGN KEY (system_id) REFERENCES system(id) ON DELETE CASCADE
);

CREATE INDEX idx_game_to_systems_system_id ON game_to_systems(system_id);

CREATE TRIGGER tr_delete_game_to_systems_system
AFTER DELETE ON game_to_systems
FOR EACH ROW
BEGIN
    DELETE FROM system 
    WHERE id = OLD.system_id 
    AND NOT EXISTS (SELECT 1 FROM game_to_systems WHERE system_id = OLD.system_id);
END;

ALTER TABLE core_update ADD COLUMN game_id INTEGER REFERENCES game(id) ON DELETE CASCADE;
CREATE UNIQUE INDEX idx_core_update_game_id ON core_update(game_id);  

CREATE TABLE IF NOT EXISTS system_update (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    has_update BOOLEAN NOT NULL DEFAULT FALSE,
    version VARCHAR(128) NOT NULL,

    game_id INTEGER UNIQUE,
    FOREIGN KEY (game_id) REFERENCES game(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS active_users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    value VARCHAR(128) NOT NULL,

    game_id INTEGER,
    FOREIGN KEY (game_id) REFERENCES game(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS game_to_module (
    game_id INTEGER,
    module_id TEXT,
    
    PRIMARY KEY (game_id, module_id),
    FOREIGN KEY (module_id) REFERENCES module(id) ON DELETE CASCADE,
    FOREIGN KEY (game_id) REFERENCES game(id) ON DELETE CASCADE
);

CREATE INDEX idx_game_to_module_module_id ON game_to_module(module_id);

CREATE TRIGGER tr_delete_game_to_module_module
AFTER DELETE ON game_to_module
FOR EACH ROW
BEGIN
    DELETE FROM module 
    WHERE id = OLD.module_id 
    AND NOT EXISTS (SELECT 1 FROM game_to_module WHERE module_id = OLD.module_id);
END;

ALTER TABLE package_warnings ADD COLUMN game_id INTEGER REFERENCES game(id) ON DELETE CASCADE;

CREATE TABLE IF NOT EXISTS game_to_pack (
    game_id INTEGER,
    pack_id TEXT,
    
    PRIMARY KEY (game_id, pack_id),
    FOREIGN KEY (pack_id) REFERENCES pack(id) ON DELETE CASCADE,
    FOREIGN KEY (game_id) REFERENCES game(id) ON DELETE CASCADE
);

CREATE INDEX idx_game_to_pack_pack_id ON game_to_pack(pack_id);

CREATE TRIGGER tr_delete_game_to_pack_pack
AFTER DELETE ON game_to_pack
FOR EACH ROW
BEGIN
    DELETE FROM pack 
    WHERE id = OLD.pack_id 
    AND NOT EXISTS (SELECT 1 FROM game_to_pack WHERE pack_id = OLD.pack_id);
END;

CREATE TABLE IF NOT EXISTS message (
    id TEXT PRIMARY KEY,

    blind BOOLEAN NOT NULL,
    emote BOOLEAN NOT NULL,
    style INTEGER NOT NULL,
    timestamp DATETIME NOT NULL,
    content VARCHAR(128) NOT NULL,
    author VARCHAR(128) NOT NULL,
    type VARCHAR(128) NOT NULL,
    flavor VARCHAR(128) NOT NULL,
    sound VARCHAR(128) NOT NULL,

    created_at TEXT DEFAULT (datetime('now')),
    updated_at TEXT DEFAULT (datetime('now')),

    game_id INTEGER,
    FOREIGN KEY (game_id) REFERENCES game(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS stats (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    core_version VARCHAR(128) NOT NULL,
    system_id VARCHAR(128) NOT NULL,
    system_version VARCHAR(128) NOT NULL,
    last_modified_by VARCHAR(128) NOT NULL,
    modified_time DATETIME NOT NULL,

    message_id TEXT UNIQUE,
    FOREIGN KEY (message_id) REFERENCES message(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS speaker (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    scene VARCHAR(128) NOT NULL,
    actor VARCHAR(128) NOT NULL,
    token VARCHAR(128) NOT NULL,
    alias VARCHAR(128) NOT NULL,

    message_id TEXT UNIQUE,
    FOREIGN KEY (message_id) REFERENCES message(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS message_whisper (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    value VARCHAR(128) NOT NULL,

    message_id TEXT,
    FOREIGN KEY (message_id) REFERENCES message(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS message_rolls (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    value VARCHAR(128) NOT NULL,

    message_id TEXT,
    FOREIGN KEY (message_id) REFERENCES message(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS combat (
    id TEXT PRIMARY KEY,

    type VARCHAR(128) NOT NULL,
    scene VARCHAR(128) NOT NULL,
    round INTEGER NOT NULL,
    turn INTEGER NOT NULL,
    sort INTEGER NOT NULL,
    active INTEGER NOT NULL,

    created_at TEXT DEFAULT (datetime('now')),
    updated_at TEXT DEFAULT (datetime('now')),

    game_id INTEGER,
    FOREIGN KEY (game_id) REFERENCES game(id) ON DELETE CASCADE
);

ALTER TABLE stats ADD COLUMN combat_id TEXT REFERENCES combat(id) ON DELETE CASCADE;
CREATE UNIQUE INDEX idx_stats_combat_id ON stats(combat_id);

CREATE TABLE IF NOT EXISTS combat_groups (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    value VARCHAR(128) NOT NULL,

    combat_id TEXT,
    FOREIGN KEY (combat_id) REFERENCES combat(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS combatant (
    id TEXT PRIMARY KEY,

    token_id VARCHAR(128) NOT NULL,
    scene_id VARCHAR(128) NOT NULL,
    actor_id VARCHAR(128) NOT NULL,
    type VARCHAR(128) NOT NULL,
    img VARCHAR(128) NOT NULL,
    group_ VARCHAR(128) NOT NULL,
    initiative INTEGER NOT NULL,
    hidden BOOLEAN NOT NULL,
    defeated BOOLEAN NOT NULL,

    created_at TEXT DEFAULT (datetime('now')),
    updated_at TEXT DEFAULT (datetime('now')),

    combat_id TEXT,
    FOREIGN KEY (combat_id) REFERENCES combat(id) ON DELETE CASCADE
);

ALTER TABLE stats ADD COLUMN combatant_id TEXT REFERENCES combatant(id) ON DELETE CASCADE;
CREATE UNIQUE INDEX idx_stats_combatant_id ON stats(combatant_id);

CREATE TABLE IF NOT EXISTS card_deck (
    id TEXT PRIMARY KEY,

    name VARCHAR(128) NOT NULL,
    type VARCHAR(128) NOT NULL,
    description VARCHAR(128) NOT NULL,
    img VARCHAR(128) NOT NULL,
    folder VARCHAR(128) NOT NULL,
    width INTEGER NOT NULL,
    height INTEGER NOT NULL,
    rotation INTEGER NOT NULL,
    sort INTEGER NOT NULL,
    display_count BOOLEAN NOT NULL,

    created_at TEXT DEFAULT (datetime('now')),
    updated_at TEXT DEFAULT (datetime('now')),

    game_id INTEGER,
    FOREIGN KEY (game_id) REFERENCES game(id) ON DELETE CASCADE
);

ALTER TABLE stats ADD COLUMN card_deck_id TEXT REFERENCES card_deck(id) ON DELETE CASCADE;
CREATE UNIQUE INDEX idx_stats_card_deck_id ON stats(card_deck_id);

CREATE TABLE IF NOT EXISTS ownership_string (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    key_ VARCHAR(128) NOT NULL,
    value INTEGER NOT NULL,

    card_deck_id TEXT,
    FOREIGN KEY (card_deck_id) REFERENCES card_deck(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS card (
    id TEXT PRIMARY KEY,

    name VARCHAR(128) NOT NULL,
    type VARCHAR(128) NOT NULL,
    suit VARCHAR(128) NOT NULL,
    description VARCHAR(128) NOT NULL,
    origin VARCHAR(128) NOT NULL,
    width INTEGER NOT NULL,
    height INTEGER NOT NULL,
    rotation INTEGER NOT NULL,
    value INTEGER NOT NULL,
    face INTEGER NOT NULL,
    sort INTEGER NOT NULL,
    drawn BOOLEAN NOT NULL,

    created_at TEXT DEFAULT (datetime('now')),
    updated_at TEXT DEFAULT (datetime('now')),

    card_deck_id TEXT,
    FOREIGN KEY (card_deck_id) REFERENCES card_deck(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS back (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    name VARCHAR(128) NOT NULL,
    text VARCHAR(128) NOT NULL,

    card_id TEXT UNIQUE,
    FOREIGN KEY (card_id) REFERENCES card(id) ON DELETE CASCADE
);

ALTER TABLE stats ADD COLUMN card_id TEXT REFERENCES card(id) ON DELETE CASCADE;
CREATE UNIQUE INDEX idx_stats_card_id ON stats(card_id);

CREATE TABLE IF NOT EXISTS face (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    name VARCHAR(128) NOT NULL,
    img VARCHAR(128) NOT NULL,
    text VARCHAR(128) NOT NULL,

    card_id TEXT,
    FOREIGN KEY (card_id) REFERENCES card(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS user (
    id TEXT PRIMARY KEY,

    name VARCHAR(128) NOT NULL,
    avatar VARCHAR(128) NOT NULL,
    character VARCHAR(128) NOT NULL,
    color VARCHAR(128) NOT NULL,
    pronouns VARCHAR(128) NOT NULL,
    role INTEGER NOT NULL,

    created_at TEXT DEFAULT (datetime('now')),
    updated_at TEXT DEFAULT (datetime('now')),

    game_id INTEGER,
    FOREIGN KEY (game_id) REFERENCES game(id) ON DELETE CASCADE
);

ALTER TABLE stats ADD COLUMN user_id TEXT REFERENCES user(id) ON DELETE CASCADE;
CREATE UNIQUE INDEX idx_stats_user_id ON stats(user_id);

CREATE TABLE IF NOT EXISTS hotbar (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    key_ INTEGER NOT NULL,
    value VARCHAR(128) NOT NULL,

    user_id TEXT,
    FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS macro (
    id TEXT PRIMARY KEY,

    command VARCHAR(128) NOT NULL,
    name VARCHAR(128) NOT NULL,
    type VARCHAR(128) NOT NULL,
    img VARCHAR(128) NOT NULL,
    author VARCHAR(128) NOT NULL,
    scope VARCHAR(128) NOT NULL,
    folder VARCHAR(128) NOT NULL,
    sort INTEGER NOT NULL,

    created_at TEXT DEFAULT (datetime('now')),
    updated_at TEXT DEFAULT (datetime('now')),

    game_id INTEGER,
    FOREIGN KEY (game_id) REFERENCES game(id) ON DELETE CASCADE
);

ALTER TABLE stats ADD COLUMN macro_id TEXT REFERENCES macro(id) ON DELETE CASCADE;
CREATE UNIQUE INDEX idx_stats_macro_id ON stats(macro_id);
ALTER TABLE ownership_string ADD COLUMN macro_id TEXT REFERENCES macro(id) ON DELETE CASCADE;

CREATE TABLE IF NOT EXISTS world_folder (
    id TEXT PRIMARY KEY,

    name VARCHAR(128) NOT NULL,
    type VARCHAR(128) NOT NULL,
    folder VARCHAR(128) NOT NULL,
    sorting VARCHAR(128) NOT NULL,
    description VARCHAR(128) NOT NULL,
    color VARCHAR(128) NOT NULL,
    sort INTEGER NOT NULL,

    created_at TEXT DEFAULT (datetime('now')),
    updated_at TEXT DEFAULT (datetime('now')),

    game_id INTEGER,
    FOREIGN KEY (game_id) REFERENCES game(id) ON DELETE CASCADE
);

ALTER TABLE stats ADD COLUMN world_folder_id TEXT REFERENCES world_folder(id) ON DELETE CASCADE;
CREATE UNIQUE INDEX idx_stats_world_folder_id ON stats(world_folder_id);

CREATE TABLE IF NOT EXISTS item (
    id TEXT PRIMARY KEY,

    img VARCHAR(128) NOT NULL,
    name VARCHAR(128) NOT NULL,
    type VARCHAR(128) NOT NULL,
    folder VARCHAR(128) NOT NULL,
    sort INTEGER NOT NULL,

    created_at TEXT DEFAULT (datetime('now')),
    updated_at TEXT DEFAULT (datetime('now')),

    game_id INTEGER,
    FOREIGN KEY (game_id) REFERENCES game(id) ON DELETE CASCADE
);

ALTER TABLE stats ADD COLUMN item_id TEXT REFERENCES item(id) ON DELETE CASCADE;
CREATE UNIQUE INDEX idx_stats_item_id ON stats(item_id);
ALTER TABLE ownership_string ADD COLUMN item_id TEXT REFERENCES item(id) ON DELETE CASCADE;

CREATE TABLE IF NOT EXISTS setting (
    id TEXT PRIMARY KEY,

    key_ VARCHAR(128) NOT NULL,
    value VARCHAR(128) NOT NULL,

    created_at TEXT DEFAULT (datetime('now')),
    updated_at TEXT DEFAULT (datetime('now')),

    game_id INTEGER,
    FOREIGN KEY (game_id) REFERENCES game(id) ON DELETE CASCADE
);

ALTER TABLE stats ADD COLUMN setting_id TEXT REFERENCES setting(id) ON DELETE CASCADE;
CREATE UNIQUE INDEX idx_stats_setting_id ON stats(setting_id);

CREATE TABLE IF NOT EXISTS journal (
    id TEXT PRIMARY KEY,

    name VARCHAR(128) NOT NULL,
    sort INTEGER NOT NULL,

    created_at TEXT DEFAULT (datetime('now')),
    updated_at TEXT DEFAULT (datetime('now')),

    game_id INTEGER,
    FOREIGN KEY (game_id) REFERENCES game(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS journal_page (
    id TEXT PRIMARY KEY,

    name VARCHAR(128) NOT NULL,
    type VARCHAR(128) NOT NULL,
    src VARCHAR(128) NOT NULL,
    sort INTEGER NOT NULL,

    created_at TEXT DEFAULT (datetime('now')),
    updated_at TEXT DEFAULT (datetime('now')),

    journal_id TEXT,
    FOREIGN KEY (journal_id) REFERENCES journal(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS journal_page_text (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    content VARCHAR(128) NOT NULL,
    markdown VARCHAR(128) NOT NULL,
    format INTEGER NOT NULL,

    journal_page_id TEXT,
    FOREIGN KEY (journal_page_id) REFERENCES journal_page(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS journal_page_title (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    show BOOLEAN NOT NULL,
    level INTEGER NOT NULL,

    journal_page_id TEXT,
    FOREIGN KEY (journal_page_id) REFERENCES journal_page(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS journal_page_video (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    controls BOOLEAN NOT NULL,
    volume REAL NOT NULL,

    journal_page_id TEXT,
    FOREIGN KEY (journal_page_id) REFERENCES journal_page(id) ON DELETE CASCADE
);

ALTER TABLE stats ADD COLUMN journal_page_id TEXT REFERENCES journal_page(id) ON DELETE CASCADE;
ALTER TABLE ownership_string ADD COLUMN journal_page_id TEXT REFERENCES journal_page(id) ON DELETE CASCADE;

ALTER TABLE ownership_string ADD COLUMN journal_id TEXT REFERENCES journal(id) ON DELETE CASCADE;

CREATE TABLE IF NOT EXISTS table_ (
    id TEXT PRIMARY KEY,

    name VARCHAR(128) NOT NULL,
    description VARCHAR(128) NOT NULL,
    formula VARCHAR(128) NOT NULL,
    img VARCHAR(128) NOT NULL,
    folder VARCHAR(128) NOT NULL,
    sort INTEGER NOT NULL,
    replacement BOOLEAN NOT NULL,
    display_roll BOOLEAN NOT NULL,

    created_at TEXT DEFAULT (datetime('now')),
    updated_at TEXT DEFAULT (datetime('now')),

    game_id INTEGER,
    FOREIGN KEY (game_id) REFERENCES game(id) ON DELETE CASCADE
);

ALTER TABLE stats ADD COLUMN table_id TEXT REFERENCES table_(id) ON DELETE CASCADE;
CREATE UNIQUE INDEX idx_stats_table_id ON stats(table_id);
ALTER TABLE ownership_string ADD COLUMN table_id TEXT REFERENCES table_(id) ON DELETE CASCADE;

CREATE TABLE IF NOT EXISTS table_result (
    id TEXT PRIMARY KEY,

    type VARCHAR(128) NOT NULL,
    img VARCHAR(128) NOT NULL,
    description VARCHAR(128) NOT NULL,
    name VARCHAR(128) NOT NULL,
    weight INTEGER NOT NULL,
    drawn BOOLEAN NOT NULL,

    created_at TEXT DEFAULT (datetime('now')),
    updated_at TEXT DEFAULT (datetime('now')),

    table_id TEXT,
    FOREIGN KEY (table_id) REFERENCES table_(id) ON DELETE CASCADE
);

ALTER TABLE stats ADD COLUMN table_result_id TEXT REFERENCES table_result(id) ON DELETE CASCADE;
CREATE UNIQUE INDEX idx_stats_table_result_id ON stats(table_result_id);

CREATE TABLE IF NOT EXISTS table_result_range (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    value INTEGER NOT NULL,

    table_result_id TEXT,
    FOREIGN KEY (table_result_id) REFERENCES table_result(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS playlist (
    id TEXT PRIMARY KEY,

    name VARCHAR(128) NOT NULL,
    folder VARCHAR(128) NOT NULL,
    sorting VARCHAR(128) NOT NULL,
    description VARCHAR(128) NOT NULL,
    channel VARCHAR(128) NOT NULL,
    mode INTEGER NOT NULL,
    fade INTEGER NOT NULL,
    seed INTEGER NOT NULL,
    sort INTEGER NOT NULL,
    playing BOOLEAN NOT NULL,

    created_at TEXT DEFAULT (datetime('now')),
    updated_at TEXT DEFAULT (datetime('now')),

    game_id INTEGER,
    FOREIGN KEY (game_id) REFERENCES game(id) ON DELETE CASCADE
);

ALTER TABLE stats ADD COLUMN playlist_id TEXT REFERENCES playlist(id) ON DELETE CASCADE;
CREATE UNIQUE INDEX idx_stats_playlist_id ON stats(playlist_id);
ALTER TABLE ownership_string ADD COLUMN playlist_id TEXT REFERENCES playlist(id) ON DELETE CASCADE;

CREATE TABLE IF NOT EXISTS sound (
    id TEXT PRIMARY KEY,

    name VARCHAR(128) NOT NULL,
    path VARCHAR(128) NOT NULL,
    channel VARCHAR(128) NOT NULL,
    description VARCHAR(128) NOT NULL,
    fade INTEGER NOT NULL,
    sort INTEGER NOT NULL,
    repeat BOOLEAN NOT NULL,
    playing BOOLEAN NOT NULL,
    volume REAL NOT NULL,
    paused_time REAL NOT NULL,

    created_at TEXT DEFAULT (datetime('now')),
    updated_at TEXT DEFAULT (datetime('now')),

    playlist_id TEXT,
    FOREIGN KEY (playlist_id) REFERENCES playlist(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS actor (
    id TEXT PRIMARY KEY,

    img VARCHAR(128) NOT NULL,
    name VARCHAR(128) NOT NULL,
    type VARCHAR(128) NOT NULL,
    folder VARCHAR(128) NOT NULL,
    sort INTEGER NOT NULL,

    created_at TEXT DEFAULT (datetime('now')),
    updated_at TEXT DEFAULT (datetime('now')),

    game_id INTEGER,
    FOREIGN KEY (game_id) REFERENCES game(id) ON DELETE CASCADE
);

ALTER TABLE stats ADD COLUMN actor_id TEXT REFERENCES actor(id) ON DELETE CASCADE;
CREATE UNIQUE INDEX idx_stats_actor_id ON stats(actor_id);
ALTER TABLE ownership_string ADD COLUMN actor_id TEXT REFERENCES actor(id) ON DELETE CASCADE;
ALTER TABLE item ADD COLUMN actor_id TEXT REFERENCES actor(id) ON DELETE CASCADE;

CREATE TABLE IF NOT EXISTS token (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    name VARCHAR(128) NOT NULL,
    actor_link BOOLEAN NOT NULL,
    append_number BOOLEAN NOT NULL,
    prepend_adjective BOOLEAN NOT NULL,
    lock_rotation BOOLEAN NOT NULL,
    random_img BOOLEAN NOT NULL,
    display_name INTEGER NOT NULL,
    display_bars INTEGER NOT NULL,
    disposition INTEGER NOT NULL,
    rotation INTEGER NOT NULL,
    alpha INTEGER NOT NULL,
    width REAL NOT NULL,
    height REAL NOT NULL,
    
    actor_id INTEGER UNIQUE,
    FOREIGN KEY (actor_id) REFERENCES actor(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS ring (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    enabled BOOLEAN NOT NULL,
    effects INTEGER NOT NULL,
    
    token_id INTEGER UNIQUE,
    FOREIGN KEY (token_id) REFERENCES token(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS ring_colors (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    ring VARCHAR(128) NOT NULL,
    background VARCHAR(128) NOT NULL,
    
    ring_id INTEGER UNIQUE,
    FOREIGN KEY (ring_id) REFERENCES ring(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS ring_subject (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    scale INTEGER NOT NULL,
    texture VARCHAR(128) NOT NULL,
    
    ring_id INTEGER UNIQUE,
    FOREIGN KEY (ring_id) REFERENCES ring(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS token_sight (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    color VARCHAR(128) NOT NULL,
    vision_mode VARCHAR(128) NOT NULL,
    range_ INTEGER NOT NULL,
    angle INTEGER NOT NULL,
    attenuation REAL NOT NULL,
    brightness REAL NOT NULL,
    enabled BOOLEAN NOT NULL,
    
    token_id INTEGER UNIQUE,
    FOREIGN KEY (token_id) REFERENCES token(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS token_texture (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    src VARCHAR(128) NOT NULL,
    fit VARCHAR(128) NOT NULL,
    tint VARCHAR(128) NOT NULL,
    scale_x REAL NOT NULL,
    scale_y REAL NOT NULL,
    offset_x REAL NOT NULL,
    offset_y REAL NOT NULL,
    rotation REAL NOT NULL,
    anchor_x REAL NOT NULL,
    anchor_y REAL NOT NULL,
    alpha_threshold REAL NOT NULL,
    
    token_id INTEGER UNIQUE,
    FOREIGN KEY (token_id) REFERENCES token(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS token_bar_1 (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    attribute VARCHAR(128) NOT NULL,
    
    token_id INTEGER UNIQUE,
    FOREIGN KEY (token_id) REFERENCES token(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS token_bar_2 (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    attribute VARCHAR(128) NOT NULL,
    
    token_id INTEGER UNIQUE,
    FOREIGN KEY (token_id) REFERENCES token(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS token_light (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    color VARCHAR(128) NOT NULL,
    priority INTEGER NOT NULL,
    angle INTEGER NOT NULL,
    negative BOOLEAN NOT NULL,
    alpha REAL NOT NULL,
    bright REAL NOT NULL,
    coloration REAL NOT NULL,
    dim REAL NOT NULL,
    attenuation REAL NOT NULL,
    luminosity REAL NOT NULL,
    saturation REAL NOT NULL,
    contrast REAL NOT NULL,
    shadows REAL NOT NULL,
    
    token_id INTEGER UNIQUE,
    FOREIGN KEY (token_id) REFERENCES token(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS token_light_animation (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    speed INTEGER NOT NULL,
    intensity INTEGER NOT NULL,
    reverse BOOLEAN NOT NULL,

    token_light_id INTEGER UNIQUE,
    FOREIGN KEY (token_light_id) REFERENCES token_light(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS token_light_darkness (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    min REAL NOT NULL,
    max REAL NOT NULL,

    token_light_id INTEGER UNIQUE,
    FOREIGN KEY (token_light_id) REFERENCES token_light(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS token_occludable (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    radius INTEGER NOT NULL,
    
    token_id INTEGER UNIQUE,
    FOREIGN KEY (token_id) REFERENCES token(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS token_turn_maker (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    mode INTEGER NOT NULL,
    animation VARCHAR(128) NOT NULL,
    src VARCHAR(128) NOT NULL,
    disposition BOOLEAN NOT NULL,
    
    token_id INTEGER UNIQUE,
    FOREIGN KEY (token_id) REFERENCES token(id) ON DELETE CASCADE
);
