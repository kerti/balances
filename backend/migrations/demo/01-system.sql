-- This SQL script resets the database and puts it back to
-- "demo-ready" state, particularly in the system category.

DELETE FROM `users`;

INSERT INTO `users` (`entity_id`, `username`, `email`, `password`, `name`, `created_by`)
VALUES ('cf4dcb72-27ed-442f-bac0-2c2871e29b1c', 'admin', 'admin@example.com', '$2a$10$blyljuKpH9.TBMeaTHMpv.O5kmoqJlE5VfMcUdwlUeuMbj5ZEsOVq', 'Admin User', 'cf4dcb72-27ed-442f-bac0-2c2871e29b1c');