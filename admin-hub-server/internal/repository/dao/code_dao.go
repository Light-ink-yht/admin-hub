// ====================================================================
// dao 包
// 提供与数据库交互的数据访问对象
// ====================================================================
package dao

import (
	"context"
	"database/sql"
	"github.com/Light-ink-yht/admin-hub/internal/domain/code_domain"
	"time"
)

// ====================================================================
// CodeDAO 验证码数据访问对象
// 负责与数据库交互，执行验证码的增删改查操作
// 包含创建表、插入、查询、更新状态等方法
// 主要功能：
// 1. 创建验证码表结构
// 2. 插入新的验证码记录
// 3. 根据业务类型和用户输入查询验证码
// 4. 更新验证码状态和验证次数
// 5. 统计指定时间范围内的验证码发送数量
// 6. 清理过期的验证码记录
// ====================================================================

type CodeDAO struct {
	db *sql.DB // 数据库连接对象
}

// ====================================================================
// NewCodeDAO 创建新的验证码DAO实例
// 参数：
// - db: 数据库连接对象
// 返回值：
// - *CodeDAO: 验证码DAO实例
// ====================================================================
func NewCodeDAO(db *sql.DB) *CodeDAO {
	return &CodeDAO{
		db: db,
	}
}

// ====================================================================
// CreateCodeTable 创建验证码表
// 在数据库中创建CD01验证码表，包含所需的字段和索引
// 参数：
// - ctx: 上下文，用于控制超时和取消
// 返回值：
// - error: 操作错误，成功时为nil
// ====================================================================
func (dao *CodeDAO) CreateCodeTable(ctx context.Context) error {
	query := `
	CREATE TABLE IF NOT EXISTS CD01 (
		-- 主键ID，用于唯一标识验证码记录
		CDA001 VARCHAR(36) PRIMARY KEY COMMENT '主键ID',
		-- 业务类型，标识验证码的用途，如email、sms等
		CDA002 VARCHAR(50) NOT NULL COMMENT '业务类型(biz)',
		-- 用户输入，存储用户的唯一标识，如邮箱地址、手机号等
		CDA003 VARCHAR(100) NOT NULL COMMENT '用户输入(input，邮箱或手机号)',
		-- 验证码内容，存储生成的验证码字符串
		CDA004 VARCHAR(20) NOT NULL COMMENT '验证码(code)',
		-- 创建时间，记录验证码生成的时间
		CDA005 DATETIME NOT NULL COMMENT '创建时间',
		-- 过期时间，记录验证码的过期时间
		CDA006 DATETIME NOT NULL COMMENT '过期时间',
		-- 验证次数限制，限制验证码的最大验证尝试次数
		CDA007 INT NOT NULL DEFAULT 3 COMMENT '验证次数限制',
		-- 已验证次数，记录验证码已被验证的次数
		CDA008 INT NOT NULL DEFAULT 0 COMMENT '已验证次数',
		-- 状态，标识验证码的当前状态
		-- 1: 有效 - 验证码已生成但尚未验证
		-- 2: 已使用 - 验证码已成功验证
		-- 3: 已过期 - 验证码已超过有效期
		-- 4: 已失效 - 验证码验证次数已达上限
		CDA009 INT NOT NULL DEFAULT 1 COMMENT '状态(1:有效, 2:已使用, 3:已过期, 4:已失效)',
		-- 备注信息，用于存储额外的说明或上下文信息
		CDA010 VARCHAR(255) DEFAULT NULL COMMENT '备注信息',
		-- 唯一索引，确保同一业务类型和用户输入下只有一条有效验证码
		UNIQUE KEY uk_biz_user_code (CDA002, CDA003, CDA004),
		-- 索引，优化按业务类型和用户输入查询的性能
		INDEX idx_biz_user (CDA002, CDA003),
		-- 索引，优化按过期时间查询的性能
		INDEX idx_expired_at (CDA006),
		-- 索引，优化按创建时间查询的性能
		INDEX idx_created_at (CDA005)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='验证码表';
	`

	// 执行SQL创建表
	_, err := dao.db.ExecContext(ctx, query)
	return err
}

