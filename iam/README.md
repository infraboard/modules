# 身份管理模块

一套简易版的RBAC用户认证与鉴权模块:

![](./docs/arch.png)

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

```sh
$ modules/iam/example ‹main*› » go run main.go init                                                                               1 ↵
? 请输入管理员用户名称: admin
? 请输入管理员密码: ******
? 再次输入管理员密码: ******

2024/03/01 14:46:19 /Users/yumaojun/Workspace/Golang/inforboard/modules/iam/apps/user/impl/mysql/user.go:26
[5.253ms] [rows:1] INSERT INTO `users` (`created_at`,`updated_at`,`username`,`password`,`role`,`label`) VALUES (1709275579,1709275579,'admin','$2a$10$87u3qGH1K6/XOERRdpD2RODJlQqLF8iODACgY.oacgQZ1Jf0JSZlm','admin','{}') RETURNING `id`
{"id":9,"created_at":1709275579,"updated_at":1709275579,"username":"admin","password":"$2a$10$87u3qGH1K6/XOERRdpD2RODJlQqLF8iODACgY.oacgQZ1Jf0JSZlm","role":"admin","label":{}}
```

### 业务接口开发

```go
func (h *ApiHandler) DBStats(ctx *gin.Context) {
	db, _ := h.db.DB()
	ctx.JSON(http.StatusOK, gin.H{
		"data": db.Stats(),
	})
}
```

### 开启认证与鉴权


```go
import (
	// 引入IAM模块组件
	_ "github.com/infraboard/modules/iam"
	// 引入IAM模块CLI工具
	_ "github.com/infraboard/modules/iam/cmd"
)
```

```go
// API路由
func (h *ApiHandler) Registry(r gin.IRouter) {
	r.Use(middleware.Auth())
	r.GET("/db_stats", middleware.Perm(user.ROLE_ADMIN), h.DBStats)
}
```

### 启动服务并验证

1. 启动服务
```sh
$ modules/iam/example ‹main*› » go run main.go start
2024-03-01T15:05:57+08:00 INFO   ioc/server/server.go:74 > loaded configs: [app.v1 trace.v1 log.v1 datasource.v1 grpc.v1 http.v1] component:SERVER
2024-03-01T15:05:57+08:00 INFO   ioc/server/server.go:75 > loaded controllers: [tokens.v1 users.v1] component:SERVER
2024-03-01T15:05:57+08:00 INFO   ioc/server/server.go:76 > loaded apis: [tokens.v1 users.v1 module_a.v1] component:SERVER
[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] POST   /exapmle/api/v1/tokens/   --> github.com/infraboard/modules/iam/apps/token/api/gin.(*TokenApiHandler).Login-fm (3 handlers)
[GIN-debug] DELETE /exapmle/api/v1/tokens/   --> github.com/infraboard/modules/iam/apps/token/api/gin.(*TokenApiHandler).Logout-fm (3 handlers)
[GIN-debug] GET    /exapmle/api/v1/users/    --> github.com/infraboard/modules/iam/apps/user/api/gin.(*UserApiHandler).QueryUser-fm (5 handlers)
[GIN-debug] GET    /exapmle/api/v1/users/:id --> github.com/infraboard/modules/iam/apps/user/api/gin.(*UserApiHandler).DescribeUser-fm (5 handlers)
[GIN-debug] POST   /exapmle/api/v1/users/    --> github.com/infraboard/modules/iam/apps/user/api/gin.(*UserApiHandler).CreateUser-fm (5 handlers)
[GIN-debug] DELETE /exapmle/api/v1/users/:id --> github.com/infraboard/modules/iam/apps/user/api/gin.(*UserApiHandler).DeleteUser-fm (5 handlers)
[GIN-debug] GET    /exapmle/api/v1/module_a/db_stats --> main.(*ApiHandler).DBStats-fm (5 handlers)
2024-03-01T15:05:57+08:00 INFO   config/http/http.go:211 > HTTP服务启动成功, 监听地址: 127.0.0.1:8020 component:HTTP
```

2.  管理员登录:
```sh
curl --location 'http://localhost:8020/exapmle/api/v1/tokens/' \
--header 'Content-Type: application/json' \
--data '{
    "username": "admin",
    "password": "123456"
}
'
```

```json
{
    "code": 0,
    "data": {
        "user_id": "9",
        "username": "admin",
        "access_token": "cngntlhus0n4irgpns0g",
        "access_token_expired_at": 7200,
        "refresh_token": "cngntlhus0n4irgpns10",
        "refresh_token_expired_at": 604800,
        "created_at": 1709276886,
        "updated_at": 1709276886,
        "role": "admin"
    }
}
```

3. 管理员测试接口权限:
```sh
curl --location --request GET 'http://localhost:8020/exapmle/api/v1/module_a/db_stats' \
--header 'Content-Type: application/json' \
--header 'Cookie: access_token=cngntlhus0n4irgpns0g' \
```

```json
{
    "data": {
        "MaxOpenConnections": 0,
        "OpenConnections": 1,
        "InUse": 0,
        "Idle": 1,
        "WaitCount": 0,
        "WaitDuration": 0,
        "MaxIdleClosed": 0,
        "MaxIdleTimeClosed": 0,
        "MaxLifetimeClosed": 0
    }
}
```

4. 创建普通账号:
```sh
curl --location 'http://localhost:8020/exapmle/api/v1/users' \
--header 'Content-Type: application/json' \
--header 'Cookie: access_token=cngnv79us0n4irgpns1g' \
--data '{
    "username": "guest",
    "password": "123456"
}
```

5. 使用普通账号登录:

```sh
curl --location 'http://localhost:8020/exapmle/api/v1/tokens/' \
--header 'Content-Type: application/json' \
--header 'Cookie: access_token=cngo23hus0n5dkkfl1p0' \
--data '{
    "username": "guest",
    "password": "123456"
}
'
```

```json
{
    "code": 0,
    "data": {
        "user_id": "10",
        "username": "guest",
        "access_token": "cngo23hus0n5dkkfl1p0",
        "access_token_expired_at": 7200,
        "refresh_token": "cngo23hus0n5dkkfl1pg",
        "refresh_token_expired_at": 604800,
        "created_at": 1709277454,
        "updated_at": 1709277454,
        "role": "member"
    }
}
```

6. 普通账号权限测试

```sh
curl --location --request GET 'http://localhost:8020/exapmle/api/v1/module_a/db_stats' \
--header 'Content-Type: application/json' \
--header 'Cookie: access_token=cngo2jhus0n5h459bvq0' \
--data '{
    "username": "guest",
    "password": "123456"
}
'
```

```json
{
    "namespace": "exapmle",
    "http_code": 403,
    "error_code": 403,
    "reason": "访问未授权",
    "message": "role member not allow ",
    "meta": null,
    "data": null
}
```

