ALTER TABLE users
ALTER COLUMN avatar_path SET DEFAULT 'default_user.png';

ALTER TABLE rooms
ALTER COLUMN avatar_path SET DEFAULT 'default_room.png';