// ====================================================================
// Insert 插入验证码记录
// 将验证码领域模型数据插入到数据库中
// 使用ON DUPLICATE KEY UPDATE确保在唯一键冲突时更新现有记录
// 参数：
// - ctx: 上下文，用于控制超时和取消
// - code: 验证码领域模型对象
// 返回值：
// - error: 操作错误，成功时为nil
// ====================================================================
func (dao *CodeDAO) Insert(ctx context.Context, code *code_domain.Code) error {
	query := `
	INSERT INTO CD01 (CDA001, CDA002, CDA003, CDA004, CDA005, CDA006, CDA007, CDA008, CDA009, CDA010)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON DUPLICATE KEY UPDATE 
		CDA004 = VALUES(CDA004),
		CDA005 = VALUES(CDA005),
		CDA006 = VALUES(CDA006),
		CDA007 = VALUES(CDA007),
		CDA008 = VALUES(CDA008),
		CDA009 = VALUES(CDA009),
		CDA010 = VALUES(CDA010)
	`

	// 执行SQL插入或更新操作
	_, err := dao.db.ExecContext(
		ctx,
		query,
		code.CDA001, // 主键ID
		code.CDA002, // 业务类型
		code.CDA003, // 用户输入
		code.CDA004, // 验证码内容
		code.CDA005, // 创建时间
		code.CDA006, // 过期时间
		code.CDA007, // 验证次数限制
		code.CDA008, // 已验证次数
		code.CDA009, // 状态
		code.CDA010, // 备注信息
	)
	return err
}

// ====================================================================
// FindByBizAndInput 根据业务类型和用户输入查询验证码
// 查询指定业务类型和用户输入的最新有效验证码
// 参数：
// - ctx: 上下文，用于控制超时和取消
// - biz: 业务类型
// - input: 用户输入
// 返回值：
// - *code_domain.Code: 验证码领域模型对象，未找到时为nil
// - error: 操作错误，成功时为nil
// ====================================================================
func (dao *CodeDAO) FindByBizAndInput(ctx context.Context, biz, input string) (*code_domain.Code, error) {
	query := `
	SELECT CDA001, CDA002, CDA003, CDA004, CDA005, CDA006, CDA007, CDA008, CDA009, CDA010 
	FROM CD01 
	WHERE CDA002 = ? AND CDA003 = ? AND CDA009 IN (1) 
	ORDER BY CDA005 DESC 
	LIMIT 1
	`

	// 初始化验证码对象
	var code code_domain.Code
	// 执行SQL查询并扫描结果
	err := dao.db.QueryRowContext(ctx, query, biz, input).Scan(
		&code.CDA001,
		&code.CDA002,
		&code.CDA003,
		&code.CDA004,
		&code.CDA005,
		&code.CDA006,
		&code.CDA007,
		&code.CDA008,
		&code.CDA009,
		&code.CDA010,
	)

	// 处理未找到记录的情况
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &code, err
}

// ====================================================================
// Update 更新验证码信息
// 更新验证码的验证次数和状态
// 参数：
// - ctx: 上下文，用于控制超时和取消
// - code: 包含更新信息的验证码领域模型对象
// 返回值：
// - error: 操作错误，成功时为nil
// ====================================================================
func (dao *CodeDAO) Update(ctx context.Context, code *code_domain.Code) error {
	query := `
	UPDATE CD01 
	SET CDA008 = ?, CDA009 = ? 
	WHERE CDA001 = ?
	`

	// 执行SQL更新操作
	_, err := dao.db.ExecContext(ctx, query, code.CDA008, code.CDA009, code.CDA001)
	return err
}

// ====================================================================
// CountByBizAndInputInTimeRange 统计指定时间范围内某业务类型和输入的验证码数量
// 用于实现发送频率限制功能
// 参数：
// - ctx: 上下文，用于控制超时和取消
// - biz: 业务类型
// - input: 用户输入
// - startTime: 开始时间
// 返回值：
// - int: 符合条件的验证码数量
// - error: 操作错误，成功时为nil
// ====================================================================
func (dao *CodeDAO) CountByBizAndInputInTimeRange(ctx context.Context, biz, input string, startTime time.Time) (int, error) {
	query := `
	SELECT COUNT(*) 
	FROM CD01 
	WHERE CDA002 = ? AND CDA003 = ? AND CDA005 >= ?
	`

	// 执行SQL查询并获取计数结果
	var count int
	err := dao.db.QueryRowContext(ctx, query, biz, input, startTime).Scan(&count)
	return count, err
}

// ====================================================================
// CleanExpiredCodes 清理过期的验证码
// 删除已过期且状态为已过期或已失效的验证码记录
// 参数：
// - ctx: 上下文，用于控制超时和取消
// 返回值：
// - int64: 被删除的记录数量
// - error: 操作错误，成功时为nil
// ====================================================================
func (dao *CodeDAO) CleanExpiredCodes(ctx context.Context) (int64, error) {
	query := `
	DELETE FROM CD01 
	WHERE CDA006 < NOW() AND CDA009 IN (3, 4)
	`

	// 执行SQL删除操作
	result, err := dao.db.ExecContext(ctx, query)
	if err != nil {
		return 0, err
	}

	// 返回被删除的记录数量
	return result.RowsAffected()
}
