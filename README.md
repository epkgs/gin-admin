# Gin Admin Template

> 该模板提供了构建企业级Web应用所需的核心功能，包括用户认证、权限管理、数据库操作、缓存等。

## 功能特性
- 提供基于 Gin 框架的轻量级项目模板
- 使用简洁的手动依赖注入，实现清晰的代码结构
- 集成 GORM 进行 ORM 映射和数据库操作
  - 支持 PostgreSQL
  - 支持 MySQL
  - 支持 SQLite
- 集成 Viper 进行配置管理
- 提供常用 Gin 中间件和工具
  - 多语言中间件：支持多语言，使用 [epkgs/i18n](https://github.com/epkgs/i18n) 模块实现
  - 跨域中间件：处理 API 跨域请求，实现 CORS 支持
  - JWT 解析中间件：从请求中解析并验证 JWT Token，用于 API 身份认证
  - Trace 中间件：记录请求的 Trace ID，用于请求链路追踪
  - 日志中间件：记录请求的日志，使用 [zap](https://go.uber.org/zap) 模块实现
  - Copy Body 中间件：复制请求的 Body 内容，用于日志记录
  - Auth 中间件：处理用户认证，基于 [jwt](https://github.com/golang-jwt/jwt) 封装实现
  - Rate Limiter 中间件：实现请求限流
  - Casbin 中间件：实现权限控制，基于 [casbin](https://github.com/casbin/casbin) 封装实现
  - Prometheus 中间件：实现 Prometheus 监控，记录请求次数、错误次数、响应时间等
- 国际化 (i18n) 支持
  - 请求参数 lang 指定语言
  - 自动识别 cookie 中的语言
  - 基于请求 Accept-Language 头自动选择语言
- 使用 Cobra 命令行框架，提供清晰的子命令结构
- Swagger 文档生成

## 架构设计
项目采用简化的架构设计，使用手动依赖注入，实现了清晰的代码结构：

## 目录结构
```
.
├── cmd                 # 命令行工具
├── configs             # 配置文件
├── internal            # 核心业务逻辑
│   ├── apis            # API控制器
│   ├── app             # 应用初始化
│   ├── configs         # 配置解析
│   ├── defines         # 常量定义
│   ├── dtos            # 数据传输对象
│   ├── errorx          # 错误处理
│   ├── models          # 数据模型
│   ├── repositories    # 数据访问层
│   ├── services        # 业务逻辑层
│   ├── swagger         # API文档
│   └── types           # 接口定义
├── locales             # 多语言文件
├── pkg                 # 公共组件库
├── scripts             # 脚本文件
├── test                # 测试文件
├── data                # 数据文件（运行时生成）
└── uploads             # 上传文件（运行时生成）
```

## 快速开始
### 环境要求
- Go 1.18+

### 安装
```bash
# 克隆项目
git clone git@github.com:epkgs/gin-admin.git

# 进入项目目录
cd gin-admin

# 安装依赖包
go mod tidy
```

### 运行参数
- start: 启动服务器
  - -c, --config: 指定配置文件路径
  - -d, --deamon: 运行为守护进程
- stop: 停止服务器
- version: 显示版本信息

## 许可证
[MIT License](https://github.com/epkgs/gin-admin/blob/master/LICENSE)