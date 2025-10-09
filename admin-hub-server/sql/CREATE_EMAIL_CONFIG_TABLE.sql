-- =================================================================================
-- 邮件服务器配置表 EM01
-- 用途：存储SMTP邮件服务器配置信息，支持系统中各类邮件发送功能
-- 设计原则：
-- 1. 支持多套邮件服务器配置
-- 2. 支持启用/禁用配置项
-- 3. 通过索引优化查询性能
-- 4. 记录配置的创建和修改信息
-- =================================================================================
CREATE TABLE IF NOT EXISTS EM01 (
    -- 邮件配置字段
    EMA001 VARCHAR(36) NOT NULL PRIMARY KEY COMMENT '配置ID（主键）',
    EMA002 VARCHAR(100) NOT NULL COMMENT '发送方邮箱地址',
    EMA003 VARCHAR(255) NOT NULL COMMENT '发送方邮箱授权码',
    EMA004 VARCHAR(100) NOT NULL COMMENT 'SMTP服务器地址',
    EMA005 INT NOT NULL COMMENT 'SMTP服务器端口',
    EMA006 VARCHAR(100) COMMENT '配置名称（如：主邮箱、备用邮箱）',
    EMA007 INT NOT NULL DEFAULT 1 COMMENT '配置状态：1-启用，2-禁用',
    EMA008 DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    EMA009 DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    EMA010 VARCHAR(255) COMMENT '备注信息',
    
    -- 索引定义
    UNIQUE INDEX IDX_UQ_EMA002_EMA004_EMA005 (EMA002, EMA004, EMA005) COMMENT '同一邮箱、服务器和端口的唯一索引',
    INDEX IDX_EMA007 (EMA007) COMMENT '按配置状态查询索引',
    INDEX IDX_EMA008 (EMA008) COMMENT '按创建时间查询索引',
    INDEX IDX_EMA009 (EMA009) COMMENT '按修改时间查询索引'
) ENGINE=INNODB DEFAULT CHARSET=utf8mb4 COMMENT='邮件服务器配置表';

-- =================================================================================
-- 邮件模板表 EM02
-- 用途：存储系统中使用的各类邮件模板，如验证码、通知、营销等邮件模板
-- 设计原则：
-- 1. 支持多类型邮件模板
-- 2. 支持模板启用/禁用管理
-- 3. 支持模板内容的富文本格式
-- 4. 通过索引优化查询性能
-- =================================================================================
CREATE TABLE IF NOT EXISTS EM02 (
    -- 邮件模板字段
    EMB001 VARCHAR(36) NOT NULL PRIMARY KEY COMMENT '模板ID（主键）',
    EMB002 VARCHAR(100) NOT NULL COMMENT '模板名称（如：验证码模板、欢迎邮件）',
    EMB003 VARCHAR(100) NOT NULL COMMENT '模板标识（唯一标识，如：signup_code、reset_password）',
    EMB004 VARCHAR(255) NOT NULL COMMENT '邮件主题',
    EMB005 TEXT NOT NULL COMMENT '邮件内容（支持HTML格式，可包含占位符如{code}）',
    EMB006 INT NOT NULL DEFAULT 1 COMMENT '模板状态：1-启用，2-禁用',
    EMB007 DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    EMB008 DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    EMB009 VARCHAR(255) COMMENT '备注信息',
    
    -- 索引定义
    UNIQUE INDEX IDX_UQ_EMB003 (EMB003) COMMENT '按模板标识唯一索引',
    INDEX IDX_EMB002 (EMB002) COMMENT '按模板名称查询索引',
    INDEX IDX_EMB006 (EMB006) COMMENT '按模板状态查询索引',
    INDEX IDX_EMB007 (EMB007) COMMENT '按创建时间查询索引',
    INDEX IDX_EMB008 (EMB008) COMMENT '按修改时间查询索引'
) ENGINE=INNODB DEFAULT CHARSET=utf8mb4 COMMENT='邮件模板表';

-- =================================================================================
-- 表结构说明：
-- 1. 表名EM01和EM02遵循项目统一的命名规范，EM表示邮件相关，01和02表示具体表序号
-- 2. 字段命名使用统一的规范，前缀EMA和EMB分别表示邮件配置和邮件模板相关
-- 3. 作为主键的ID字段（EMA001、EMB001）使用VARCHAR(36)类型，兼容UUID
-- 4. 索引设计充分考虑了查询需求，包括唯一索引和普通索引
-- 5. 状态字段定义了启用/禁用两种状态，便于配置和模板的管理
-- 6. 存储引擎使用InnoDB，支持事务和行级锁，确保数据一致性
-- 7. 字符集使用utf8mb4，支持更广泛的字符范围
-- =================================================================================