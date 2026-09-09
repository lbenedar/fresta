DROP TRIGGER IF EXISTS tr_delete_package_warnings_to_data_package_warnings_data;
DROP TRIGGER IF EXISTS tr_delete_game_to_systems_system;
DROP TRIGGER IF EXISTS tr_delete_game_to_pack_pack; 
DROP TRIGGER IF EXISTS tr_delete_game_to_module_module; 
DROP TRIGGER IF EXISTS tr_delete_pack_to_index_index;
DROP TABLE IF EXISTS token_turn_maker;
DROP TABLE IF EXISTS token_occludable;
DROP TABLE IF EXISTS token_light_darkness;
DROP TABLE IF EXISTS token_light_animation;
DROP TABLE IF EXISTS token_light;
DROP TABLE IF EXISTS token_bar_2;
DROP TABLE IF EXISTS token_bar_1;
DROP TABLE IF EXISTS token_texture;
DROP TABLE IF EXISTS token_sight;
DROP TABLE IF EXISTS ring_subject;
DROP TABLE IF EXISTS ring_colors;
DROP TABLE IF EXISTS ring;
DROP TABLE IF EXISTS token;
DROP TABLE IF EXISTS actor;
DROP TABLE IF EXISTS sound;
DROP TABLE IF EXISTS playlist;
DROP TABLE IF EXISTS table_result_range;
DROP TABLE IF EXISTS table_result;
DROP TABLE IF EXISTS table_;
DROP TABLE IF EXISTS journal_page_video;
DROP TABLE IF EXISTS journal_page_title;
DROP TABLE IF EXISTS journal_page_text;
DROP TABLE IF EXISTS journal_page;
DROP TABLE IF EXISTS journal;
DROP TABLE IF EXISTS setting;
DROP TABLE IF EXISTS item;
DROP TABLE IF EXISTS world_folder;
DROP TABLE IF EXISTS macro;
DROP TABLE IF EXISTS hotbar;
DROP TABLE IF EXISTS user;
DROP TABLE IF EXISTS face;
DROP TABLE IF EXISTS back;
DROP TABLE IF EXISTS card;
DROP TABLE IF EXISTS ownership_string;
DROP TABLE IF EXISTS card_deck;
DROP TABLE IF EXISTS combatant;
DROP TABLE IF EXISTS combat_groups;
DROP TABLE IF EXISTS combat;
DROP TABLE IF EXISTS message_rolls;
DROP TABLE IF EXISTS message_whisper;
DROP TABLE IF EXISTS speaker;
DROP TABLE IF EXISTS stats;
DROP TABLE IF EXISTS message;

CREATE TABLE IF NOT EXISTS pack_new (
    id VARCHAR(128) PRIMARY KEY,

    name VARCHAR(128) NOT NULL,
    label VARCHAR(128) NOT NULL,
    banner VARCHAR(128) NOT NULL,
    path VARCHAR(128) NOT NULL,
    type VARCHAR(128) NOT NULL,
    system VARCHAR(128) NOT NULL,
    package_type VARCHAR(128) NOT NULL,
    package_name VARCHAR(128) NOT NULL,

    module_id TEXT,
    FOREIGN KEY (module_id) REFERENCES module(id) ON DELETE CASCADE
);
ALTER TABLE pack_new ADD COLUMN system_id TEXT REFERENCES system(id) ON DELETE CASCADE;
ALTER TABLE pack_new ADD COLUMN world_id TEXT REFERENCES world(id) ON DELETE CASCADE;
INSERT INTO pack_new (id, name, label, banner, path, type, system, package_type, package_name, module_id, system_id, world_id)
SELECT id, name, label, banner, path, type, system, package_type, package_name, module_id, system_id, world_id FROM pack;
DROP TABLE IF EXISTS pack;
ALTER TABLE pack_new RENAME TO pack;

CREATE TABLE IF NOT EXISTS package_warnings_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    key_ TEXT NOT NULL,

    setup_id INTEGER,
    FOREIGN KEY (setup_id) REFERENCES setup(id) ON DELETE CASCADE
);
INSERT INTO package_warnings_new (id, key_, setup_id)
SELECT id, key_, setup_id FROM package_warnings;
DROP TABLE IF EXISTS package_warnings;
ALTER TABLE package_warnings_new RENAME TO package_warnings;

DROP TABLE IF EXISTS game_to_pack;
DROP TABLE IF EXISTS game_to_module;
DROP TABLE IF EXISTS active_users;
DROP TABLE IF EXISTS system_update;
DROP INDEX IF EXISTS idx_core_update_game_id;

