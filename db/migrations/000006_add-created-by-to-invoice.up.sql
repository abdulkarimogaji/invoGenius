ALTER TABLE `invoice` ADD `created_by` INT NOT NULL;
ALTER TABLE `invoice` ADD FOREIGN KEY (`created_by`) REFERENCES `user` (`id`);