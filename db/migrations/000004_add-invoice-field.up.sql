ALTER TABLE `invoice` ADD `deadline` DATETIME NOT NULL;
ALTER TABLE `invoice` ADD `invoice_file` TEXT;
CREATE TABLE `invoice_activity` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `user_id` int NOT NULL,
  `invoice_id` int NOT NULL,
  `action_type` varchar(255) NOT NULL,
  `resource_id` int NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `attachment` text
);

ALTER TABLE `invoice_activity` ADD FOREIGN KEY (`user_id`) REFERENCES `user` (`id`);

ALTER TABLE `invoice_activity` ADD FOREIGN KEY (`invoice_id`) REFERENCES `invoice` (`id`);