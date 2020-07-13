CREATE TABLE IF NOT EXISTS `bank_accounts` (
  `entity_id` CHAR(36) NOT NULL,
  `account_name` VARCHAR(255) NOT NULL,
  `bank_name` VARCHAR(255) NOT NULL,
  `account_holder_name` VARCHAR(255) NOT NULL,
  `account_number` VARCHAR(255) NOT NULL,
  `last_balance` DECIMAL(18,2) NOT NULL,
  `last_balance_date` TIMESTAMP NOT NULL,
  `status` ENUM('active', 'inactive') NOT NULL DEFAULT 'active',
  `created` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created_by` CHAR(36) NOT NULL,
  `updated` TIMESTAMP NULL DEFAULT NULL,
  `updated_by` CHAR(36) NULL DEFAULT NULL,
  `deleted` TIMESTAMP NULL DEFAULT NULL,
  `deleted_by` CHAR(36) NULL DEFAULT NULL,
  PRIMARY KEY (`entity_id`),
  UNIQUE KEY `bank_account_idx_1` (`bank_name`, `account_number`),
  INDEX `bank_account_idx_2` (`bank_name`),
  INDEX `bank_account_idx_3` (`account_holder_name`),
  INDEX `bank_account_idx_4` (`account_number`),
  INDEX `bank_account_idx_5` (`last_balance`),
  INDEX `bank_account_idx_6` (`last_balance_date`),
  INDEX `bank_account_idx_7` (`status`),
  INDEX `bank_account_idx_8` (`created`),
  INDEX `bank_account_idx_9` (`created_by`)
) ENGINE=InnoDB
DEFAULT CHARSET=utf8;

INSERT INTO `bank_accounts` 
(`entity_id`, `account_name`, `bank_name`, `account_holder_name`, `account_number`, `last_balance`, `last_balance_date`, `status`, `created_by`)
VALUES
('d631d2f6-c051-4ba7-9584-5665d03cda6e', 'Current Account', 'First National Bank', 'Admin User', '5647382910', 1000000, DATE_SUB(NOW(), INTERVAL 30 DAY), 'active', 'cf4dcb72-27ed-442f-bac0-2c2871e29b1c'),
('2d17f0c4-3ed9-4b50-84be-3cf9bc94cfaf', 'Savings Account', 'Second National Bank', 'Admin User', '1029384756', 500000, DATE_SUB(NOW(), INTERVAL 3 DAY), 'active', 'cf4dcb72-27ed-442f-bac0-2c2871e29b1c'),
('277a1795-d9e8-42f0-b343-a4658bd64ba4', 'Retirement Account', 'Third National Bank', 'Admin User', '1627384950', 2000000, DATE_SUB(NOW(), INTERVAL 1 DAY), 'active', 'cf4dcb72-27ed-442f-bac0-2c2871e29b1c');


CREATE TABLE IF NOT EXISTS `bank_account_balances` (
  `entity_id` CHAR(36) NOT NULL,
  `bank_account_entity_id` CHAR(36) NOT NULL,
  `date` TIMESTAMP NOT NULL,
  `balance` DECIMAL(18,2) NOT NULL, 
  `created` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created_by` CHAR(36) NOT NULL,
  `updated` TIMESTAMP NULL DEFAULT NULL,
  `updated_by` CHAR(36) NULL DEFAULT NULL,
  `deleted` TIMESTAMP NULL DEFAULT NULL,
  `deleted_by` CHAR(36) NULL DEFAULT NULL,
  PRIMARY KEY (`entity_id`),
  CONSTRAINT `fk_bab_bank_account_entity_id` FOREIGN KEY (`bank_account_entity_id`)
    REFERENCES `bank_accounts`(`entity_id`)
    ON UPDATE NO ACTION
    ON DELETE NO ACTION,
  INDEX `bank_account_balances_idx_1` (`date`),
  INDEX `bank_account_balances_idx_2` (`created`),
  INDEX `bank_account_balances_idx_3` (`created_by`)
) ENGINE=InnoDB
DEFAULT CHARSET=utf8;

