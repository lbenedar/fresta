CREATE TABLE IF NOT EXISTS setup (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    
    is_admin BOOLEAN NOT NULL DEFAULT FALSE,
    is_setup BOOLEAN NOT NULL DEFAULT TRUE,

    created_at DATETIME NOT NULL DEFAULT current_timestamp
);

CREATE TABLE IF NOT EXISTS core_update (
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

CREATE TABLE IF NOT EXISTS featured_content (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    title VARCHAR(128) NOT NULL,
    caption VARCHAR(128) NOT NULL,
    url VARCHAR(128) NOT NULL,
    image VARCHAR(128) NOT NULL,

    setup_id INTEGER UNIQUE,
    FOREIGN KEY (setup_id) REFERENCES setup(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS files (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    setup_id INTEGER UNIQUE,
    FOREIGN KEY (setup_id) REFERENCES setup(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS files_storage (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    storage TEXT NOT NULL,

    files_id INTEGER,
    FOREIGN KEY (files_id) REFERENCES files(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS setup_options (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    css_theme VARCHAR(128) NOT NULL,
    data_path VARCHAR(128) NOT NULL,
    hostname VARCHAR(128) NOT NULL,
    language VARCHAR(128) NOT NULL,
    local_hostname VARCHAR(128) NOT NULL,
    update_channel VARCHAR(128) NOT NULL,
    port INTEGER NOT NULL,
    compress_socket BOOLEAN NOT NULL DEFAULT FALSE,
    compress_static BOOLEAN NOT NULL DEFAULT FALSE,
    fullscreen BOOLEAN NOT NULL DEFAULT FALSE,
    hot_reload BOOLEAN NOT NULL DEFAULT FALSE,
    proxy_ssl BOOLEAN NOT NULL,
    telemetry BOOLEAN NOT NULL,
    upnp BOOLEAN NOT NULL,
    delete_nedb BOOLEAN NOT NULL,
    no_backups BOOLEAN NOT NULL,

    setup_id INTEGER UNIQUE,
    FOREIGN KEY (setup_id) REFERENCES setup(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS release_ (
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

CREATE TABLE IF NOT EXISTS setup_language (
    id INTEGER PRIMARY KEY,

    label VARCHAR(128) NOT NULL,

    setup_id INTEGER,
    FOREIGN KEY (setup_id) REFERENCES setup(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS setup_language_module (
    id TEXT PRIMARY KEY,

    label VARCHAR(128) NOT NULL,
    path VARCHAR(128) NOT NULL,

    setup_language_id INTEGER,
    FOREIGN KEY (setup_language_id) REFERENCES setup_language(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS module (
    id TEXT PRIMARY KEY,
    
    title VARCHAR(128) NOT NULL,
    description VARCHAR(128) NOT NULL,
    url VARCHAR(128) NOT NULL,
    license VARCHAR(128) NOT NULL,
    readme VARCHAR(128) NOT NULL,
    bugs VARCHAR(128) NOT NULL,
    changelog VARCHAR(128) NOT NULL,
    version VARCHAR(128) NOT NULL,
    manifest VARCHAR(128) NOT NULL,
    download VARCHAR(128) NOT NULL,
    socket BOOLEAN NOT NULL DEFAULT FALSE,
    protected BOOLEAN NOT NULL DEFAULT FALSE,
    exclusive_ BOOLEAN NOT NULL DEFAULT FALSE,
    persistent_storage BOOLEAN NOT NULL DEFAULT FALSE,
    core_translation BOOLEAN NOT NULL DEFAULT FALSE,
    library BOOLEAN NOT NULL DEFAULT FALSE,
    locked BOOLEAN NOT NULL DEFAULT FALSE,
    owned BOOLEAN NOT NULL DEFAULT FALSE,
    has_storage BOOLEAN NOT NULL DEFAULT FALSE,
    active BOOLEAN NOT NULL DEFAULT FALSE,
    availability INTEGER NOT NULL DEFAULT FALSE,

    setup_id INTEGER,
    FOREIGN KEY (setup_id) REFERENCES setup(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS document_types (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    module_id TEXT UNIQUE,
    FOREIGN KEY (module_id) REFERENCES module(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS document_types_data (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    type VARCHAR(64) NOT NULL,

    document_types_id INTEGER,
    FOREIGN KEY (document_types_id) REFERENCES document_types(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS document_types_data_html (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    value TEXT NOT NULL,

    document_types_data_id INTEGER,
    FOREIGN KEY (document_types_data_id) REFERENCES document_types_data(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS relationships (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    module_id TEXT UNIQUE,
    FOREIGN KEY (module_id) REFERENCES module(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS relationships_systems (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    key_ TEXT NOT NULL,
    type TEXT NOT NULL,
    manifest TEXT NOT NULL,

    relationships_id INTEGER,
    FOREIGN KEY (relationships_id) REFERENCES relationships(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS relationships_requires (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    key_ INTEGER NOT NULL,
    type TEXT NOT NULL,
    manifest TEXT NOT NULL,

    relationships_id INTEGER,
    FOREIGN KEY (relationships_id) REFERENCES relationships(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS relationships_recommends (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    key_ INTEGER NOT NULL,
    type TEXT NOT NULL,
    manifest TEXT NOT NULL,

    relationships_id INTEGER,
    FOREIGN KEY (relationships_id) REFERENCES relationships(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS relationships_conflicts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    key_ INTEGER NOT NULL,
    type TEXT NOT NULL,
    manifest TEXT NOT NULL,

    relationships_id INTEGER,
    FOREIGN KEY (relationships_id) REFERENCES relationships(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS compatibility (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    minimum VARCHAR(64) NOT NULL,
    verified VARCHAR(64) NOT NULL,
    maximum VARCHAR(64) NOT NULL,

    module_id TEXT UNIQUE,
    FOREIGN KEY (module_id) REFERENCES module(id) ON DELETE CASCADE
);

ALTER TABLE compatibility ADD COLUMN relationships_systems_id TEXT REFERENCES relationships_systems(id) ON DELETE CASCADE;
ALTER TABLE compatibility ADD COLUMN relationships_requires_id TEXT REFERENCES relationships_requires(id) ON DELETE CASCADE;
ALTER TABLE compatibility ADD COLUMN relationships_recommends_id TEXT REFERENCES relationships_recommends(id) ON DELETE CASCADE;
ALTER TABLE compatibility ADD COLUMN relationships_conflicts_id TEXT REFERENCES relationships_conflicts(id) ON DELETE CASCADE;

CREATE TABLE IF NOT EXISTS scripts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    value TEXT NOT NULL,

    module_id TEXT,
    FOREIGN KEY (module_id) REFERENCES module(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS es_modules (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    value TEXT NOT NULL,

    module_id TEXT,
    FOREIGN KEY (module_id) REFERENCES module(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS tags (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    value VARCHAR(128) NOT NULL,

    module_id TEXT,
    FOREIGN KEY (module_id) REFERENCES module(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS author (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    name VARCHAR(128) NOT NULL,
    url VARCHAR(128) NOT NULL,
    email VARCHAR(128) NOT NULL,
    discord VARCHAR(128) NOT NULL,

    module_id TEXT,
    FOREIGN KEY (module_id) REFERENCES module(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS media (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    type VARCHAR(128) NOT NULL,
    url VARCHAR(128) NOT NULL,
    caption VARCHAR(128) NOT NULL,

    module_id TEXT,
    FOREIGN KEY (module_id) REFERENCES module(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS style (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    src VARCHAR(128) NOT NULL,

    module_id TEXT,
    FOREIGN KEY (module_id) REFERENCES module(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS language (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    lang VARCHAR(128) NOT NULL,
    name VARCHAR(128) NOT NULL,
    path VARCHAR(128) NOT NULL,

    module_id TEXT,
    FOREIGN KEY (module_id) REFERENCES module(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS pack (
    id TEXT PRIMARY KEY,

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

CREATE TABLE IF NOT EXISTS ownership (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    player VARCHAR(128) NOT NULL,
    trusted VARCHAR(128) NOT NULL,
    assistant VARCHAR(128) NOT NULL,

    pack_id TEXT,
    FOREIGN KEY (pack_id) REFERENCES pack(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS index_ (
    id TEXT PRIMARY KEY,

    folder VARCHAR(128) NOT NULL,
    img VARCHAR(128) NOT NULL,
    name VARCHAR(128) NOT NULL,
    type VARCHAR(128) NOT NULL
);

CREATE TABLE IF NOT EXISTS pack_to_index (
    pack_id TEXT,
    index_id TEXT,
    
    PRIMARY KEY (pack_id, index_id),
    FOREIGN KEY (index_id) REFERENCES index_(id) ON DELETE CASCADE,
    FOREIGN KEY (pack_id) REFERENCES pack(id) ON DELETE CASCADE
);

CREATE INDEX idx_pack_to_index_index_id ON pack_to_index(index_id);

CREATE TRIGGER tr_delete_pack_to_index_index
AFTER DELETE ON pack_to_index
FOR EACH ROW
BEGIN
    DELETE FROM index_ 
    WHERE id = OLD.index_id 
    AND NOT EXISTS (SELECT 1 FROM pack_to_index WHERE index_id = OLD.index_id);
END;

CREATE TABLE IF NOT EXISTS pack_folder (
    id TEXT PRIMARY KEY,

    description VARCHAR(128) NOT NULL,
    name VARCHAR(128) NOT NULL,
    sorting VARCHAR(128) NOT NULL,
    type VARCHAR(128) NOT NULL,
    sort INTEGER NOT NULL,

    pack_id TEXT,
    FOREIGN KEY (pack_id) REFERENCES pack(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS folder (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    name VARCHAR(128) NOT NULL,
    sorting VARCHAR(128) NOT NULL,
    color VARCHAR(128) NOT NULL,

    folder_id INTEGER,
    FOREIGN KEY (folder_id) REFERENCES folder(id) ON DELETE CASCADE
);

ALTER TABLE folder ADD COLUMN module_id TEXT REFERENCES module(id) ON DELETE CASCADE;
CREATE UNIQUE INDEX idx_folder_module_id ON folder(module_id);

CREATE TABLE IF NOT EXISTS folder_packs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    value VARCHAR(128) NOT NULL,

    folder_id INTEGER,
    FOREIGN KEY (folder_id) REFERENCES folder(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS news (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    title VARCHAR(64) NOT NULL,
    caption VARCHAR(64) NOT NULL,
    url VARCHAR(64) NOT NULL,
    image VARCHAR(64) NOT NULL,

    setup_id INTEGER,
    FOREIGN KEY (setup_id) REFERENCES setup(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS package_warnings (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    key_ TEXT NOT NULL,

    setup_id INTEGER,
    FOREIGN KEY (setup_id) REFERENCES setup(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS package_warnings_to_data (
    package_warnings_id INTEGER,
    package_warnings_data_id TEXT,
    
    PRIMARY KEY (package_warnings_id, package_warnings_data_id),
    FOREIGN KEY (package_warnings_id) REFERENCES package_warnings(id) ON DELETE CASCADE,
    FOREIGN KEY (package_warnings_data_id) REFERENCES package_warnings_data(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS package_warnings_data (
    id TEXT PRIMARY KEY,

    type VARCHAR(128) NOT NULL,
    manifest VARCHAR(128) NOT NULL,
    reinstallable BOOLEAN NOT NULL
);

CREATE INDEX idx_package_warnings_to_data_package_warnings_data_id ON package_warnings_to_data(package_warnings_data_id);

CREATE TRIGGER tr_delete_package_warnings_to_data_package_warnings_data
AFTER DELETE ON package_warnings_to_data
FOR EACH ROW
BEGIN
    DELETE FROM package_warnings_data 
    WHERE id = OLD.package_warnings_data_id 
    AND NOT EXISTS (SELECT 1 FROM package_warnings_to_data WHERE package_warnings_data_id = OLD.package_warnings_data_id);
END;

CREATE TABLE IF NOT EXISTS package_warnings_data_warning (
    id INTEGER PRIMARY KEY,

    value VARCHAR(128) NOT NULL,

    package_warnings_data_id TEXT,
    FOREIGN KEY (package_warnings_data_id) REFERENCES package_warnings_data(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS package_warnings_data_error (
    id INTEGER PRIMARY KEY,

    value VARCHAR(128) NOT NULL,

    package_warnings_data_id TEXT,
    FOREIGN KEY (package_warnings_data_id) REFERENCES package_warnings_data(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS system (
    id TEXT PRIMARY KEY,
    
    title VARCHAR(128) NOT NULL,
    description VARCHAR(128) NOT NULL,
    url VARCHAR(128) NOT NULL,
    license VARCHAR(128) NOT NULL,
    bugs VARCHAR(128) NOT NULL,
    changelog VARCHAR(128) NOT NULL,
    version VARCHAR(128) NOT NULL,
    manifest VARCHAR(128) NOT NULL,
    download VARCHAR(128) NOT NULL,
    background VARCHAR(128) NOT NULL DEFAULT FALSE,
    primary_token_attribute VARCHAR(128) NOT NULL DEFAULT FALSE,
    availability INTEGER NOT NULL DEFAULT FALSE,
    socket BOOLEAN NOT NULL DEFAULT FALSE,
    protected BOOLEAN NOT NULL DEFAULT FALSE,
    exclusive_ BOOLEAN NOT NULL DEFAULT FALSE,
    persistent_storage BOOLEAN NOT NULL DEFAULT FALSE,
    locked BOOLEAN NOT NULL DEFAULT FALSE,
    owned BOOLEAN NOT NULL DEFAULT FALSE,
    has_storage BOOLEAN NOT NULL DEFAULT FALSE,

    setup_id INTEGER,
    FOREIGN KEY (setup_id) REFERENCES setup(id) ON DELETE CASCADE
);

ALTER TABLE compatibility ADD COLUMN system_id TEXT REFERENCES system(id) ON DELETE CASCADE;
CREATE UNIQUE INDEX idx_compatibility_system_id ON compatibility(system_id);
ALTER TABLE relationships ADD COLUMN system_id TEXT REFERENCES system(id) ON DELETE CASCADE; 
CREATE UNIQUE INDEX idx_relationships_system_id ON relationships(system_id);
ALTER TABLE document_types ADD COLUMN system_id TEXT REFERENCES system(id) ON DELETE CASCADE;
CREATE UNIQUE INDEX idx_document_types_system_id ON document_types(system_id);

CREATE TABLE IF NOT EXISTS grid (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    type INTEGER NOT NULL,
    size INTEGER NOT NULL,
    distance INTEGER NOT NULL,
    diagonals INTEGER NOT NULL,
    thickness INTEGER NOT NULL,
    alpha REAL NOT NULL,
    color VARCHAR(128) NOT NULL,
    units VARCHAR(128) NOT NULL,
    style VARCHAR(128) NOT NULL,

    system_id TEXT UNIQUE,
    FOREIGN KEY (system_id) REFERENCES system(id) ON DELETE CASCADE
);

ALTER TABLE es_modules ADD COLUMN system_id TEXT REFERENCES system(id) ON DELETE CASCADE;
ALTER TABLE scripts ADD COLUMN system_id TEXT REFERENCES system(id) ON DELETE CASCADE; 
ALTER TABLE tags ADD COLUMN system_id TEXT REFERENCES system(id) ON DELETE CASCADE;
ALTER TABLE author ADD COLUMN system_id TEXT REFERENCES system(id) ON DELETE CASCADE;
ALTER TABLE media ADD COLUMN system_id TEXT REFERENCES system(id) ON DELETE CASCADE;
ALTER TABLE pack ADD COLUMN system_id TEXT REFERENCES system(id) ON DELETE CASCADE;
ALTER TABLE style ADD COLUMN system_id TEXT REFERENCES system(id) ON DELETE CASCADE;
ALTER TABLE language ADD COLUMN system_id TEXT REFERENCES system(id) ON DELETE CASCADE;
ALTER TABLE folder ADD COLUMN system_id TEXT REFERENCES system(id) ON DELETE CASCADE;

CREATE TABLE IF NOT EXISTS world (
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

ALTER TABLE compatibility ADD COLUMN world_id TEXT REFERENCES world(id) ON DELETE CASCADE;
CREATE UNIQUE INDEX idx_compatibility_world_id ON compatibility(world_id);
ALTER TABLE relationships ADD COLUMN world_id TEXT REFERENCES world(id) ON DELETE CASCADE;
CREATE UNIQUE INDEX idx_relationships_world_id ON relationships(world_id);

ALTER TABLE tags ADD COLUMN world_id TEXT REFERENCES world(id) ON DELETE CASCADE;
ALTER TABLE scripts ADD COLUMN world_id TEXT REFERENCES world(id) ON DELETE CASCADE;
ALTER TABLE es_modules ADD COLUMN world_id TEXT REFERENCES world(id) ON DELETE CASCADE;
ALTER TABLE author ADD COLUMN world_id TEXT REFERENCES world(id) ON DELETE CASCADE;
ALTER TABLE media ADD COLUMN world_id TEXT REFERENCES world(id) ON DELETE CASCADE;
ALTER TABLE style ADD COLUMN world_id TEXT REFERENCES world(id) ON DELETE CASCADE;
ALTER TABLE language ADD COLUMN world_id TEXT REFERENCES world(id) ON DELETE CASCADE;
ALTER TABLE pack ADD COLUMN world_id TEXT REFERENCES world(id) ON DELETE CASCADE;
ALTER TABLE folder ADD COLUMN world_id TEXT REFERENCES world(id) ON DELETE CASCADE;
