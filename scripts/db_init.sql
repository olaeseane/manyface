DROP TABLE IF EXISTS connection;
DROP TABLE IF EXISTS face;
DROP TABLE IF EXISTS userV1beta1;
DROP TABLE IF EXISTS userV2beta1;
DROP TABLE IF EXISTS faceV2beta1;
DROP TABLE IF EXISTS session;
-- DROP TABLE IF EXISTS wordlist;
PRAGMA foreign_keys = ON;
CREATE TABLE IF NOT EXISTS userV1beta1 (
    user_id INTEGER PRIMARY KEY,
    username TEXT NOT NULL,
    password BLOB NOT NULL
);
CREATE TABLE IF NOT EXISTS userV2beta1 (
    user_id INTEGER PRIMARY KEY,
    seed TEXT NOT NULL,
    password BLOB NOT NULL
);
CREATE TABLE IF NOT EXISTS faceV2beta1 (
    face_id TEXT PRIMARY KEY,
    purpose TEXT NOT NULL,
    nick TEXT NOT NULL,
    bio TEXT,
    comments TEXT,
    user_id INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES user (user_id) ON UPDATE RESTRICT ON DELETE RESTRICT
) WITHOUT ROWID;
CREATE TABLE IF NOT EXISTS face (
    face_id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    user_id INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES user (user_id) ON UPDATE RESTRICT ON DELETE RESTRICT
) WITHOUT ROWID;
CREATE TABLE IF NOT EXISTS connection (
    conn_id INTEGER PRIMARY KEY,
    mtrx_user_id TEXT NOT NULL,
    mtrx_password TEXT NOT NULL,
    mtrx_access_token TEXT,
    mtrx_room_id TEXT NOT NULL,
    mtrx_peer_id TEXT NOT NULL,
    face_id TEXT NOT NULL,
    face_peer_id TEXT NOT NULL,
    FOREIGN KEY (face_id) REFERENCES face (face_id) ON UPDATE RESTRICT ON DELETE RESTRICT
);
CREATE TABLE IF NOT EXISTS session (
    sess_id TEXT PRIMARY KEY,
    user_id INTEGER NOT NULL
);
CREATE TABLE IF NOT EXISTS wordlist (word TEXT NOT NULL);
INSERT INTO userV1beta1 (username, password)
VALUES ('user1', 'welcome');
INSERT INTO userV1beta1 (username, password)
VALUES ('user2', 'welcome');
INSERT INTO face (face_id, name, description, user_id)
VALUES (
        'ada2da68aace03fa2891efbb9314f2c1',
        'Willy_Wanka',
        'for chatting on the internet',
        1
    );
INSERT INTO face (face_id, name, description, user_id)
VALUES (
        '32f0df4913547929dd69ed63e7b8cb3a',
        'Mr._Shaquille_Oatmeal',
        'my offical login',
        1
    );
INSERT INTO face (face_id, name, description, user_id)
VALUES (
        'be585b5d359a7fe4deef8466d940c6d3',
        'Nameless_Faceless',
        'bad guy',
        2
    );