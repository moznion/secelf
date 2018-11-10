CREATE TABLE IF NOT EXISTS `file` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `file_name` VARCHAR(255),
    `salt` VARCHAR(32)
);

CREATE INDEX `idx_file_name` ON `file`(`file_name`);
