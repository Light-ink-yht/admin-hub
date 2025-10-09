package dao

import (
	"context"
	"database/sql"
	"time"

	"github.com/Light-ink-yht/admin-hub/internal/domain/email_domain"
)

// EmailConfigDAO 邮件服务器配置数据访问对象
// 负责与数据库交互，执行邮件服务器配置的增删改查操作
// 包含创建表、插入、查询、更新状态等方法
// 主要功能：
// 1. 创建邮件服务器配置表结构
// 2. 插入新的邮件服务器配置记录
// 3. 根据ID查询邮件服务器配置
// 4. 查询所有邮件服务器配置
// 5. 更新邮件服务器配置
// 6. 删除邮件服务器配置
// 7. 启用/禁用邮件服务器配置
// 8. 设置/获取默认邮件服务器配置

type EmailConfigDAO struct {
	db *sql.DB // 数据库连接对象
}

// NewEmailConfigDAO 创建新的邮件服务器配置DAO实例
// 参数：
// - db: 数据库连接对象
// 返回值：
// - *EmailConfigDAO: 邮件服务器配置DAO实例
func NewEmailConfigDAO(db *sql.DB) *EmailConfigDAO {
	return &EmailConfigDAO{
		db: db,
	}
}

// CreateEmailConfigTable 创建邮件服务器配置表
// 在数据库中创建EM01邮件服务器配置表，包含所需的字段和索引
// 参数：
// - ctx: 上下文，用于控制超时和取消
// 返回值：
// - error: 操作错误，成功时为nil
func (dao *EmailConfigDAO) CreateEmailConfigTable(ctx context.Context) error {
	query := `
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
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='邮件服务器配置表';
	`

	// 执行SQL创建表
	_, err := dao.db.ExecContext(ctx, query)
	return err
}

// Insert 插入邮件服务器配置记录
// 将邮件服务器配置领域模型数据插入到数据库中
// 使用ON DUPLICATE KEY UPDATE确保在唯一键冲突时更新现有记录
// 参数：
// - ctx: 上下文，用于控制超时和取消
// - config: 邮件服务器配置领域模型对象
// 返回值：
// - error: 操作错误，成功时为nil
func (dao *EmailConfigDAO) Insert(ctx context.Context, config *email_domain.EmailConfig) error {
	query := `
	INSERT INTO EM01 (EMA001, EMA002, EMA003, EMA004, EMA005, EMA006, EMA007, EMA008, EMA009, EMA010)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON DUPLICATE KEY UPDATE 
		EMA002 = VALUES(EMA002),
		EMA003 = VALUES(EMA003),
		EMA004 = VALUES(EMA004),
		EMA005 = VALUES(EMA005),
		EMA006 = VALUES(EMA006),
		EMA007 = VALUES(EMA007),
		EMA008 = VALUES(EMA008),
		EMA009 = VALUES(EMA009),
		EMA010 = VALUES(EMA010)
	`

	// 执行SQL插入或更新操作
	_, err := dao.db.ExecContext(
		ctx,
		query,
		config.EMA001, // 主键ID
		config.EMA002, // 发送方邮箱地址
		config.EMA003, // 发送方邮箱授权码
		config.EMA004, // SMTP服务器地址
		config.EMA005, // SMTP服务器端口
		config.EMA006, // 配置名称
		config.EMA007, // 配置状态
		config.EMA008, // 创建时间
		config.EMA009, // 修改时间
		config.EMA010, // 备注信息
	)
	return err
}

// Update 更新邮件服务器配置信息
// 更新邮件服务器配置的各项信息
// 参数：
// - ctx: 上下文，用于控制超时和取消
// - config: 包含更新信息的邮件服务器配置领域模型对象
// 返回值：
// - error: 操作错误，成功时为nil
func (dao *EmailConfigDAO) Update(ctx context.Context, config *email_domain.EmailConfig) error {
	query := `
	UPDATE EM01 
	SET EMA002 = ?, EMA003 = ?, EMA004 = ?, EMA005 = ?, EMA006 = ?, EMA007 = ?, EMA009 = ?, EMA010 = ? 
	WHERE EMA001 = ?
	`

	// 执行SQL更新操作
	_, err := dao.db.ExecContext(
		ctx,
		query,
		config.EMA002, // 发送方邮箱地址
		config.EMA003, // 发送方邮箱授权码
		config.EMA004, // SMTP服务器地址
		config.EMA005, // SMTP服务器端口
		config.EMA006, // 配置名称
		config.EMA007, // 配置状态
		config.EMA009, // 修改时间
		config.EMA010, // 备注信息
		config.EMA001, // 主键ID
	)
	return err
}

