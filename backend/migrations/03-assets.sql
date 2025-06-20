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

CREATE TABLE IF NOT EXISTS `vehicles` (
  `entity_id` CHAR(36) NOT NULL,
  `name` VARCHAR(255) NOT NULL,
  `make` VARCHAR(255) NOT NULL,
  `model` VARCHAR(255) NOT NULL,
  `year` INT NOT NULL,
  `type` ENUM('car', 'truck', 'motorcycle', 'bicycle', 'other') NOT NULL,
  `title_holder` VARCHAR(255) NOT NULL,
  `license_plate_number` VARCHAR(255) NOT NULL,
  `purchase_date` TIMESTAMP NOT NULL,
  `initial_value` DECIMAL(18,2) NOT NULL,
  `initial_value_date` TIMESTAMP NOT NULL,
  `current_value` DECIMAL(18,2) NOT NULL,
  `current_value_date` TIMESTAMP NOT NULL,
  `annual_depreciation_percent` DECIMAL(12,4) NOT NULL,
  `status` ENUM('in_use', 'retired', 'sold') NOT NULL DEFAULT 'in_use',
  `created` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created_by` CHAR(36) NOT NULL,
  `updated` TIMESTAMP NULL DEFAULT NULL,
  `updated_by` CHAR(36) NULL DEFAULT NULL,
  `deleted` TIMESTAMP NULL DEFAULT NULL,
  `deleted_by` CHAR(36) NULL DEFAULT NULL,
  PRIMARY KEY (`entity_id`),
  UNIQUE KEY `vehicles_idx_1` (`make`, `model`, `year`, `license_plate_number`),
  INDEX `vehicles_idx_2` (`name`),
  INDEX `vehicles_idx_3` (`make`),
  INDEX `vehicles_idx_4` (`model`),
  INDEX `vehicles_idx_5` (`year`),
  INDEX `vehicles_idx_6` (`type`),
  INDEX `vehicles_idx_7` (`title_holder`),
  INDEX `vehicles_idx_8` (`license_plate_number`),
  INDEX `vehicles_idx_9` (`purchase_date`),
  INDEX `vehicles_idx_10` (`initial_value`),
  INDEX `vehicles_idx_11` (`initial_value_date`),
  INDEX `vehicles_idx_12` (`current_value`),
  INDEX `vehicles_idx_13` (`current_value_date`),
  INDEX `vehicles_idx_14` (`annual_depreciation_percent`),
  INDEX `vehicles_idx_15` (`status`),
  INDEX `vehicles_idx_16` (`created`),
  INDEX `vehicles_idx_17` (`created_by`)
) ENGINE=InnoDB
DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `vehicle_values` (
  `entity_id` CHAR(36) NOT NULL,
  `vehicle_entity_id` CHAR(36) NOT NULL,
  `date` TIMESTAMP NOT NULL,
  `value` DECIMAL(18,2) NOT NULL,
  `created` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created_by` CHAR(36) NOT NULL,
  `updated` TIMESTAMP NULL DEFAULT NULL,
  `updated_by` CHAR(36) NULL DEFAULT NULL,
  `deleted` TIMESTAMP NULL DEFAULT NULL,
  `deleted_by` CHAR(36) NULL DEFAULT NULL,
  PRIMARY KEY (`entity_id`),
  CONSTRAINT `fk_vv_vehicle_entity_id` FOREIGN KEY (`vehicle_entity_id`)
    REFERENCES `vehicles`(`entity_id`)
    ON UPDATE NO ACTION
    ON DELETE NO ACTION,
  INDEX `vehicle_values_idx_1` (`date`),
  INDEX `vehicle_values_idx_2` (`created`),
  INDEX `vehicle_values_idx_3` (`created_by`)
) ENGINE=InnoDB
DEFAULT CHARSET=utf8;