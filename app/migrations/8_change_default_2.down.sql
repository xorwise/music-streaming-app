ALTER TABLE users
ALTER COLUMN avatar_path SET DEFAULT 'media/default_user.png';

ALTER TABLE rooms
ALTER COLUMN avatar_path SET DEFAULT 'media/default_room.png';
