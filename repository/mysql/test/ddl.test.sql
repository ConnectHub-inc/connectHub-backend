-- Description: テスト用のDDLを記述します
CREATE DATABASE IF NOT EXISTS `connecthubdb` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;

-- テーブル作成
USE `connecthubdb` ;

DROP TABLE IF EXISTS TestItems CASCADE;

CREATE TABLE TestItems (
    id CHAR(36) PRIMARY KEY, -- UUIDは36文字の文字列として格納されます
    user_id CHAR(36) NOT NULL,
    text TEXT NOT NULL,
    count INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);