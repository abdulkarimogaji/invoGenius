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
  `updated_at` datetime NOT NULL
);

CREATE TABLE `user` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `first_name` varchar(255) NOT NULL,
  `last_name` varchar(255) NOT NULL,
  `email` varchar(255) UNIQUE NOT NULL,
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
  `updated_at` datetime NOT NULL
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