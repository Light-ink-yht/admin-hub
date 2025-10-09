-- 短信配置表 SM01
-- 存储短信服务商的配置信息，如腾讯云、阿里云等不同服务商的密钥信息
CREATE TABLE IF NOT EXISTS SM01 (
    -- 短信配置字段
    SMA001 VARCHAR(36) NOT NULL PRIMARY KEY COMMENT '配置ID（主键）',
    SMA002 VARCHAR(100) NOT NULL COMMENT '配置名称（如：腾讯云主账号、阿里云备用账号）',
    SMA003 VARCHAR(50) NOT NULL COMMENT '短信服务商类型（tencent、aliyun等）',
    SMA004 VARCHAR(100) NOT NULL COMMENT '应用ID（AppId/SdkAppId）',
    SMA005 VARCHAR(255) NOT NULL COMMENT '应用密钥（SecretId/AccessKeyId）',
    SMA006 VARCHAR(255) NOT NULL COMMENT '应用密钥（SecretKey/AccessKeySecret）',
    SMA007 VARCHAR(100) COMMENT '短信签名',
    SMA008 INT NOT NULL DEFAULT 1 COMMENT '配置状态：1-启用，2-禁用',
    SMA009 DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    SMA010 DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    SMA011 VARCHAR(255) COMMENT '备注信息',

    -- 索引定义
    INDEX IDX_SMA002 (SMA002) COMMENT '按配置名称查询索引',
    INDEX IDX_SMA003 (SMA003) COMMENT '按服务商类型查询索引',
    INDEX IDX_SMA008 (SMA008) COMMENT '按配置状态查询索引',
    INDEX IDX_SMA009 (SMA009) COMMENT '按创建时间查询索引',
    INDEX IDX_SMA010 (SMA010) COMMENT '按修改时间查询索引'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='短信服务商配置表';

-- 短信模板表 SM02
-- 存储短信模板信息，包括模板ID和模板内容描述
CREATE TABLE IF NOT EXISTS SM02 (
    -- 短信模板字段
    SMB001 VARCHAR(36) NOT NULL PRIMARY KEY COMMENT '模板ID（主键）',
    SMB002 VARCHAR(100) NOT NULL COMMENT '模板名称',
    SMB003 VARCHAR(50) NOT NULL COMMENT '模板标识（唯一标识）',
    SMB004 VARCHAR(50) NOT NULL COMMENT '短信服务商模板ID',
    SMB005 VARCHAR(255) COMMENT '模板内容描述',
    SMB006 INT NOT NULL DEFAULT 1 COMMENT '模板状态：1-启用，2-禁用',
    SMB007 DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    SMB008 DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    SMB009 VARCHAR(255) COMMENT '备注信息',

    -- 索引定义
    UNIQUE INDEX IDX_UQ_SMB003 (SMB003) COMMENT '按模板标识唯一索引',
    INDEX IDX_SMB002 (SMB002) COMMENT '按模板名称查询索引',
    INDEX IDX_SMB006 (SMB006) COMMENT '按模板状态查询索引',
    INDEX IDX_SMB007 (SMB007) COMMENT '按创建时间查询索引',
    INDEX IDX_SMB008 (SMB008) COMMENT '按修改时间查询索引'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='短信模板表';

-- ========================================================
-- 表结构说明：
-- 1. 表名SM01和SM02遵循项目统一的命名规范，SM表示短信相关，01和02表示具体表序号
-- 2. 字段命名使用统一的规范，前缀SMA和SMB分别表示短信配置和短信模板相关
-- 3. 作为主键的ID字段（SMA001、SMB001）使用VARCHAR(36)类型，兼容UUID
-- 4. 索引设计充分考虑了查询需求，包括唯一索引和普通索引
-- 5. 状态字段定义了启用/禁用两种状态，便于配置和模板的管理
-- 6. 存储引擎使用InnoDB，支持事务和行级锁，确保数据一致性
-- 7. 字符集使用utf8mb4，支持更广泛的字符范围
-- ========================================================