CREATE TABLE IF NOT EXISTS `users` (
  `entity_id` CHAR(36) NOT NULL,
  `username` VARCHAR(512) NOT NULL,
  `email` VARCHAR(512) NOT NULL,
  `password` VARCHAR(1000) NOT NULL,
  `name` VARCHAR(512) NOT NULL,
  `created` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created_by` CHAR(36) NOT NULL,
  `updated` TIMESTAMP NULL DEFAULT NULL,
  `updated_by` CHAR(36) NULL DEFAULT NULL,
  PRIMARY KEY (`entity_id`)
) ENGINE=InnoDB
DEFAULT CHARSET=utf8;

INSERT INTO `users` (`entity_id`, `username`, `email`, `password`, `name`, `created_by`)
VALUES ('cf4dcb72-27ed-442f-bac0-2c2871e29b1c', 'admin', 'admin@example.com', '$2a$10$blyljuKpH9.TBMeaTHMpv.O5kmoqJlE5VfMcUdwlUeuMbj5ZEsOVq', 'Admin User', 'cf4dcb72-27ed-442f-bac0-2c2871e29b1c');
