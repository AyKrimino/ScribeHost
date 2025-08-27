-- +goose Up

CREATE TABLE `users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `email` varchar(100) NOT NULL,
  `password_hash` varchar(255) NOT NULL,
  `role` varchar(50) DEFAULT 'author',
  `email_verified` boolean DEFAULT FALSE,
  `is_active` boolean DEFAULT TRUE,
  `otp_secret` varchar(100),
  `otp` varchar(64),
  `otp_expiry` datetime(3),
  `reset_token` varchar(100),
  `reset_token_expiry` datetime(3),
  `last_login` datetime(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uq_users_email` (`email`),
  KEY `idx_users_role` (`role`),
  KEY `idx_users_email_verified` (`email_verified`),
  KEY `idx_users_is_active` (`is_active`)
);

CREATE TABLE `profiles` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `user_id` bigint unsigned NOT NULL,
  `first_name` varchar(100),
  `last_name` varchar(100),
  `avatar_url` text,
  `phone` varchar(20),
  `address` text,
  `date_of_birth` datetime(3),
  `gender` varchar(20) DEFAULT 'male',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uq_profiles_user_id` (`user_id`),
  KEY `idx_profiles_phone` (`phone`),
  CONSTRAINT `fk_profiles_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE `refresh_tokens` (
  `token_hash` varchar(64) NOT NULL,
  `user_id` bigint unsigned NOT NULL,
  `expiry` datetime(3) NOT NULL,
  `issued_at` datetime(3) NOT NULL,
  `user_agent` text,
  `ip_address` varchar(45),
  `is_revoked` boolean DEFAULT FALSE,
  PRIMARY KEY (`token_hash`),
  KEY `idx_refresh_tokens_user_id` (`user_id`),
  KEY `idx_refresh_tokens_expiry` (`expiry`),
  KEY `idx_refresh_tokens_is_revoked` (`is_revoked`),
  CONSTRAINT `fk_refresh_tokens_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
);

-- +goose Down

DROP TABLE IF EXISTS `refresh_tokens`;
DROP TABLE IF EXISTS `profiles`;
DROP TABLE IF EXISTS `users`;
