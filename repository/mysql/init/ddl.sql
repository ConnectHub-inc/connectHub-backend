CREATE DATABASE IF NOT EXISTS `connecthubdb` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE `connecthubdb` ;

DROP TABLE IF EXISTS Messages CASCADE;
DROP TABLE IF EXISTS User_Rooms CASCADE;
DROP TABLE IF EXISTS User_Workspaces CASCADE;
DROP TABLE IF EXISTS Users CASCADE;
DROP TABLE IF EXISTS ActionTags CASCADE;
DROP TABLE IF EXISTS Rooms CASCADE;
DROP TABLE IF EXISTS Workspaces CASCADE;

CREATE TABLE  Workspaces (
    id CHAR(36) PRIMARY KEY, -- UUIDは36文字の文字列として格納されます
    name VARCHAR(50) NOT NULL,
    description TEXT
);

CREATE TABLE Rooms (
    id CHAR(36) PRIMARY KEY, -- UUIDは36文字の文字列として格納されます
    workspace_id CHAR(36) NOT NULL,
    name VARCHAR(50) NOT NULL,
    private BOOLEAN NOT NULL,
    description TEXT,
    FOREIGN KEY (workspace_id) REFERENCES Workspaces(id) ON DELETE CASCADE
);

CREATE TABLE ActionTags (
    id CHAR(36) PRIMARY KEY, -- UUIDは36文字の文字列として格納されます
    name VARCHAR(50) NOT NULL,
    description TEXT
);

CREATE TABLE Users (
    id CHAR(36) PRIMARY KEY, -- UUIDは36文字の文字列として格納されます
    name VARCHAR(50) NOT NULL,
    email VARCHAR(150) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,  -- 暗号化されたパスワードを格納
    profile_image_url VARCHAR(255) NOT NULL
);

CREATE TABLE User_Workspaces (
    user_id CHAR(36) NOT NULL,
    workspace_id CHAR(36) NOT NULL,
    PRIMARY KEY (user_id, workspace_id),
    FOREIGN KEY (user_id) REFERENCES Users(id) ON DELETE CASCADE,
    FOREIGN KEY (workspace_id) REFERENCES Workspaces(id) ON DELETE CASCADE
);

CREATE TABLE User_Rooms (
    user_id CHAR(36) NOT NULL,
    room_id CHAR(36) NOT NULL,
    PRIMARY KEY (user_id, room_id),
    FOREIGN KEY (user_id) REFERENCES Users(id) ON DELETE CASCADE,
    FOREIGN KEY (room_id) REFERENCES Rooms(id) ON DELETE CASCADE
);

CREATE TABLE Messages (
    id CHAR(36) PRIMARY KEY, -- UUIDは36文字の文字列として格納されます
    user_id CHAR(36) NOT NULL,
    room_id CHAR(36) NOT NULL,
    action_tag_id  CHAR(36) NOT NULL,
    text TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES Users(id) ON DELETE CASCADE,
    FOREIGN KEY (room_id) REFERENCES Rooms(id) ON DELETE CASCADE,
    FOREIGN KEY (action_tag_id) REFERENCES ActionTags(id) ON DELETE CASCADE
);