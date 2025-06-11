-- This SQL script resets the database and puts it back to
-- "demo-ready" state, particularly in the assets category.

-- -- Bank Accounts and Bank Acount Balances

DELETE FROM `bank_account_balances`;
DELETE FROM `bank_accounts`;

INSERT INTO `bank_accounts` 
(`entity_id`, `account_name`, `bank_name`, `account_holder_name`, `account_number`, `last_balance`, `last_balance_date`, `status`, `created_by`)
VALUES
('d631d2f6-c051-4ba7-9584-5665d03cda6e', 'Current Account', 'First National Bank', 'Admin User', '5647382910', 2500000, timestamp(last_day(curdate() - interval 1 month) + interval 1 day - interval 7 hour - interval 1 minute), 'active', 'cf4dcb72-27ed-442f-bac0-2c2871e29b1c'),
('2d17f0c4-3ed9-4b50-84be-3cf9bc94cfaf', 'Savings Account', 'Second National Bank', 'Admin User', '1029384756', 500000, timestamp(last_day(curdate() - interval 1 month) + interval 1 day - interval 7 hour - interval 1 minute), 'active', 'cf4dcb72-27ed-442f-bac0-2c2871e29b1c'),
('277a1795-d9e8-42f0-b343-a4658bd64ba4', 'Shopping Account', 'Third National Bank', 'Admin User', '1627384950', 2000000, timestamp(last_day(curdate() - interval 1 month) + interval 1 day - interval 7 hour - interval 1 minute), 'active', 'cf4dcb72-27ed-442f-bac0-2c2871e29b1c');

