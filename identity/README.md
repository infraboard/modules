# 身份管理模块

一套简易版的RBAC用户认证与鉴权模块:

包含如下接口:
+ 令牌管理:
	+ 令牌颁发
	+ 令牌撤销
+ 用户管理
	+ 用户创建
	+ 用户查询
	+ 用户删除

## 快速使用

### 初始化

1. 初始化SQL
```sql
CREATE TABLE `users` (
	`id` int unsigned NOT NULL AUTO_INCREMENT,
	`created_at` int NOT NULL COMMENT '创建时间',
	`updated_at` int NOT NULL COMMENT '更新时间',
	`username` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '用户名, 用户名不允许重复的',
	`password` varchar(255) COLLATE utf8mb4_general_ci NOT NULL COMMENT '不能保持用户的明文密码',
	`label` varchar(255) COLLATE utf8mb4_general_ci NOT NULL COMMENT '用户标签',
	`role` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '用户的角色',
	PRIMARY KEY (`id`) USING BTREE,
	UNIQUE KEY `idx_user` (`username`)
  ) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE `tokens` (
	`created_at` int NOT NULL COMMENT '创建时间',
	`updated_at` int NOT NULL COMMENT '更新时间',
	`user_id` int NOT NULL COMMENT '用户的Id',
	`username` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '用户名, 用户名不允许重复的',
	`access_token` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '用户的访问令牌',
	`access_token_expired_at` int NOT NULL COMMENT '令牌过期时间',
	`refresh_token` varchar(255) COLLATE utf8mb4_general_ci NOT NULL COMMENT '刷新令牌',
	`refresh_token_expired_at` int NOT NULL COMMENT '刷新令牌过期时间',
	PRIMARY KEY (`access_token`) USING BTREE,
	UNIQUE KEY `idx_token` (`access_token`) USING BTREE
  ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
```

2. 初始化管理员


### 业务接口开发




### 开启认证与鉴权





