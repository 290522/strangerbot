CREATE TABLE `matched_detail` (
                                  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
                                  `chat_id` bigint(20) NOT NULL,
                                  `match_chat_id` bigint(20) NOT NULL,
                                  `is_del` tinyint(4) unsigned NOT NULL DEFAULT '0' COMMENT 'is del 0.false 1.true',
                                  `deleted_time` int(10) unsigned NOT NULL DEFAULT '0' COMMENT 'deleted time',
                                  `create_time` int(10) unsigned NOT NULL DEFAULT '0' COMMENT 'create time',
                                  `modify_time` int(10) unsigned NOT NULL DEFAULT '0' COMMENT 'update time',
                                  PRIMARY KEY (`id`) USING BTREE,
                                  KEY `idx_chat_id` (`chat_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='matched detail table';


ALTER TABLE `users` ADD COLUMN `custom_rate_limit_enabled` tinyint(4) unsigned NOT NULL DEFAULT '0' COMMENT 'custom user match rate limit  0.false 1.true';
ALTER TABLE `users` ADD COLUMN `rate_limit_unit` VARCHAR(32) NOT NULL DEFAULT 'day' COMMENT 'rate limit unit only support:day';
ALTER TABLE `users` ADD COLUMN `rate_limit_unit_period` int(10) unsigned NOT NULL DEFAULT '7' COMMENT 'rate limit per unit period,dont change this value.';
ALTER TABLE `users` ADD COLUMN `match_per_rate` tinyint(4) unsigned NOT NULL DEFAULT '2' COMMENT 'match per rate';