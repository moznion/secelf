CREATE TABLE IF NOT EXISTS `file` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `filename` VARCHAR(255),
    `encrypted_cek` BLOB
);

CREATE INDEX `idx_filename` ON `file`(`filename`);