// Delete 删除邮件服务器配置
// 根据ID删除邮件服务器配置记录
// 参数：
// - ctx: 上下文，用于控制超时和取消
// - id: 邮件服务器配置ID
// 返回值：
// - error: 操作错误，成功时为nil
func (dao *EmailConfigDAO) Delete(ctx context.Context, id string) error {
	query := `
	DELETE FROM EM01 WHERE EMA001 = ?
	`

	// 执行SQL删除操作
	_, err := dao.db.ExecContext(ctx, query, id)
	return err
}

// FindByID 根据ID查询邮件服务器配置
// 查询指定ID的邮件服务器配置信息
// 参数：
// - ctx: 上下文，用于控制超时和取消
// - id: 邮件服务器配置ID
// 返回值：
// - *email_domain.EmailConfig: 邮件服务器配置领域模型对象，未找到时为nil
// - error: 操作错误，成功时为nil
func (dao *EmailConfigDAO) FindByID(ctx context.Context, id string) (*email_domain.EmailConfig, error) {
	query := `
	SELECT EMA001, EMA002, EMA003, EMA004, EMA005, EMA006, EMA007, EMA008, EMA009, EMA010
	FROM EM01 
	WHERE EMA001 = ?
	`

	// 初始化邮件服务器配置对象
	var config email_domain.EmailConfig
	// 执行SQL查询并扫描结果
	err := dao.db.QueryRowContext(ctx, query, id).Scan(
		&config.EMA001,
		&config.EMA002,
		&config.EMA003,
		&config.EMA004,
		&config.EMA005,
		&config.EMA006,
		&config.EMA007,
		&config.EMA008,
		&config.EMA009,
		&config.EMA010,
	)

	// 处理未找到记录的情况
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &config, err
}

// FindAll 查询所有邮件服务器配置
// 查询所有邮件服务器配置信息
// 参数：
// - ctx: 上下文，用于控制超时和取消
// 返回值：
// - []*email_domain.EmailConfig: 邮件服务器配置领域模型对象数组，未找到时为空数组
// - error: 操作错误，成功时为nil
func (dao *EmailConfigDAO) FindAll(ctx context.Context) ([]*email_domain.EmailConfig, error) {
	query := `
	SELECT EMA001, EMA002, EMA003, EMA004, EMA005, EMA006, EMA007, EMA008, EMA009, EMA010
	FROM EM01 
	ORDER BY EMA009 DESC, EMA010 DESC
	`

	// 执行SQL查询
	rows, err := dao.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 初始化结果数组
	configs := make([]*email_domain.EmailConfig, 0)
	// 遍历查询结果
	for rows.Next() {
		var config email_domain.EmailConfig
		// 扫描结果到配置对象
		if err := rows.Scan(
			&config.EMA001,
			&config.EMA002,
			&config.EMA003,
			&config.EMA004,
			&config.EMA005,
			&config.EMA006,
			&config.EMA007,
			&config.EMA008,
			&config.EMA009,
			&config.EMA010,
		); err != nil {
			return nil, err
		}
		// 添加到结果数组
		configs = append(configs, &config)
	}

	// 检查遍历过程中是否有错误
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return configs, nil
}

// Enable 启用邮件服务器配置
// 根据ID启用邮件服务器配置
// 参数：
// - ctx: 上下文，用于控制超时和取消
// - id: 邮件服务器配置ID
// 返回值：
// - error: 操作错误，成功时为nil
func (dao *EmailConfigDAO) Enable(ctx context.Context, id string) error {
	query := `
	UPDATE EM01 
	SET EMA007 = 1, EMA009 = ? 
	WHERE EMA001 = ?
	`

	// 执行SQL更新操作
	_, err := dao.db.ExecContext(ctx, query, time.Now(), id)
	return err
}

// Disable 禁用邮件服务器配置
// 根据ID禁用邮件服务器配置
// 参数：
// - ctx: 上下文，用于控制超时和取消
// - id: 邮件服务器配置ID
// 返回值：
// - error: 操作错误，成功时为nil
func (dao *EmailConfigDAO) Disable(ctx context.Context, id string) error {
	query := `
	UPDATE EM01 
	SET EMA007 = 2, EMA009 = ? 
	WHERE EMA001 = ?
	`

	// 执行SQL更新操作
	_, err := dao.db.ExecContext(ctx, query, time.Now(), id)
	return err
}

