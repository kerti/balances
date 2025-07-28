CREATE TABLE IF NOT EXISTS `personal_debts` (
  `entity_id` CHAR(36) NOT NULL,
  `name` VARCHAR(255) NOT NULL,
  `creditor` VARCHAR(255) NOT NULL,
  `contact_information` VARCHAR(255) NOT NULL,
  `principal` DECIMAL(18,2) NOT NULL,
  `interest` DECIMAL(18,2) NOT NULL,
  `interest_type` ENUM ('nominal', 'percentage') NOT NULL,
  `date` TIMESTAMP NOT NULL,
  `current_balance` DECIMAL(18,2) NOT NULL,
  `current_balance_date` TIMESTAMP NOT NULL,
  `status` ENUM ('active', 'paid', 'defaulted', 'written_off'),
  `created` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created_by` CHAR(36) NOT NULL,
  `updated` TIMESTAMP NULL DEFAULT NULL,
  `updated_by` CHAR(36) NULL DEFAULT NULL,
  `deleted` TIMESTAMP NULL DEFAULT NULL,
  `deleted_by` CHAR(36) NULL DEFAULT NULL,
  PRIMARY KEY (`entity_id`),
  UNIQUE KEY `personal_debt_idx_1` (`creditor`, `principal`, `interest`, `interest_type`, `date`),
  INDEX `personal_debt_idx_2` (`name`),
  INDEX `personal_debt_idx_3` (`creditor`),
  INDEX `personal_debt_idx_4` (`contact_information`),
  INDEX `personal_debt_idx_5` (`principal`),
  INDEX `personal_debt_idx_6` (`interest`),
  INDEX `personal_debt_idx_7` (`interest_type`),
  INDEX `personal_debt_idx_8` (`date`),
  INDEX `personal_debt_idx_9` (`current_balance`),
  INDEX `personal_debt_idx_10` (`current_balance_date`),
  INDEX `personal_debt_idx_11` (`status`),
  INDEX `personal_debt_idx_12` (`created`),
  INDEX `personal_debt_idx_13` (`created_by`)
) ENGINE=InnoDB
DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `personal_debt_balances` (
  `entity_id` CHAR(36) NOT NULL,
  `personal_debt_entity_id` CHAR(36) NOT NULL,
  `date` TIMESTAMP NOT NULL,
  `balance` DECIMAL(18,2) NOT NULL, 
  `created` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created_by` CHAR(36) NOT NULL,
  `updated` TIMESTAMP NULL DEFAULT NULL,
  `updated_by` CHAR(36) NULL DEFAULT NULL,
  `deleted` TIMESTAMP NULL DEFAULT NULL,
  `deleted_by` CHAR(36) NULL DEFAULT NULL,
  PRIMARY KEY (`entity_id`),
  CONSTRAINT `fk_pdb_personal_debt_entity_id` FOREIGN KEY (`personal_debt_entity_id`)
    REFERENCES `personal_debts`(`entity_id`)
    ON UPDATE NO ACTION
    ON DELETE NO ACTION,
  INDEX `personal_debt_balances_idx_1` (`date`),
  INDEX `personal_debt_balances_idx_2` (`created`),
  INDEX `personal_debt_balances_idx_3` (`created_by`)
) ENGINE=InnoDB
DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `personal_debt_payments` (
  `entity_id` CHAR(36) NOT NULL,
  `personal_debt_entity_id` CHAR(36) NOT NULL,
  `payment_date` TIMESTAMP NOT NULL,
  `payment_amount` DECIMAL(18,2) NOT NULL,
  `created` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created_by` CHAR(36) NOT NULL,
  `updated` TIMESTAMP NULL DEFAULT NULL,
  `updated_by` CHAR(36) NULL DEFAULT NULL,
  `deleted` TIMESTAMP NULL DEFAULT NULL,
  `deleted_by` CHAR(36) NULL DEFAULT NULL,
  PRIMARY KEY (`entity_id`),
  CONSTRAINT `fk_pdr_personal_debt_entity_id` FOREIGN KEY (`personal_debt_entity_id`)
    REFERENCES `personal_debts`(`entity_id`)
    ON UPDATE NO ACTION
    ON DELETE NO ACTION,
  INDEX `personal_debt_payments_idx_1` (`payment_date`),
  INDEX `personal_debt_payments_idx_2` (`payment_amount`),
  INDEX `personal_debt_payments_idx_3` (`created`),
  INDEX `personal_debt_payments_idx_4` (`created_by`)
) ENGINE=InnoDB
DEFAULT CHARSET=utf8;