INSERT INTO `bank_account_balances` (`entity_id`, `bank_account_entity_id`, `date`, `balance`, `created_by`)
VALUES
-- first account
('14bac85d-57fc-41a5-9e58-6119151d0b10', 'd631d2f6-c051-4ba7-9584-5665d03cda6e', timestamp(last_day(curdate() - interval 12 month) + interval 1 day - interval 7 hour - interval 1 minute), RAND()*(2500000-20)+20, 'cf4dcb72-27ed-442f-bac0-2c2871e29b1c'),
('acc05778-dca1-4dfb-980f-2fddab7fd177', 'd631d2f6-c051-4ba7-9584-5665d03cda6e', timestamp(last_day(curdate() - interval 11 month) + interval 1 day - interval 7 hour - interval 1 minute), RAND()*(2500000-20)+20, 'cf4dcb72-27ed-442f-bac0-2c2871e29b1c'),
('92b53ef9-b49b-4cd5-aeda-1de44268a8d2', 'd631d2f6-c051-4ba7-9584-5665d03cda6e', timestamp(last_day(curdate() - interval 10 month) + interval 1 day - interval 7 hour - interval 1 minute), RAND()*(2500000-20)+20, 'cf4dcb72-27ed-442f-bac0-2c2871e29b1c'),
('7296caac-2dfb-453a-857e-963f98b8d7f2', 'd631d2f6-c051-4ba7-9584-5665d03cda6e', timestamp(last_day(curdate() - interval 9 month) + interval 1 day - interval 7 hour - interval 1 minute), RAND()*(2500000-20)+20, 'cf4dcb72-27ed-442f-bac0-2c2871e29b1c'),
('f225dfa6-c53b-4f6c-8fd3-048829233025', 'd631d2f6-c051-4ba7-9584-5665d03cda6e', timestamp(last_day(curdate() - interval 8 month) + interval 1 day - interval 7 hour - interval 1 minute), RAND()*(2500000-20)+20, 'cf4dcb72-27ed-442f-bac0-2c2871e29b1c'),
('70066bf2-7a0f-4874-a38f-c5e7a6e39aaa', 'd631d2f6-c051-4ba7-9584-5665d03cda6e', timestamp(last_day(curdate() - interval 7 month) + interval 1 day - interval 7 hour - interval 1 minute), RAND()*(2500000-20)+20, 'cf4dcb72-27ed-442f-bac0-2c2871e29b1c'),
('b27fde47-7b6a-4894-b874-a71261b1986c', 'd631d2f6-c051-4ba7-9584-5665d03cda6e', timestamp(last_day(curdate() - interval 6 month) + interval 1 day - interval 7 hour - interval 1 minute), RAND()*(2500000-20)+20, 'cf4dcb72-27ed-442f-bac0-2c2871e29b1c'),
('9f76e502-b469-4563-b581-13599c858e1e', 'd631d2f6-c051-4ba7-9584-5665d03cda6e', timestamp(last_day(curdate() - interval 5 month) + interval 1 day - interval 7 hour - interval 1 minute), RAND()*(2500000-20)+20, 'cf4dcb72-27ed-442f-bac0-2c2871e29b1c'),
('a8f61305-473f-4e34-a0b6-d7a9376dc190', 'd631d2f6-c051-4ba7-9584-5665d03cda6e', timestamp(last_day(curdate() - interval 4 month) + interval 1 day - interval 7 hour - interval 1 minute), RAND()*(2500000-20)+20, 'cf4dcb72-27ed-442f-bac0-2c2871e29b1c'),
('61abc088-8453-409b-8660-847dc37fd307', 'd631d2f6-c051-4ba7-9584-5665d03cda6e', timestamp(last_day(curdate() - interval 3 month) + interval 1 day - interval 7 hour - interval 1 minute), RAND()*(2500000-20)+20, 'cf4dcb72-27ed-442f-bac0-2c2871e29b1c'),
('9112152c-21a9-4157-bd8d-e2cc7eddf00e', 'd631d2f6-c051-4ba7-9584-5665d03cda6e', timestamp(last_day(curdate() - interval 2 month) + interval 1 day - interval 7 hour - interval 1 minute), RAND()*(2500000-20)+20, 'cf4dcb72-27ed-442f-bac0-2c2871e29b1c'),
('177eb2a7-ca35-4b7d-a32c-3d0e1b3d103a', 'd631d2f6-c051-4ba7-9584-5665d03cda6e', timestamp(last_day(curdate() - interval 1 month) + interval 1 day - interval 7 hour - interval 1 minute), 2500000, 'cf4dcb72-27ed-442f-bac0-2c2871e29b1c'),
-- second account
('1173a06c-0d53-4be7-bc6b-7f6a19536c71', '2d17f0c4-3ed9-4b50-84be-3cf9bc94cfaf', timestamp(last_day(curdate() - interval 12 month) + interval 1 day - interval 7 hour - interval 1 minute), RAND()*(2500000-20)+20, 'cf4dcb72-27ed-442f-bac0-2c2871e29b1c'),
('f2436052-9337-4e52-8b35-8806660842a1', '2d17f0c4-3ed9-4b50-84be-3cf9bc94cfaf', timestamp(last_day(curdate() - interval 11 month) + interval 1 day - interval 7 hour - interval 1 minute), RAND()*(2500000-20)+20, 'cf4dcb72-27ed-442f-bac0-2c2871e29b1c'),
('773c96f4-1bc9-411f-b7b5-fdd171e81e00', '2d17f0c4-3ed9-4b50-84be-3cf9bc94cfaf', timestamp(last_day(curdate() - interval 10 month) + interval 1 day - interval 7 hour - interval 1 minute), RAND()*(2500000-20)+20, 'cf4dcb72-27ed-442f-bac0-2c2871e29b1c'),
('b97b00d5-ee71-4348-821f-614eb0c2f36c', '2d17f0c4-3ed9-4b50-84be-3cf9bc94cfaf', timestamp(last_day(curdate() - interval 9 month) + interval 1 day - interval 7 hour - interval 1 minute), RAND()*(2500000-20)+20, 'cf4dcb72-27ed-442f-bac0-2c2871e29b1c'),
('c25e6a84-f059-48e6-abd1-a9e24b1ac2cb', '2d17f0c4-3ed9-4b50-84be-3cf9bc94cfaf', timestamp(last_day(curdate() - interval 8 month) + interval 1 day - interval 7 hour - interval 1 minute), RAND()*(2500000-20)+20, 'cf4dcb72-27ed-442f-bac0-2c2871e29b1c'),
('325f47c1-1066-439b-85d9-c97c249674a6', '2d17f0c4-3ed9-4b50-84be-3cf9bc94cfaf', timestamp(last_day(curdate() - interval 7 month) + interval 1 day - interval 7 hour - interval 1 minute), RAND()*(2500000-20)+20, 'cf4dcb72-27ed-442f-bac0-2c2871e29b1c'),
('ad7c0a73-a333-445d-a000-bf6e99a93f94', '2d17f0c4-3ed9-4b50-84be-3cf9bc94cfaf', timestamp(last_day(curdate() - interval 6 month) + interval 1 day - interval 7 hour - interval 1 minute), RAND()*(2500000-20)+20, 'cf4dcb72-27ed-442f-bac0-2c2871e29b1c'),
('613f234e-49b6-433f-9691-7d9335ad2a39', '2d17f0c4-3ed9-4b50-84be-3cf9bc94cfaf', timestamp(last_day(curdate() - interval 5 month) + interval 1 day - interval 7 hour - interval 1 minute), RAND()*(2500000-20)+20, 'cf4dcb72-27ed-442f-bac0-2c2871e29b1c'),
('bfe5d2af-b232-420c-a477-76e294f8bf44', '2d17f0c4-3ed9-4b50-84be-3cf9bc94cfaf', timestamp(last_day(curdate() - interval 4 month) + interval 1 day - interval 7 hour - interval 1 minute), RAND()*(2500000-20)+20, 'cf4dcb72-27ed-442f-bac0-2c2871e29b1c'),
('fe2730a0-9a5a-46bf-a165-0b414dfd11cc', '2d17f0c4-3ed9-4b50-84be-3cf9bc94cfaf', timestamp(last_day(curdate() - interval 3 month) + interval 1 day - interval 7 hour - interval 1 minute), RAND()*(2500000-20)+20, 'cf4dcb72-27ed-442f-bac0-2c2871e29b1c'),
('ebacd3ae-e39e-40f6-9c87-ebb1023b67a1', '2d17f0c4-3ed9-4b50-84be-3cf9bc94cfaf', timestamp(last_day(curdate() - interval 2 month) + interval 1 day - interval 7 hour - interval 1 minute), RAND()*(2500000-20)+20, 'cf4dcb72-27ed-442f-bac0-2c2871e29b1c'),
('5fd67e7b-1da8-49f3-b33a-de039abc1bf2', '2d17f0c4-3ed9-4b50-84be-3cf9bc94cfaf', timestamp(last_day(curdate() - interval 1 month) + interval 1 day - interval 7 hour - interval 1 minute), 500000, 'cf4dcb72-27ed-442f-bac0-2c2871e29b1c'),
-- third account
('ff9d8a1a-f2d2-4d22-b5ca-457d4f6bc6ef', '277a1795-d9e8-42f0-b343-a4658bd64ba4', timestamp(last_day(curdate() - interval 12 month) + interval 1 day - interval 7 hour - interval 1 minute), RAND()*(2500000-20)+20, 'cf4dcb72-27ed-442f-bac0-2c2871e29b1c'),
('aae357f4-c0a2-4026-abbf-f966e73b0e23', '277a1795-d9e8-42f0-b343-a4658bd64ba4', timestamp(last_day(curdate() - interval 11 month) + interval 1 day - interval 7 hour - interval 1 minute), RAND()*(2500000-20)+20, 'cf4dcb72-27ed-442f-bac0-2c2871e29b1c'),
('60c8e783-49a5-4372-8057-274f55ab489c', '277a1795-d9e8-42f0-b343-a4658bd64ba4', timestamp(last_day(curdate() - interval 10 month) + interval 1 day - interval 7 hour - interval 1 minute), RAND()*(2500000-20)+20, 'cf4dcb72-27ed-442f-bac0-2c2871e29b1c'),
('b2c5dc3a-535d-488a-84bb-d046c2b55a66', '277a1795-d9e8-42f0-b343-a4658bd64ba4', timestamp(last_day(curdate() - interval 9 month) + interval 1 day - interval 7 hour - interval 1 minute), RAND()*(2500000-20)+20, 'cf4dcb72-27ed-442f-bac0-2c2871e29b1c'),
('d5f14b2d-b954-4873-85eb-809675009c0a', '277a1795-d9e8-42f0-b343-a4658bd64ba4', timestamp(last_day(curdate() - interval 8 month) + interval 1 day - interval 7 hour - interval 1 minute), RAND()*(2500000-20)+20, 'cf4dcb72-27ed-442f-bac0-2c2871e29b1c'),
('719444f2-7ea1-4526-8a2f-771c588c5067', '277a1795-d9e8-42f0-b343-a4658bd64ba4', timestamp(last_day(curdate() - interval 7 month) + interval 1 day - interval 7 hour - interval 1 minute), RAND()*(2500000-20)+20, 'cf4dcb72-27ed-442f-bac0-2c2871e29b1c'),
('4ebeb79b-d488-46cc-ba9f-7a61333384c0', '277a1795-d9e8-42f0-b343-a4658bd64ba4', timestamp(last_day(curdate() - interval 6 month) + interval 1 day - interval 7 hour - interval 1 minute), RAND()*(2500000-20)+20, 'cf4dcb72-27ed-442f-bac0-2c2871e29b1c'),
('cdcd696f-39d6-4493-9d04-84011f2e99ce', '277a1795-d9e8-42f0-b343-a4658bd64ba4', timestamp(last_day(curdate() - interval 5 month) + interval 1 day - interval 7 hour - interval 1 minute), RAND()*(2500000-20)+20, 'cf4dcb72-27ed-442f-bac0-2c2871e29b1c'),
('bd1985d0-4405-4352-9717-9e42750421c7', '277a1795-d9e8-42f0-b343-a4658bd64ba4', timestamp(last_day(curdate() - interval 4 month) + interval 1 day - interval 7 hour - interval 1 minute), RAND()*(2500000-20)+20, 'cf4dcb72-27ed-442f-bac0-2c2871e29b1c'),
('74055928-5d27-484c-90e4-583ed3f9e977', '277a1795-d9e8-42f0-b343-a4658bd64ba4', timestamp(last_day(curdate() - interval 3 month) + interval 1 day - interval 7 hour - interval 1 minute), RAND()*(2500000-20)+20, 'cf4dcb72-27ed-442f-bac0-2c2871e29b1c'),
('a8ba1eec-38d4-42fb-848c-ba6474ba890b', '277a1795-d9e8-42f0-b343-a4658bd64ba4', timestamp(last_day(curdate() - interval 2 month) + interval 1 day - interval 7 hour - interval 1 minute), RAND()*(2500000-20)+20, 'cf4dcb72-27ed-442f-bac0-2c2871e29b1c'),
('ac17de77-34ef-47de-b11a-ab04ae65f7cd', '277a1795-d9e8-42f0-b343-a4658bd64ba4', timestamp(last_day(curdate() - interval 1 month) + interval 1 day - interval 7 hour - interval 1 minute), 2000000, 'cf4dcb72-27ed-442f-bac0-2c2871e29b1c');