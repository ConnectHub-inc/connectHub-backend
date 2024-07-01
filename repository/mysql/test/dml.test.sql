-- dml.test.sql
USE `connecthubTestDB`;

-- ワークスペースデータを挿入
INSERT INTO Workspaces (id, name, description) VALUES
('5fe0e237-6b49-11ee-b686-0242c0a87001', 'Workspace 1', 'Description for Workspace 1'),
('5fe0e238-6b49-11ee-b686-0242c0a87001', 'Workspace 2', 'Description for Workspace 2');

-- ルームデータを挿入
INSERT INTO Rooms (id, workspace_id, name, private, description) VALUES
('5fe0e239-6b49-11ee-b686-0242c0a87001', '5fe0e237-6b49-11ee-b686-0242c0a87001', 'Room 1', false, 'Description for Room 1 in Workspace 1'),
('5fe0e23a-6b49-11ee-b686-0242c0a87001', '5fe0e237-6b49-11ee-b686-0242c0a87001', 'Room 2', true, 'Description for Room 2 in Workspace 1'),
('5fe0e23b-6b49-11ee-b686-0242c0a87001', '5fe0e238-6b49-11ee-b686-0242c0a87001', 'Room 3', false, 'Description for Room 3 in Workspace 2');

-- アクションタグデータを挿入
INSERT INTO ActionTags (id, name, description) VALUES
('5fe0e23c-6b49-11ee-b686-0242c0a87001', 'Tag 1', 'Description for Tag 1'),
('5fe0e23d-6b49-11ee-b686-0242c0a87001', 'Tag 2', 'Description for Tag 2');

-- ユーザーデータを挿入
INSERT INTO Users (id, email, password) VALUES
('5fe0e23e-6b49-11ee-b686-0242c0a87001', 'john.doe@example.com', 'hashed_password_1'),
('5fe0e23f-6b49-11ee-b686-0242c0a87001', 'jane.smith@example.com', 'hashed_password_2'),
('5fe0e240-6b49-11ee-b686-0242c0a87001', 'alice.johnson@example.com', 'hashed_password_3');

-- ユーザーとワークスペースの関係データを挿入
INSERT INTO User_Workspaces (user_id, workspace_id, name, profile_image_url, is_admin, is_deleted) VALUES
('5fe0e23e-6b49-11ee-b686-0242c0a87001', '5fe0e237-6b49-11ee-b686-0242c0a87001', 'John Doe', 'https://example.com/profile_image_1', false, false),
('5fe0e23f-6b49-11ee-b686-0242c0a87001', '5fe0e237-6b49-11ee-b686-0242c0a87001', 'Jane Smith', 'https://example.com/profile_image_2', false, false),
('5fe0e240-6b49-11ee-b686-0242c0a87001', '5fe0e238-6b49-11ee-b686-0242c0a87001', 'Alice Johnson', 'https://example.com/profile_image_3', false, false);

-- ユーザーとルームの関係データを挿入
INSERT INTO User_Rooms (user_workspace_id, room_id) VALUES
(CONCAT('5fe0e23e-6b49-11ee-b686-0242c0a87001', '_', '5fe0e237-6b49-11ee-b686-0242c0a87001'), '5fe0e239-6b49-11ee-b686-0242c0a87001'),
(CONCAT('5fe0e23f-6b49-11ee-b686-0242c0a87001', '_', '5fe0e237-6b49-11ee-b686-0242c0a87001'), '5fe0e23a-6b49-11ee-b686-0242c0a87001'),
(CONCAT('5fe0e240-6b49-11ee-b686-0242c0a87001', '_', '5fe0e238-6b49-11ee-b686-0242c0a87001'), '5fe0e23b-6b49-11ee-b686-0242c0a87001');

-- メッセージデータを挿入
INSERT INTO Messages (id, user_workspace_id, room_id, action_tag_id, text, created_at, updated_at) VALUES
('5fe0e241-6b49-11ee-b686-0242c0a87001', CONCAT('5fe0e23e-6b49-11ee-b686-0242c0a87001', '_', '5fe0e237-6b49-11ee-b686-0242c0a87001'), '5fe0e239-6b49-11ee-b686-0242c0a87001', '5fe0e23c-6b49-11ee-b686-0242c0a87001', 'Hello from user 1 in room 1', '2023-01-01 10:00:00', '2023-01-01 10:00:00'),
('5fe0e242-6b49-11ee-b686-0242c0a87001', CONCAT('5fe0e23f-6b49-11ee-b686-0242c0a87001', '_', '5fe0e237-6b49-11ee-b686-0242c0a87001'), '5fe0e23a-6b49-11ee-b686-0242c0a87001', '5fe0e23d-6b49-11ee-b686-0242c0a87001', 'Hello from user 2 in room 2', '2023-01-01 11:00:00', '2023-01-01 11:00:00'),
('5fe0e243-6b49-11ee-b686-0242c0a87001', CONCAT('5fe0e240-6b49-11ee-b686-0242c0a87001', '_', '5fe0e238-6b49-11ee-b686-0242c0a87001'), '5fe0e23b-6b49-11ee-b686-0242c0a87001', '5fe0e23c-6b49-11ee-b686-0242c0a87001', 'Hello from user 3 in room 3', '2023-01-01 12:00:00', '2023-01-01 12:00:00');
