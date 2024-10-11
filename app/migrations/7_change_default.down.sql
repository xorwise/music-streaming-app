ALTER TABLE users
ALTER COLUMN avatar_path SET DEFAULT 'default_user.jpg';

ALTER TABLE rooms
ALTER COLUMN avatar_path SET DEFAULT 'default_room.jpg';