// FindDefault 查询默认邮件服务器配置
// 查询系统默认的邮件服务器配置
// 参数：
// - ctx: 上下文，用于控制超时和取消
// 返回值：
// - *email_domain.EmailConfig: 邮件服务器配置领域模型对象，未找到时为nil
// - error: 操作错误，成功时为nil
func (dao *EmailConfigDAO) FindDefault(ctx context.Context) (*email_domain.EmailConfig, error) {
	query := `
	SELECT EMA001, EMA002, EMA003, EMA004, EMA005, EMA006, EMA007, EMA008, EMA009, EMA010
	FROM EM01 
	WHERE EMA010 = 'default'
	`

	// 初始化邮件服务器配置对象
	var config email_domain.EmailConfig
	// 执行SQL查询并扫描结果
	err := dao.db.QueryRowContext(ctx, query).Scan(
		&config.EMA001,
		&config.EMA002,
		&config.EMA003,
		&config.EMA004,
		&config.EMA005,
		&config.EMA006,
		&config.EMA007,
		&config.EMA008,
		&config.EMA009,
		&config.EMA010,
	)

	// 处理未找到记录的情况
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &config, err
}

// SetDefault 设置默认邮件服务器配置
// 将指定ID的邮件服务器配置设置为默认，并取消其他配置的默认状态
// 参数：
// - ctx: 上下文，用于控制超时和取消
// - id: 邮件服务器配置ID
// 返回值：
// - error: 操作错误，成功时为nil
func (dao *EmailConfigDAO) SetDefault(ctx context.Context, id string) error {
	// 开始事务
	tx, err := dao.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 取消所有配置的默认状态
	_, err = tx.ExecContext(ctx, "UPDATE EM01 SET EMA010 = NULL WHERE EMA010 = 'default'")
	if err != nil {
		tx.Rollback()
		return err
	}

	// 设置指定配置为默认
	now := time.Now()
	_, err = tx.ExecContext(ctx, "UPDATE EM01 SET EMA010 = 'default', EMA009 = ? WHERE EMA001 = ?", now, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 提交事务
	return tx.Commit()
}

// EmailTemplateDAO 邮件模板数据访问对象
// 负责与数据库交互，执行邮件模板的增删改查操作
// 包含创建表、插入、查询、更新等方法
// 主要功能：
// 1. 创建邮件模板表结构
// 2. 插入新的邮件模板记录
// 3. 根据ID查询邮件模板
// 4. 查询所有邮件模板
// 5. 根据类型查询邮件模板
// 6. 更新邮件模板
// 7. 删除邮件模板
// 8. 搜索邮件模板

type EmailTemplateDAO struct {
	db *sql.DB // 数据库连接对象
}

// NewEmailTemplateDAO 创建新的邮件模板DAO实例
// 参数：
// - db: 数据库连接对象
// 返回值：
// - *EmailTemplateDAO: 邮件模板DAO实例
func NewEmailTemplateDAO(db *sql.DB) *EmailTemplateDAO {
	return &EmailTemplateDAO{
		db: db,
	}
}

// CreateEmailTemplateTable 创建邮件模板表
// 在数据库中创建EM02邮件模板表，包含所需的字段和索引
// 参数：
// - ctx: 上下文，用于控制超时和取消
// 返回值：
// - error: 操作错误，成功时为nil
func (dao *EmailTemplateDAO) CreateEmailTemplateTable(ctx context.Context) error {
	query := `
	CREATE TABLE IF NOT EXISTS EM02 (
		-- 主键ID，用于唯一标识邮件模板记录
		EMB001 VARCHAR(36) PRIMARY KEY COMMENT '主键ID',
		-- 模板名称，用于标识邮件模板
		EMB002 VARCHAR(100) NOT NULL COMMENT '模板名称',
		-- 模板类型，用于区分不同用途的模板（如验证码、通知、营销等）
		EMB003 VARCHAR(50) NOT NULL COMMENT '模板类型(type)',
		-- 邮件主题
		EMB004 VARCHAR(255) NOT NULL COMMENT '邮件主题',
		-- 邮件内容（HTML格式）
		EMB005 TEXT NOT NULL COMMENT '邮件内容(HTML)',
		-- 状态，标识邮件模板的当前状态
		-- 1: 启用 - 邮件模板可用
		-- 0: 禁用 - 邮件模板不可用
		EMB006 TINYINT NOT NULL DEFAULT 1 COMMENT '状态(1:启用, 0:禁用)',
		-- 创建时间，记录邮件模板的创建时间
		EMB007 DATETIME NOT NULL COMMENT '创建时间',
		-- 更新时间，记录邮件模板的最后更新时间
		EMB008 DATETIME NOT NULL COMMENT '更新时间',
		-- 备注信息，用于存储额外的说明或上下文信息
		EMB009 VARCHAR(255) DEFAULT NULL COMMENT '备注信息',
		-- 唯一索引，确保模板名称和类型的组合唯一
		UNIQUE KEY uk_name_type (EMB002, EMB003),
		-- 索引，优化按模板类型查询的性能
		INDEX idx_template_type (EMB003),
		-- 索引，优化按状态查询的性能
		INDEX idx_status (EMB006),
		-- 索引，优化按创建时间查询的性能
		INDEX idx_created_at (EMB007)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='邮件模板表';
	`

	// 执行SQL创建表
	_, err := dao.db.ExecContext(ctx, query)
	return err
}

// Insert 插入邮件模板记录
// 将邮件模板领域模型数据插入到数据库中
// 使用ON DUPLICATE KEY UPDATE确保在唯一键冲突时更新现有记录
// 参数：
// - ctx: 上下文，用于控制超时和取消
// - template: 邮件模板领域模型对象
// 返回值：
// - error: 操作错误，成功时为nil
func (dao *EmailTemplateDAO) Insert(ctx context.Context, template *email_domain.EmailTemplate) error {
	query := `
	INSERT INTO EM02 (EMB001, EMB002, EMB003, EMB004, EMB005, EMB006, EMB007, EMB008, EMB009)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON DUPLICATE KEY UPDATE 
		EMB002 = VALUES(EMB002),
		EMB003 = VALUES(EMB003),
		EMB004 = VALUES(EMB004),
		EMB005 = VALUES(EMB005),
		EMB006 = VALUES(EMB006),
		EMB007 = VALUES(EMB007),
		EMB008 = VALUES(EMB008),
		EMB009 = VALUES(EMB009)
	`

	// 执行SQL插入或更新操作
	_, err := dao.db.ExecContext(
		ctx,
		query,
		template.EMB001, // 主键ID
		template.EMB002, // 模板名称
		template.EMB003, // 模板类型
		template.EMB004, // 邮件主题
		template.EMB005, // 邮件内容
		template.EMB006, // 状态
		template.EMB007, // 创建时间
		template.EMB008, // 更新时间
		template.EMB009, // 备注信息
	)
	return err
}

// Update 更新邮件模板信息
// 更新邮件模板的各项信息
// 参数：
// - ctx: 上下文，用于控制超时和取消
// - template: 包含更新信息的邮件模板领域模型对象
// 返回值：
// - error: 操作错误，成功时为nil
func (dao *EmailTemplateDAO) Update(ctx context.Context, template *email_domain.EmailTemplate) error {
	query := `
	UPDATE EM02 
	SET EMB002 = ?, EMB003 = ?, EMB004 = ?, EMB005 = ?, EMB006 = ?, EMB008 = ?, EMB009 = ? 
	WHERE EMB001 = ?
	`

	// 执行SQL更新操作
	_, err := dao.db.ExecContext(
		ctx,
		query,
		template.EMB002, // 模板名称
		template.EMB003, // 模板类型
		template.EMB004, // 邮件主题
		template.EMB005, // 邮件内容
		template.EMB006, // 状态
		template.EMB008, // 更新时间
		template.EMB009, // 备注信息
		template.EMB001, // 主键ID
	)
	return err
}

// Delete 删除邮件模板
// 根据ID删除邮件模板记录
// 参数：
// - ctx: 上下文，用于控制超时和取消
// - id: 邮件模板ID
// 返回值：
// - error: 操作错误，成功时为nil
func (dao *EmailTemplateDAO) Delete(ctx context.Context, id string) error {
	query := `
	DELETE FROM EM02 WHERE EMB001 = ?
	`

	// 执行SQL删除操作
	_, err := dao.db.ExecContext(ctx, query, id)
	return err
}

// FindByID 根据ID查询邮件模板
// 查询指定ID的邮件模板信息
// 参数：
// - ctx: 上下文，用于控制超时和取消
// - id: 邮件模板ID
// 返回值：
// - *email_domain.EmailTemplate: 邮件模板领域模型对象，未找到时为nil
// - error: 操作错误，成功时为nil
func (dao *EmailTemplateDAO) FindByID(ctx context.Context, id string) (*email_domain.EmailTemplate, error) {
	query := `
	SELECT EMB001, EMB002, EMB003, EMB004, EMB005, EMB006, EMB007, EMB008, EMB009 
	FROM EM02 
	WHERE EMB001 = ?
	`

	// 初始化邮件模板对象
	var template email_domain.EmailTemplate
	// 执行SQL查询并扫描结果
	err := dao.db.QueryRowContext(ctx, query, id).Scan(
		&template.EMB001,
		&template.EMB002,
		&template.EMB003,
		&template.EMB004,
		&template.EMB005,
		&template.EMB006,
		&template.EMB007,
		&template.EMB008,
		&template.EMB009,
	)

	// 处理未找到记录的情况
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &template, err
}

// FindAll 查询所有邮件模板
// 查询所有邮件模板信息
// 参数：
// - ctx: 上下文，用于控制超时和取消
// 返回值：
// - []*email_domain.EmailTemplate: 邮件模板领域模型对象数组，未找到时为空数组
// - error: 操作错误，成功时为nil
func (dao *EmailTemplateDAO) FindAll(ctx context.Context) ([]*email_domain.EmailTemplate, error) {
	query := `
	SELECT EMB001, EMB002, EMB003, EMB004, EMB005, EMB006, EMB007, EMB008, EMB009 
	FROM EM02 
	ORDER BY EMB003, EMB007 DESC
	`

	// 执行SQL查询
	rows, err := dao.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 初始化结果数组
	templates := make([]*email_domain.EmailTemplate, 0)
	// 遍历查询结果
	for rows.Next() {
		var template email_domain.EmailTemplate
		// 扫描结果到模板对象
		if err := rows.Scan(
			&template.EMB001,
			&template.EMB002,
			&template.EMB003,
			&template.EMB004,
			&template.EMB005,
			&template.EMB006,
			&template.EMB007,
			&template.EMB008,
			&template.EMB009,
		); err != nil {
			return nil, err
		}
		// 添加到结果数组
		templates = append(templates, &template)
	}

	// 检查遍历过程中是否有错误
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return templates, nil
}

// FindByType 根据类型查询邮件模板
// 查询指定类型的邮件模板信息，优先返回启用状态的模板
// 参数：
// - ctx: 上下文，用于控制超时和取消
// - templateType: 模板类型
// 返回值：
// - *email_domain.EmailTemplate: 邮件模板领域模型对象，未找到时为nil
// - error: 操作错误，成功时为nil
func (dao *EmailTemplateDAO) FindByType(ctx context.Context, templateType string) (*email_domain.EmailTemplate, error) {
	query := `
	SELECT EMB001, EMB002, EMB003, EMB004, EMB005, EMB006, EMB007, EMB008, EMB009 
	FROM EM02 
	WHERE EMB003 = ? AND EMB006 = 1 
	ORDER BY EMB007 DESC 
	LIMIT 1
	`

	// 初始化邮件模板对象
	var template email_domain.EmailTemplate
	// 执行SQL查询并扫描结果
	err := dao.db.QueryRowContext(ctx, query, templateType).Scan(
		&template.EMB001,
		&template.EMB002,
		&template.EMB003,
		&template.EMB004,
		&template.EMB005,
		&template.EMB006,
		&template.EMB007,
		&template.EMB008,
		&template.EMB009,
	)

	// 处理未找到记录的情况
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &template, err
}

// Search 搜索邮件模板
// 根据关键词搜索邮件模板（匹配模板名称、类型、主题）
// 参数：
// - ctx: 上下文，用于控制超时和取消
// - keyword: 搜索关键词
// 返回值：
// - []*email_domain.EmailTemplate: 匹配的邮件模板领域模型对象数组，未找到时为空数组
// - error: 操作错误，成功时为nil
func (dao *EmailTemplateDAO) Search(ctx context.Context, keyword string) ([]*email_domain.EmailTemplate, error) {
	query := `
	SELECT EMB001, EMB002, EMB003, EMB004, EMB005, EMB006, EMB007, EMB008, EMB009 
	FROM EM02 
	WHERE EMB002 LIKE ? OR EMB003 LIKE ? OR EMB004 LIKE ? 
	ORDER BY EMB006 DESC, EMB007 DESC
	`

	// 构建搜索参数
	searchParam := "%" + keyword + "%"

	// 执行SQL查询
	rows, err := dao.db.QueryContext(ctx, query, searchParam, searchParam, searchParam)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 初始化结果数组
	templates := make([]*email_domain.EmailTemplate, 0)
	// 遍历查询结果
	for rows.Next() {
		var template email_domain.EmailTemplate
		// 扫描结果到模板对象
		if err := rows.Scan(
			&template.EMB001,
			&template.EMB002,
			&template.EMB003,
			&template.EMB004,
			&template.EMB005,
			&template.EMB006,
			&template.EMB007,
			&template.EMB008,
			&template.EMB009,
		); err != nil {
			return nil, err
		}
		// 添加到结果数组
		templates = append(templates, &template)
	}

	// 检查遍历过程中是否有错误
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return templates, nil
}