CREATE TABLE IF NOT EXISTS core_update_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    has_update BOOLEAN NOT NULL DEFAULT FALSE,
    can_update BOOLEAN NOT NULL DEFAULT FALSE,
    could_reach_website BOOLEAN NOT NULL DEFAULT FALSE,
    slow_response BOOLEAN NOT NULL DEFAULT FALSE,
    will_disable_modules BOOLEAN NOT NULL DEFAULT FALSE,
    version VARCHAR(64) NOT NULL,
    channel VARCHAR(64) NOT NULL,

    setup_id INTEGER UNIQUE,
    FOREIGN KEY (setup_id) REFERENCES setup(id) ON DELETE CASCADE
);
INSERT INTO core_update_new (id, has_update, can_update, could_reach_website, slow_response, will_disable_modules, version, channel, setup_id)
SELECT id, has_update, can_update, could_reach_website, slow_response, will_disable_modules, version, channel, setup_id FROM core_update;
DROP TABLE IF EXISTS core_update;
ALTER TABLE core_update_new RENAME TO core_update;

DROP TABLE IF EXISTS game_to_systems;
DROP INDEX IF EXISTS idx_world_game_id;

CREATE TABLE IF NOT EXISTS world_new (
    id TEXT PRIMARY KEY,
    
    title VARCHAR(128) NOT NULL,
    description VARCHAR(128) NOT NULL,
    version VARCHAR(128) NOT NULL,
    system VARCHAR(128) NOT NULL,
    background VARCHAR(128) NOT NULL,
    join_theme VARCHAR(128) NOT NULL,
    core_version VARCHAR(128) NOT NULL,
    system_version VARCHAR(128) NOT NULL,
    last_played VARCHAR(128) NOT NULL,
    playtime INTEGER NOT NULL,
    availability INTEGER NOT NULL,
    next_session DATETIME NOT NULL,
    socket BOOLEAN NOT NULL DEFAULT FALSE,
    protected BOOLEAN NOT NULL DEFAULT FALSE,
    exclusive_ BOOLEAN NOT NULL DEFAULT FALSE,
    persistent_storage BOOLEAN NOT NULL DEFAULT FALSE,
    locked BOOLEAN NOT NULL DEFAULT FALSE,
    owned BOOLEAN NOT NULL DEFAULT FALSE,
    has_storage BOOLEAN NOT NULL DEFAULT FALSE,

    created_at TEXT DEFAULT (datetime('now')),
    updated_at TEXT DEFAULT (datetime('now')),

    setup_id INTEGER,
    FOREIGN KEY (setup_id) REFERENCES setup(id) ON DELETE CASCADE
);
INSERT INTO world_new (id, title, description, version, system, background, join_theme, core_version, system_version, last_played, playtime, availability, next_session, socket, protected, exclusive_, persistent_storage, locked, owned, has_storage, created_at, updated_at, setup_id)
SELECT id, title, description, version, system, background, join_theme, core_version, system_version, last_played, playtime, availability, next_session, socket, protected, exclusive_, persistent_storage, locked, owned, has_storage, created_at, updated_at, setup_id FROM world;
DROP TABLE IF EXISTS world;
ALTER TABLE world_new RENAME TO world;

DROP INDEX IF EXISTS idx_release_game_id;

CREATE TABLE IF NOT EXISTS release_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    generation INTEGER NOT NULL,
    build INTEGER NOT NULL,
    node_version INTEGER NOT NULL,
    max_generation INTEGER NOT NULL,
    max_stable_generation INTEGER NOT NULL,
    time DATETIME NOT NULL,
    channel VARCHAR(128) NOT NULL,
    suffix VARCHAR(128) NOT NULL,

    setup_id INTEGER UNIQUE,
    FOREIGN KEY (setup_id) REFERENCES setup(id) ON DELETE CASCADE
);
INSERT INTO release_new (id, generation, build, node_version, max_generation, max_stable_generation, time, channel, suffix, setup_id)
SELECT id, generation, build, node_version, max_generation, max_stable_generation, time, channel, suffix, setup_id FROM release_;
DROP TABLE IF EXISTS release_;
ALTER TABLE release_new RENAME TO release_;

DROP TABLE IF EXISTS game_options;
DROP INDEX IF EXISTS idx_files_game_id;

DROP TABLE IF EXISTS files_new;
CREATE TABLE IF NOT EXISTS files_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    setup_id INTEGER UNIQUE,
    FOREIGN KEY (setup_id) REFERENCES setup(id) ON DELETE CASCADE
);
INSERT INTO files_new (id, setup_id)
SELECT id, setup_id FROM files;
DROP TABLE IF EXISTS files;
ALTER TABLE files_new RENAME TO files;

DROP TABLE IF EXISTS addresses;
DROP TABLE IF EXISTS game;

CREATE TRIGGER tr_delete_package_warnings_to_data_package_warnings_data
AFTER DELETE ON package_warnings_to_data
FOR EACH ROW
BEGIN
    DELETE FROM package_warnings_data 
    WHERE id = OLD.package_warnings_data_id 
    AND NOT EXISTS (SELECT 1 FROM package_warnings_to_data WHERE package_warnings_data_id = OLD.package_warnings_data_id);
END;

CREATE TRIGGER tr_delete_pack_to_index_index
AFTER DELETE ON pack_to_index
FOR EACH ROW
BEGIN
    DELETE FROM index_ 
    WHERE id = OLD.index_id 
    AND NOT EXISTS (SELECT 1 FROM pack_to_index WHERE index_id = OLD.index_id);
END;