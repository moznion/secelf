CREATE TABLE IF NOT EXISTS `file` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `file_name` VARCHAR(255),
    `encrypted_cek` BLOB
);

CREATE INDEX `idx_file_name` ON `file`(`file_name`);
