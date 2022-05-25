CREATE TABLE `queue` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `data` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT '消息',
  `lock_id` char(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '锁id',
  `retry_times` int unsigned NOT NULL DEFAULT '0' COMMENT '重试次数',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='消息队列';