INSERT INTO `bank_account_balances` (`entity_id`, `bank_account_entity_id`, `date`, `balance`, `created_by`)
VALUES
-- first account
('14bac85d-57fc-41a5-9e58-6119151d0b10', 'd631d2f6-c051-4ba7-9584-5665d03cda6e', DATE_SUB(NOW(), INTERVAL 360 DAY), 220000, 'cf4dcb72-27ed-442f-bac0-2c2871e29b1c'),
('acc05778-dca1-4dfb-980f-2fddab7fd177', 'd631d2f6-c051-4ba7-9584-5665d03cda6e', DATE_SUB(NOW(), INTERVAL 330 DAY), 253000, 'cf4dcb72-27ed-442f-bac0-2c2871e29b1c'),
('92b53ef9-b49b-4cd5-aeda-1de44268a8d2', 'd631d2f6-c051-4ba7-9584-5665d03cda6e', DATE_SUB(NOW(), INTERVAL 300 DAY), 254000, 'cf4dcb72-27ed-442f-bac0-2c2871e29b1c'),
('7296caac-2dfb-453a-857e-963f98b8d7f2', 'd631d2f6-c051-4ba7-9584-5665d03cda6e', DATE_SUB(NOW(), INTERVAL 270 DAY), 255000, 'cf4dcb72-27ed-442f-bac0-2c2871e29b1c'),
('f225dfa6-c53b-4f6c-8fd3-048829233025', 'd631d2f6-c051-4ba7-9584-5665d03cda6e', DATE_SUB(NOW(), INTERVAL 240 DAY), 260000, 'cf4dcb72-27ed-442f-bac0-2c2871e29b1c'),
('70066bf2-7a0f-4874-a38f-c5e7a6e39aaa', 'd631d2f6-c051-4ba7-9584-5665d03cda6e', DATE_SUB(NOW(), INTERVAL 210 DAY), 265000, 'cf4dcb72-27ed-442f-bac0-2c2871e29b1c'),
('b27fde47-7b6a-4894-b874-a71261b1986c', 'd631d2f6-c051-4ba7-9584-5665d03cda6e', DATE_SUB(NOW(), INTERVAL 180 DAY), 275000, 'cf4dcb72-27ed-442f-bac0-2c2871e29b1c'),
('9f76e502-b469-4563-b581-13599c858e1e', 'd631d2f6-c051-4ba7-9584-5665d03cda6e', DATE_SUB(NOW(), INTERVAL 150 DAY), 300000, 'cf4dcb72-27ed-442f-bac0-2c2871e29b1c'),
('a8f61305-473f-4e34-a0b6-d7a9376dc190', 'd631d2f6-c051-4ba7-9584-5665d03cda6e', DATE_SUB(NOW(), INTERVAL 120 DAY), 400000, 'cf4dcb72-27ed-442f-bac0-2c2871e29b1c'),
('61abc088-8453-409b-8660-847dc37fd307', 'd631d2f6-c051-4ba7-9584-5665d03cda6e', DATE_SUB(NOW(), INTERVAL 90 DAY), 500000, 'cf4dcb72-27ed-442f-bac0-2c2871e29b1c'),
('9112152c-21a9-4157-bd8d-e2cc7eddf00e', 'd631d2f6-c051-4ba7-9584-5665d03cda6e', DATE_SUB(NOW(), INTERVAL 60 DAY), 800000, 'cf4dcb72-27ed-442f-bac0-2c2871e29b1c'),
('177eb2a7-ca35-4b7d-a32c-3d0e1b3d103a', 'd631d2f6-c051-4ba7-9584-5665d03cda6e', DATE_SUB(NOW(), INTERVAL 30 DAY), 950000, 'cf4dcb72-27ed-442f-bac0-2c2871e29b1c'),
-- second account
('168316ad-b953-49e3-a027-e4301c0885e8', '2d17f0c4-3ed9-4b50-84be-3cf9bc94cfaf', DATE_SUB(NOW(), INTERVAL 31 DAY), 312500, 'cf4dcb72-27ed-442f-bac0-2c2871e29b1c'),
('6d54e995-b72d-48a5-9dcc-fc6988570e68', '2d17f0c4-3ed9-4b50-84be-3cf9bc94cfaf', DATE_SUB(NOW(), INTERVAL 24 DAY), 325000, 'cf4dcb72-27ed-442f-bac0-2c2871e29b1c'),
('16cffd41-8f32-41e5-b456-ab0617909147', '2d17f0c4-3ed9-4b50-84be-3cf9bc94cfaf', DATE_SUB(NOW(), INTERVAL 17 DAY), 350000, 'cf4dcb72-27ed-442f-bac0-2c2871e29b1c'),
('76c1b005-2971-4a9d-9b92-9df492aa2efa', '2d17f0c4-3ed9-4b50-84be-3cf9bc94cfaf', DATE_SUB(NOW(), INTERVAL 10 DAY), 400000, 'cf4dcb72-27ed-442f-bac0-2c2871e29b1c'),
('0b877c4f-dbd0-4942-b655-5aadc5a6958b', '2d17f0c4-3ed9-4b50-84be-3cf9bc94cfaf', DATE_SUB(NOW(), INTERVAL 3 DAY), 500000, 'cf4dcb72-27ed-442f-bac0-2c2871e29b1c'),
-- third account
('269a633b-2752-47db-a8b7-845eaf41637f', '277a1795-d9e8-42f0-b343-a4658bd64ba4', DATE_SUB(NOW(), INTERVAL 720 DAY), 700000, 'cf4dcb72-27ed-442f-bac0-2c2871e29b1c'),
('c3fd1691-fe54-4261-9485-328357f312c9', '277a1795-d9e8-42f0-b343-a4658bd64ba4', DATE_SUB(NOW(), INTERVAL 360 DAY), 1000000, 'cf4dcb72-27ed-442f-bac0-2c2871e29b1c'),
('035f2e92-8036-4935-bd8e-9808fecd67fb', '277a1795-d9e8-42f0-b343-a4658bd64ba4', DATE_SUB(NOW(), INTERVAL 1 DAY), 2000000, 'cf4dcb72-27ed-442f-bac0-2c2871e29b1c');
