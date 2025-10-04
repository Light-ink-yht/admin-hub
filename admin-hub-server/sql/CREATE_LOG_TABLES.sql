-- =================================================================================
-- 常用日志表 LL01
-- 用途：记录系统中各类通用操作日志，如系统启动、配置加载、服务状态变化等
-- =================================================================================
CREATE TABLE IF NOT EXISTS LL01 (
    -- 基础字段（所有表共用）
    LAN001 BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    LAN002 DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    LAN003 DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    LAN004 DATETIME COMMENT '删除时间',
    LAN005 VARCHAR(1) NOT NULL DEFAULT '1' COMMENT '是否删除 1 否，2 是',
    LAN006 BIGINT NOT NULL DEFAULT 0 COMMENT '创建人',
    LAN007 BIGINT NOT NULL DEFAULT 0 COMMENT '修改人',
    
    -- 日志特有字段
    LLA001 VARCHAR(50) NOT NULL COMMENT '日志ID',
    LLA002 VARCHAR(10) NOT NULL COMMENT '日志级别（DEBUG/INFO/WARN/ERROR/PANIC/FATAL）',
    LLA003 VARCHAR(255) NOT NULL COMMENT '日志消息',
    LLA004 TEXT COMMENT '日志字段（JSON格式，存储结构化数据）',
    LLA005 VARCHAR(255) COMMENT '调用位置（代码文件:行号）',
    LLA006 VARCHAR(20) COMMENT '环境标识（dev/prod）',
    LLA007 DATETIME NOT NULL COMMENT '日志时间',
    
    -- 索引定义
    INDEX IDX_LLA007 (LLA007) COMMENT '按日志时间查询索引',
    INDEX IDX_LLA001 (LLA001) COMMENT '按日志ID查询索引',
    INDEX IDX_LLA002 (LLA002) COMMENT '按日志级别查询索引'
) ENGINE=INNODB DEFAULT CHARSET=utf8mb4 COMMENT='常用日志表';

-- =================================================================================
-- 登录日志表 LL02
-- 用途：记录用户登录相关的日志信息，用于安全审计和用户行为追踪
-- =================================================================================
CREATE TABLE IF NOT EXISTS LL02 (
    -- 基础字段（所有表共用）
    LAN001 BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    LAN002 DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    LAN003 DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    LAN004 DATETIME COMMENT '删除时间',
    LAN005 VARCHAR(1) NOT NULL DEFAULT '1' COMMENT '是否删除 1 否，2 是',
    LAN006 BIGINT NOT NULL DEFAULT 0 COMMENT '创建人',
    LAN007 BIGINT NOT NULL DEFAULT 0 COMMENT '修改人',
    
    -- 登录日志特有字段
    LLA001 VARCHAR(50) NOT NULL COMMENT '日志ID',
    LLA002 VARCHAR(10) NOT NULL COMMENT '日志级别（INFO/ERROR）',
    LLA003 VARCHAR(255) NOT NULL COMMENT '日志消息',
    LLA004 VARCHAR(50) COMMENT '用户ID',
    LLA005 VARCHAR(50) COMMENT '用户名',
    LLA006 VARCHAR(50) COMMENT '登录IP',
    LLA007 VARCHAR(255) COMMENT '登录设备信息（浏览器、操作系统等）',
    LLA008 VARCHAR(10) COMMENT '登录状态（成功/失败）',
    LLA009 VARCHAR(255) COMMENT '失败原因（如密码错误、账户锁定等）',
    LLA010 TEXT COMMENT '日志字段（JSON格式，存储其他登录相关信息）',
    LLA011 DATETIME NOT NULL COMMENT '登录时间',
    
    -- 索引定义
    INDEX IDX_LLA011 (LLA011) COMMENT '按登录时间查询索引',
    INDEX IDX_LLA001 (LLA001) COMMENT '按日志ID查询索引',
    INDEX IDX_LLA004 (LLA004) COMMENT '按用户ID查询索引',
    INDEX IDX_LLA006 (LLA006) COMMENT '按登录IP查询索引'
) ENGINE=INNODB DEFAULT CHARSET=utf8mb4 COMMENT='登录日志表';

-- =================================================================================
-- 业务日志表 LL03
-- 用途：记录用户在各业务模块的操作日志，用于业务审计和操作追溯
-- =================================================================================
CREATE TABLE IF NOT EXISTS LL03 (
    -- 基础字段（所有表共用）
    LAN001 BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    LAN002 DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    LAN003 DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    LAN004 DATETIME COMMENT '删除时间',
    LAN005 VARCHAR(1) NOT NULL DEFAULT '1' COMMENT '是否删除 1 否，2 是',
    LAN006 BIGINT NOT NULL DEFAULT 0 COMMENT '创建人',
    LAN007 BIGINT NOT NULL DEFAULT 0 COMMENT '修改人',
    
    -- 业务日志特有字段
    LLA001 VARCHAR(50) NOT NULL COMMENT '日志ID',
    LLA002 VARCHAR(10) NOT NULL COMMENT '日志级别（INFO/ERROR）',
    LLA003 VARCHAR(50) COMMENT '业务模块（如用户管理、系统配置等）',
    LLA004 VARCHAR(50) COMMENT '操作类型（如新增、修改、删除、查询等）',
    LLA005 VARCHAR(50) COMMENT '业务ID（关联具体业务对象）',
    LLA006 VARCHAR(50) COMMENT '操作人ID',
    LLA007 VARCHAR(50) COMMENT '操作人姓名',
    LLA008 VARCHAR(255) COMMENT '操作内容（描述具体执行了什么操作）',
    LLA009 VARCHAR(10) COMMENT '操作状态（成功/失败）',
    LLA010 VARCHAR(255) COMMENT '失败原因（如参数错误、权限不足等）',
    LLA011 VARCHAR(50) COMMENT '操作IP',
    LLA012 TEXT COMMENT '日志字段（JSON格式，存储操作前后的数据快照等）',
    LLA013 DATETIME NOT NULL COMMENT '操作时间',
    
    -- 索引定义
    INDEX IDX_LLA013 (LLA013) COMMENT '按操作时间查询索引',
    INDEX IDX_LLA001 (LLA001) COMMENT '按日志ID查询索引',
    INDEX IDX_LLA003 (LLA003) COMMENT '按业务模块查询索引',
    INDEX IDX_LLA006 (LLA006) COMMENT '按操作人ID查询索引'
) ENGINE=INNODB DEFAULT CHARSET=utf8mb4 COMMENT='业务日志表';