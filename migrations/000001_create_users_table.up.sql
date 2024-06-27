CREATE TABLE IF NOT EXISTS `users` (
    `id` VARCHAR(36) PRIMARY KEY,
    `username` VARCHAR(64) NOT NULL,
    `email` VARCHAR(255) UNIQUE NOT NULL,
    `password_hash` VARCHAR(255) NOT NULL
);
