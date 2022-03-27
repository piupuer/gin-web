-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE IF NOT EXISTS `tb_sys_api` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'auto increment id',
  `created_at` datetime(3) DEFAULT NULL COMMENT 'create time',
  `updated_at` datetime(3) DEFAULT NULL COMMENT 'update time',
  `deleted_at` datetime(3) DEFAULT NULL COMMENT 'soft delete time',
  `method` longtext COLLATE utf8mb4_general_ci COMMENT 'request method',
  `path` longtext COLLATE utf8mb4_general_ci COMMENT 'api path',
  `category` longtext COLLATE utf8mb4_general_ci COMMENT 'api group category',
  `desc` longtext COLLATE utf8mb4_general_ci COMMENT 'api description',
  PRIMARY KEY (`id`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=72 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `tb_sys_casbin` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `ptype` varchar(100) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT 'enforer type',
  `v0` varchar(100) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT 'role keyword(SysRole.Keyword)',
  `v1` varchar(100) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT 'resource name',
  `v2` varchar(100) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT 'request method',
  `v3` varchar(100) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `v4` varchar(100) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `v5` varchar(100) COLLATE utf8mb4_general_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_casbin` (`ptype`,`v0`,`v1`,`v2`,`v3`,`v4`,`v5`),
  UNIQUE KEY `idx_tb_sys_casbin` (`ptype`,`v0`,`v1`,`v2`,`v3`,`v4`,`v5`)
) ENGINE=InnoDB AUTO_INCREMENT=94 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `tb_sys_dict` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'auto increment id',
  `created_at` datetime(3) DEFAULT NULL COMMENT 'create time',
  `updated_at` datetime(3) DEFAULT NULL COMMENT 'update time',
  `deleted_at` datetime(3) DEFAULT NULL COMMENT 'soft delete time',
  `name` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT 'name',
  `desc` longtext COLLATE utf8mb4_general_ci COMMENT 'description',
  `status` tinyint(1) DEFAULT '1' COMMENT 'status(0: disabled, 1: enabled)',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_name_unique` (`name`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `tb_sys_dict_data` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'auto increment id',
  `created_at` datetime(3) DEFAULT NULL COMMENT 'create time',
  `updated_at` datetime(3) DEFAULT NULL COMMENT 'update time',
  `deleted_at` datetime(3) DEFAULT NULL COMMENT 'soft delete time',
  `key` longtext COLLATE utf8mb4_general_ci COMMENT 'key',
  `val` longtext COLLATE utf8mb4_general_ci COMMENT 'val',
  `addition` longtext COLLATE utf8mb4_general_ci COMMENT 'custom addition params',
  `sort` bigint unsigned DEFAULT NULL COMMENT 'sort',
  `status` tinyint(1) DEFAULT '1' COMMENT 'status(0: disabled, 1: enabled)',
  `dict_id` bigint unsigned DEFAULT NULL COMMENT 'dict id',
  PRIMARY KEY (`id`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=14 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `tb_sys_user` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'auto increment id',
  `created_at` datetime(3) DEFAULT NULL COMMENT 'create time',
  `updated_at` datetime(3) DEFAULT NULL COMMENT 'update time',
  `deleted_at` datetime(3) DEFAULT NULL COMMENT 'soft delete time',
  `username` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT 'user login name',
  `password` longtext COLLATE utf8mb4_general_ci COMMENT 'password',
  `mobile` longtext COLLATE utf8mb4_general_ci COMMENT 'mobile number',
  `avatar` longtext COLLATE utf8mb4_general_ci COMMENT 'avatar url',
  `nickname` longtext COLLATE utf8mb4_general_ci COMMENT 'nickname',
  `introduction` longtext COLLATE utf8mb4_general_ci COMMENT 'introduction',
  `status` tinyint(1) DEFAULT '1' COMMENT 'status(0: disabled, 1: enable)',
  `role_id` bigint unsigned DEFAULT NULL COMMENT 'role id',
  `last_login` datetime(3) DEFAULT NULL COMMENT 'last login time',
  `locked` tinyint(1) DEFAULT '0' COMMENT 'locked(0: unlock, 1: locked)',
  `lock_expire` bigint DEFAULT NULL COMMENT 'lock expiration time',
  `wrong` bigint DEFAULT NULL COMMENT 'type wrong password count',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_username` (`username`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `tb_sys_role` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'auto increment id',
  `created_at` datetime(3) DEFAULT NULL COMMENT 'create time',
  `updated_at` datetime(3) DEFAULT NULL COMMENT 'update time',
  `deleted_at` datetime(3) DEFAULT NULL COMMENT 'soft delete time',
  `name` longtext COLLATE utf8mb4_general_ci COMMENT 'name',
  `keyword` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT 'keyword(unique str)',
  `desc` longtext COLLATE utf8mb4_general_ci COMMENT 'description',
  `status` tinyint(1) DEFAULT '1' COMMENT 'status(0: disabled, 1: enable)',
  `sort` bigint unsigned DEFAULT '1' COMMENT 'sort(>=0, the smaller the sort, the greater the permission, sort=0 is a super admin)',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_keyword` (`keyword`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `tb_sys_menu` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'auto increment id',
  `created_at` datetime(3) DEFAULT NULL COMMENT 'create time',
  `updated_at` datetime(3) DEFAULT NULL COMMENT 'update time',
  `deleted_at` datetime(3) DEFAULT NULL COMMENT 'soft delete time',
  `name` longtext COLLATE utf8mb4_general_ci COMMENT 'name',
  `title` longtext COLLATE utf8mb4_general_ci COMMENT 'title',
  `icon` longtext COLLATE utf8mb4_general_ci COMMENT 'icon',
  `path` longtext COLLATE utf8mb4_general_ci COMMENT 'url path',
  `redirect` longtext COLLATE utf8mb4_general_ci COMMENT 'redirect url',
  `component` longtext COLLATE utf8mb4_general_ci COMMENT 'ui component name',
  `permission` longtext COLLATE utf8mb4_general_ci COMMENT 'permission',
  `sort` int unsigned DEFAULT NULL COMMENT 'sort(>=0)',
  `status` tinyint(1) DEFAULT '1' COMMENT 'status(0: disabled 1: enabled)',
  `visible` tinyint(1) DEFAULT '1' COMMENT 'visible(0: hidden 1: visible)',
  `breadcrumb` tinyint(1) DEFAULT '1' COMMENT 'breadcrumb(0: disabled 1: enabled)',
  `parent_id` bigint unsigned DEFAULT '0' COMMENT 'parent menu id',
  PRIMARY KEY (`id`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=21 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `tb_sys_menu_role_relation` (
  `menu_id` bigint unsigned DEFAULT NULL,
  `role_id` bigint unsigned DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `tb_leave` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'auto increment id',
  `created_at` datetime(3) DEFAULT NULL COMMENT 'create time',
  `updated_at` datetime(3) DEFAULT NULL COMMENT 'update time',
  `deleted_at` datetime(3) DEFAULT NULL COMMENT 'soft delete time',
  `user_id` bigint unsigned DEFAULT NULL COMMENT 'user id(SysUser.Id)',
  `fsm_uuid` varchar(100) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT 'finite state machine uuid',
  `status` bigint unsigned DEFAULT '0' COMMENT 'status(0:submitted 1:approved 2:refused 3:cancel 4:approving 5:waiting confirm)',
  `approval_opinion` longtext COLLATE utf8mb4_general_ci COMMENT 'approval opinion or remark',
  `desc` longtext COLLATE utf8mb4_general_ci COMMENT 'submitter description',
  `start_time` datetime(3) DEFAULT NULL COMMENT 'start time',
  `end_time` datetime(3) DEFAULT NULL COMMENT 'end time',
  PRIMARY KEY (`id`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `tb_sys_machine` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'auto increment id',
  `created_at` datetime(3) DEFAULT NULL COMMENT 'create time',
  `updated_at` datetime(3) DEFAULT NULL COMMENT 'update time',
  `deleted_at` datetime(3) DEFAULT NULL COMMENT 'soft delete time',
  `host` longtext COLLATE utf8mb4_general_ci COMMENT 'host(IP/Domain)',
  `ssh_port` bigint DEFAULT NULL COMMENT 'ssh port',
  `version` longtext COLLATE utf8mb4_general_ci COMMENT 'os version',
  `name` longtext COLLATE utf8mb4_general_ci COMMENT 'os name',
  `arch` longtext COLLATE utf8mb4_general_ci COMMENT 'os architecture',
  `cpu` longtext COLLATE utf8mb4_general_ci COMMENT 'CPU model',
  `memory` longtext COLLATE utf8mb4_general_ci COMMENT 'memory size',
  `disk` longtext COLLATE utf8mb4_general_ci COMMENT 'disk size',
  `login_name` longtext COLLATE utf8mb4_general_ci COMMENT 'login name',
  `login_pwd` longtext COLLATE utf8mb4_general_ci COMMENT 'login password',
  `status` bigint unsigned DEFAULT '0' COMMENT 'status(0:unhealthy 1:healthy)',
  `remark` longtext COLLATE utf8mb4_general_ci COMMENT 'remark',
  PRIMARY KEY (`id`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `tb_sys_message` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'auto increment id',
  `created_at` datetime(3) DEFAULT NULL COMMENT 'create time',
  `updated_at` datetime(3) DEFAULT NULL COMMENT 'update time',
  `deleted_at` datetime(3) DEFAULT NULL COMMENT 'soft delete time',
  `from_user_id` bigint unsigned DEFAULT NULL COMMENT 'sender user id',
  `title` longtext COLLATE utf8mb4_general_ci COMMENT 'title',
  `content` longtext COLLATE utf8mb4_general_ci COMMENT 'content',
  `type` tinyint DEFAULT '0' COMMENT 'type(0: one2one, 1: one2more, 2: system(one2all))',
  `role_id` bigint unsigned DEFAULT NULL COMMENT 'role id',
  `expired_at` datetime(3) DEFAULT NULL COMMENT 'expire time',
  PRIMARY KEY (`id`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `tb_sys_message_log` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'auto increment id',
  `created_at` datetime(3) DEFAULT NULL COMMENT 'create time',
  `updated_at` datetime(3) DEFAULT NULL COMMENT 'update time',
  `deleted_at` datetime(3) DEFAULT NULL COMMENT 'soft delete time',
  `to_user_id` bigint unsigned DEFAULT NULL COMMENT 'receiver user id',
  `message_id` bigint unsigned DEFAULT NULL COMMENT 'message id',
  `status` tinyint DEFAULT '0' COMMENT 'status(0: unread, 1: read, 2: deleted)',
  PRIMARY KEY (`id`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `tb_sys_operation_log` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'auto increment id',
  `created_at` datetime(3) DEFAULT NULL COMMENT 'create time',
  `updated_at` datetime(3) DEFAULT NULL COMMENT 'update time',
  `deleted_at` datetime(3) DEFAULT NULL COMMENT 'soft delete time',
  `api_desc` longtext COLLATE utf8mb4_general_ci COMMENT 'api description',
  `path` longtext COLLATE utf8mb4_general_ci COMMENT 'url path',
  `method` longtext COLLATE utf8mb4_general_ci COMMENT 'api method',
  `header` blob COMMENT 'request header',
  `body` blob COMMENT 'request body',
  `params` blob COMMENT 'request params',
  `resp` blob COMMENT 'response data',
  `status` bigint DEFAULT NULL COMMENT 'response status',
  `username` longtext COLLATE utf8mb4_general_ci COMMENT 'login username',
  `role_name` longtext COLLATE utf8mb4_general_ci COMMENT 'login role name',
  `ip` longtext COLLATE utf8mb4_general_ci COMMENT 'IP',
  `ip_location` longtext COLLATE utf8mb4_general_ci COMMENT 'real location of the IP',
  `latency` bigint DEFAULT NULL COMMENT 'request time(ms)',
  `user_agent` longtext COLLATE utf8mb4_general_ci COMMENT 'browser user agent',
  PRIMARY KEY (`id`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
