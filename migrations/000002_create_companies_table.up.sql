CREATE TABLE IF NOT EXISTS `companies` (
    `id` VARCHAR(36) PRIMARY KEY,
    `name` VARCHAR(15) UNIQUE NOT NULL,
    `description` VARCHAR(3000) DEFAULT '',
    `amount_of_employees` INT NOT NULL,
    `registered` BOOLEAN NOT NULL,
    `type` ENUM('Corporations', 'NonProfit', 'Cooperative', 'Sole Proprietorship') NOT NULL
);
