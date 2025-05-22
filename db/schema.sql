CREATE TABLE `invoice` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `user_id` int NOT NULL,
  `amount` double NOT NULL,
  `vat` double NOT NULL,
  `type` varchar(255) NOT NULL,
  `issued_at` datetime NOT NULL,
  `from_date` date NOT NULL,
  `until_date` date NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `currency` varchar(255) NOT NULL,
  `deadline` datetime NOT NULL,
  `invoice_file` text
);

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

CREATE TABLE `user` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `first_name` varchar(255) NOT NULL,
  `last_name` varchar(255) NOT NULL,
  `email` varchar(255) UNIQUE NOT NULL,
  `phone` varchar(255) NOT NULL,
  `photo` text,
  `role` varchar(255) NOT NULL,
  `password` varchar(255),
  `status` varchar(255) NOT NULL DEFAULT 'inactive',
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL
);

CREATE TABLE `transaction` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `invoice_id` int NOT NULL,
  `payment_method` varchar(255),
  `paid_at` datetime,
  `status` varchar(255) NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `amount` double NOT NULL
);

CREATE TABLE `receipt` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `transaction_id` int NOT NULL,
  `uploaded_by` int,
  `filename` varchar(255),
  `file` text,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL
);

CREATE TABLE `setting` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `setting_key` varchar(255) UNIQUE NOT NULL,
  `setting_value` varchar(255) NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL
);

ALTER TABLE `invoice` ADD FOREIGN KEY (`user_id`) REFERENCES `user` (`id`);

ALTER TABLE `transaction` ADD FOREIGN KEY (`invoice_id`) REFERENCES `invoice` (`id`);

ALTER TABLE `receipt` ADD FOREIGN KEY (`transaction_id`) REFERENCES `transaction` (`id`);

ALTER TABLE `receipt` ADD FOREIGN KEY (`uploaded_by`) REFERENCES `user` (`id`);

ALTER TABLE `invoice_activity` ADD FOREIGN KEY (`user_id`) REFERENCES `user` (`id`);

ALTER TABLE `invoice_activity` ADD FOREIGN KEY (`invoice_id`) REFERENCES `invoice` (`id`);
