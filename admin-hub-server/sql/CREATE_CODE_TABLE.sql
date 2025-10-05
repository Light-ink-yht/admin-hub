-- =================================================================================
-- 验证码表 CD01
-- 用途：存储系统中生成的各类验证码信息，支持多种业务场景下的验证码验证，如注册、登录、密码重置等
-- 设计原则：
-- 1. 支持多种业务类型（如邮箱验证码、短信验证码等）
-- 2. 支持验证码有效期管理
-- 3. 支持验证次数限制
-- 4. 支持状态追踪（未使用、已使用、已过期、已失效）
-- 5. 通过索引优化查询性能
-- =================================================================================
CREATE TABLE IF NOT EXISTS CD01 (
    -- 验证码字段
    CDA001 VARCHAR(36) NOT NULL PRIMARY KEY COMMENT '验证码ID（主键）',
    CDA002 VARCHAR(50) NOT NULL COMMENT '业务类型(biz)，如email、sms等',
    CDA003 VARCHAR(100) NOT NULL COMMENT '用户输入(input)，如邮箱地址、手机号等',
    CDA004 VARCHAR(20) NOT NULL COMMENT '验证码内容',
    CDA005 DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '验证码创建时间',
    CDA006 DATETIME NOT NULL COMMENT '验证码过期时间',
    CDA007 INT NOT NULL DEFAULT 3 COMMENT '最大验证次数',
    CDA008 INT NOT NULL DEFAULT 0 COMMENT '已验证次数',
    CDA009 INT NOT NULL DEFAULT 1 COMMENT '验证码状态：1-有效（未使用），2-已使用，3-已过期，4-已失效（验证次数过多）',
    CDA010 VARCHAR(255) COMMENT '备注信息',
    
    -- 索引定义
    UNIQUE INDEX IDX_UQ_CDA002_CDA003_CDA004 (CDA002, CDA003, CDA004) COMMENT '同一业务类型、用户输入和验证码内容下的唯一索引',
    INDEX IDX_CDA002 (CDA002) COMMENT '按业务类型查询索引',
    INDEX IDX_CDA003 (CDA003) COMMENT '按用户输入查询索引',
    INDEX IDX_CDA002_CDA003 (CDA002, CDA003) COMMENT '按业务类型和用户输入组合查询索引',
    INDEX IDX_CDA006 (CDA006) COMMENT '按过期时间查询索引，用于定期清理过期数据',
    INDEX IDX_CDA005 (CDA005) COMMENT '按创建时间查询索引',
    INDEX IDX_CDA009 (CDA009) COMMENT '按验证码状态查询索引'
) ENGINE=INNODB DEFAULT CHARSET=utf8mb4 COMMENT='验证码表';

-- =================================================================================
-- 表结构说明：
-- 1. 表名CD01遵循项目统一的命名规范，CD表示验证码相关，01表示具体表序号
-- 2. 验证码字段（CDA001-CDA010）使用统一的命名规范，前缀CDA表示验证码相关，后三位数字表示具体字段序号
-- 3. CDA001作为主键，唯一标识每条验证码记录
-- 4. 索引设计充分考虑了查询需求，包括唯一索引和组合索引，优化了查询性能
-- 5. 状态字段CDA009定义了四种状态，便于追踪验证码的完整生命周期
-- 6. 存储引擎使用InnoDB，支持事务和行级锁，确保数据一致性
-- 7. 字符集使用utf8mb4，支持更广泛的字符范围
-- =================================================================================