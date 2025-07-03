-- This SQL script resets the database and puts it back to
-- "demo-ready" state, particularly in the system category.

START TRANSACTION;

-- -- Users

DELETE FROM `users`;

INSERT INTO `users` (`entity_id`, `username`, `email`, `password`, `name`, `created_by`)
VALUES ('cf4dcb72-27ed-442f-bac0-2c2871e29b1c', 'admin', 'johndoe@example.com', '$2a$10$blyljuKpH9.TBMeaTHMpv.O5kmoqJlE5VfMcUdwlUeuMbj5ZEsOVq', 'John Fitzgerald Doe', 'cf4dcb72-27ed-442f-bac0-2c2871e29b1c');

COMMIT;