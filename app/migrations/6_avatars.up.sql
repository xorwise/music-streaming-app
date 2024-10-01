ALTER TABLE rooms
ADD COLUMN IF NOT EXISTS avatar_path varchar(255) NOT NULL DEFAULT 'default_room.jpg';

ALTER TABLE users
ADD COLUMN IF NOT EXISTS avatar_path varchar(255) NOT NULL DEFAULT 'default_user.jpg';
