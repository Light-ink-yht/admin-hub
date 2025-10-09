package dao

import (
	"context"
	"database/sql"
	"time"

	"github.com/Light-ink-yht/admin-hub/internal/domain/sms_domain"
)

// SmsConfigDAO 短信服务商配置数据访问对象
// 负责与数据库交互，执行短信服务商配置的增删改查操作
// 包含创建表、插入、查询、更新状态等方法
type SmsConfigDAO struct {
	db *sql.DB // 数据库连接对象
}

// NewSmsConfigDAO 创建新的短信服务商配置DAO实例
// 参数：
// - db: 数据库连接对象
// 返回值：
// - *SmsConfigDAO: 短信服务商配置DAO实例
func NewSmsConfigDAO(db *sql.DB) *SmsConfigDAO {
	return &SmsConfigDAO{
		db: db,
	}
}

// CreateSmsConfigTable 创建短信服务商配置表
// 在数据库中创建SM01短信服务商配置表，包含所需的字段和索引
// 参数：
// - ctx: 上下文，用于控制超时和取消
// 返回值：
// - error: 操作错误，成功时为nil
func (dao *SmsConfigDAO) CreateSmsConfigTable(ctx context.Context) error {
	query := `
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
	`

	// 执行SQL创建表
	_, err := dao.db.ExecContext(ctx, query)
	return err
}

// CreateSmsTemplateTable 创建短信模板表
// 在数据库中创建SM02短信模板表，包含所需的字段和索引
// 参数：
// - ctx: 上下文，用于控制超时和取消
// 返回值：
// - error: 操作错误，成功时为nil
func (dao *SmsConfigDAO) CreateSmsTemplateTable(ctx context.Context) error {
	query := `
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
	`

	// 执行SQL创建表
	_, err := dao.db.ExecContext(ctx, query)
	return err
}

// InsertConfig 插入短信服务商配置记录
// 将短信服务商配置领域模型数据插入到数据库中
// 使用ON DUPLICATE KEY UPDATE确保在唯一键冲突时更新现有记录
// 参数：
// - ctx: 上下文，用于控制超时和取消
// - config: 短信服务商配置领域模型对象
// 返回值：
// - error: 操作错误，成功时为nil
func (dao *SmsConfigDAO) InsertConfig(ctx context.Context, config *sms_domain.SmsConfig) error {
	query := `
	INSERT INTO SM01 (SMA001, SMA002, SMA003, SMA004, SMA005, SMA006, SMA007, SMA008, SMA009, SMA010, SMA011)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON DUPLICATE KEY UPDATE 
		SMA002 = VALUES(SMA002),
		SMA003 = VALUES(SMA003),
		SMA004 = VALUES(SMA004),
		SMA005 = VALUES(SMA005),
		SMA006 = VALUES(SMA006),
		SMA007 = VALUES(SMA007),
		SMA008 = VALUES(SMA008),
		SMA009 = VALUES(SMA009),
		SMA010 = VALUES(SMA010),
		SMA011 = VALUES(SMA011)
	`

	// 执行SQL插入或更新操作
	_, err := dao.db.ExecContext(
		ctx,
		query,
		config.SMA001, // 主键ID
		config.SMA002, // 配置名称
		config.SMA003, // 服务商类型
		config.SMA004, // 应用ID
		config.SMA005, // 应用密钥
		config.SMA006, // 应用密钥
		config.SMA007, // 短信签名
		config.SMA008, // 配置状态
		config.SMA009, // 创建时间
		config.SMA010, // 修改时间
		config.SMA011, // 备注信息
	)
	return err
}

// UpdateConfig 更新短信服务商配置信息
// 更新短信服务商配置的各项信息
// 参数：
// - ctx: 上下文，用于控制超时和取消
// - config: 包含更新信息的短信服务商配置领域模型对象
// 返回值：
// - error: 操作错误，成功时为nil
func (dao *SmsConfigDAO) UpdateConfig(ctx context.Context, config *sms_domain.SmsConfig) error {
	query := `
	UPDATE SM01 
	SET SMA002 = ?, SMA003 = ?, SMA004 = ?, SMA005 = ?, SMA006 = ?, SMA007 = ?, SMA008 = ?, SMA010 = ?, SMA011 = ? 
	WHERE SMA001 = ?
	`

	// 执行SQL更新操作
	_, err := dao.db.ExecContext(
		ctx,
		query,
		config.SMA002, // 配置名称
		config.SMA003, // 服务商类型
		config.SMA004, // 应用ID
		config.SMA005, // 应用密钥
		config.SMA006, // 应用密钥
		config.SMA007, // 短信签名
		config.SMA008, // 配置状态
		config.SMA010, // 修改时间
		config.SMA011, // 备注信息
		config.SMA001, // 主键ID
	)
	return err
}

// DeleteConfig 删除短信服务商配置
// 根据ID删除短信服务商配置记录
// 参数：
// - ctx: 上下文，用于控制超时和取消
// - id: 短信服务商配置ID
// 返回值：
// - error: 操作错误，成功时为nil
func (dao *SmsConfigDAO) DeleteConfig(ctx context.Context, id string) error {
	query := `
	DELETE FROM SM01 WHERE SMA001 = ?
	`

	// 执行SQL删除操作
	_, err := dao.db.ExecContext(ctx, query, id)
	return err
}

// FindConfigByID 根据ID查询短信服务商配置
// 查询指定ID的短信服务商配置信息
// 参数：
// - ctx: 上下文，用于控制超时和取消
// - id: 短信服务商配置ID
// 返回值：
// - *sms_domain.SmsConfig: 短信服务商配置领域模型对象，未找到时为nil
// - error: 操作错误，成功时为nil
func (dao *SmsConfigDAO) FindConfigByID(ctx context.Context, id string) (*sms_domain.SmsConfig, error) {
	query := `
	SELECT SMA001, SMA002, SMA003, SMA004, SMA005, SMA006, SMA007, SMA008, SMA009, SMA010, SMA011
	FROM SM01 
	WHERE SMA001 = ?
	`

	// 初始化短信服务商配置对象
	var config sms_domain.SmsConfig
	// 执行SQL查询并扫描结果
	err := dao.db.QueryRowContext(ctx, query, id).Scan(
		&config.SMA001,
		&config.SMA002,
		&config.SMA003,
		&config.SMA004,
		&config.SMA005,
		&config.SMA006,
		&config.SMA007,
		&config.SMA008,
		&config.SMA009,
		&config.SMA010,
		&config.SMA011,
	)

	// 处理未找到记录的情况
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &config, err
}

// FindAllConfigs 查询所有短信服务商配置
// 查询所有短信服务商配置信息
// 参数：
// - ctx: 上下文，用于控制超时和取消
// 返回值：
// - []*sms_domain.SmsConfig: 短信服务商配置领域模型对象列表
// - error: 操作错误，成功时为nil
func (dao *SmsConfigDAO) FindAllConfigs(ctx context.Context) ([]*sms_domain.SmsConfig, error) {
	query := `
	SELECT SMA001, SMA002, SMA003, SMA004, SMA005, SMA006, SMA007, SMA008, SMA009, SMA010, SMA011
	FROM SM01 
	ORDER BY SMA009 DESC
	`

	// 执行SQL查询
	rows, err := dao.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 初始化结果列表
	var configs []*sms_domain.SmsConfig
	// 遍历查询结果
	for rows.Next() {
		var config sms_domain.SmsConfig
		// 扫描结果到对象
		err := rows.Scan(
			&config.SMA001,
			&config.SMA002,
			&config.SMA003,
			&config.SMA004,
			&config.SMA005,
			&config.SMA006,
			&config.SMA007,
			&config.SMA008,
			&config.SMA009,
			&config.SMA010,
			&config.SMA011,
		)
		if err != nil {
			return nil, err
		}
		configs = append(configs, &config)
	}

	return configs, nil
}

// FindEnabledConfigs 查询所有启用的短信服务商配置
// 查询所有启用的短信服务商配置信息
// 参数：
// - ctx: 上下文，用于控制超时和取消
// 返回值：
// - []*sms_domain.SmsConfig: 短信服务商配置领域模型对象列表
// - error: 操作错误，成功时为nil
func (dao *SmsConfigDAO) FindEnabledConfigs(ctx context.Context) ([]*sms_domain.SmsConfig, error) {
	query := `
	SELECT SMA001, SMA002, SMA003, SMA004, SMA005, SMA006, SMA007, SMA008, SMA009, SMA010, SMA011
	FROM SM01 
	WHERE SMA008 = 1
	ORDER BY SMA009 DESC
	`

	// 执行SQL查询
	rows, err := dao.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 初始化结果列表
	var configs []*sms_domain.SmsConfig
	// 遍历查询结果
	for rows.Next() {
		var config sms_domain.SmsConfig
		// 扫描结果到对象
		err := rows.Scan(
			&config.SMA001,
			&config.SMA002,
			&config.SMA003,
			&config.SMA004,
			&config.SMA005,
			&config.SMA006,
			&config.SMA007,
			&config.SMA008,
			&config.SMA009,
			&config.SMA010,
			&config.SMA011,
		)
		if err != nil {
			return nil, err
		}
		configs = append(configs, &config)
	}

	return configs, nil
}

// EnableConfig 启用短信服务商配置
// 根据ID启用短信服务商配置
// 参数：
// - ctx: 上下文，用于控制超时和取消
// - id: 短信服务商配置ID
// 返回值：
// - error: 操作错误，成功时为nil
func (dao *SmsConfigDAO) EnableConfig(ctx context.Context, id string) error {
	query := `
	UPDATE SM01 SET SMA008 = 1, SMA010 = ? WHERE SMA001 = ?
	`

	// 执行SQL更新操作
	_, err := dao.db.ExecContext(ctx, query, time.Now(), id)
	return err
}

// DisableConfig 禁用短信服务商配置
// 根据ID禁用短信服务商配置
// 参数：
// - ctx: 上下文，用于控制超时和取消
// - id: 短信服务商配置ID
// 返回值：
// - error: 操作错误，成功时为nil
func (dao *SmsConfigDAO) DisableConfig(ctx context.Context, id string) error {
	query := `
	UPDATE SM01 SET SMA008 = 2, SMA010 = ? WHERE SMA001 = ?
	`

	// 执行SQL更新操作
	_, err := dao.db.ExecContext(ctx, query, time.Now(), id)
	return err
}

// InsertTemplate 插入短信模板记录
// 将短信模板领域模型数据插入到数据库中
// 参数：
// - ctx: 上下文，用于控制超时和取消
// - template: 短信模板领域模型对象
// 返回值：
// - error: 操作错误，成功时为nil
func (dao *SmsConfigDAO) InsertTemplate(ctx context.Context, template *sms_domain.SmsTemplate) error {
	query := `
	INSERT INTO SM02 (SMB001, SMB002, SMB003, SMB004, SMB005, SMB006, SMB007, SMB008, SMB009)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	// 执行SQL插入操作
	_, err := dao.db.ExecContext(
		ctx,
		query,
		template.SMB001, // 主键ID
		template.SMB002, // 模板名称
		template.SMB003, // 模板标识
		template.SMB004, // 服务商模板ID
		template.SMB005, // 模板内容描述
		template.SMB006, // 模板状态
		template.SMB007, // 创建时间
		template.SMB008, // 修改时间
		template.SMB009, // 备注信息
	)
	return err
}

// UpdateTemplate 更新短信模板信息
// 更新短信模板的各项信息
// 参数：
// - ctx: 上下文，用于控制超时和取消
// - template: 包含更新信息的短信模板领域模型对象
// 返回值：
// - error: 操作错误，成功时为nil
func (dao *SmsConfigDAO) UpdateTemplate(ctx context.Context, template *sms_domain.SmsTemplate) error {
	query := `
	UPDATE SM02 
	SET SMB002 = ?, SMB003 = ?, SMB004 = ?, SMB005 = ?, SMB006 = ?, SMB008 = ?, SMB009 = ? 
	WHERE SMB001 = ?
	`

	// 执行SQL更新操作
	_, err := dao.db.ExecContext(
		ctx,
		query,
		template.SMB002, // 模板名称
		template.SMB003, // 模板标识
		template.SMB004, // 服务商模板ID
		template.SMB005, // 模板内容描述
		template.SMB006, // 模板状态
		template.SMB008, // 修改时间
		template.SMB009, // 备注信息
		template.SMB001, // 主键ID
	)
	return err
}

// DeleteTemplate 删除短信模板
// 根据ID删除短信模板记录
// 参数：
// - ctx: 上下文，用于控制超时和取消
// - id: 短信模板ID
// 返回值：
// - error: 操作错误，成功时为nil
func (dao *SmsConfigDAO) DeleteTemplate(ctx context.Context, id string) error {
	query := `
	DELETE FROM SM02 WHERE SMB001 = ?
	`

	// 执行SQL删除操作
	_, err := dao.db.ExecContext(ctx, query, id)
	return err
}

// FindTemplateByID 根据ID查询短信模板
// 查询指定ID的短信模板信息
// 参数：
// - ctx: 上下文，用于控制超时和取消
// - id: 短信模板ID
// 返回值：
// - *sms_domain.SmsTemplate: 短信模板领域模型对象，未找到时为nil
// - error: 操作错误，成功时为nil
func (dao *SmsConfigDAO) FindTemplateByID(ctx context.Context, id string) (*sms_domain.SmsTemplate, error) {
	query := `
	SELECT SMB001, SMB002, SMB003, SMB004, SMB005, SMB006, SMB007, SMB008, SMB009
	FROM SM02 
	WHERE SMB001 = ?
	`

	// 初始化短信模板对象
	var template sms_domain.SmsTemplate
	// 执行SQL查询并扫描结果
	err := dao.db.QueryRowContext(ctx, query, id).Scan(
		&template.SMB001,
		&template.SMB002,
		&template.SMB003,
		&template.SMB004,
		&template.SMB005,
		&template.SMB006,
		&template.SMB007,
		&template.SMB008,
		&template.SMB009,
	)

	// 处理未找到记录的情况
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &template, err
}

// FindTemplateByIdentifier 根据标识查询短信模板
// 根据模板标识查询短信模板信息
// 参数：
// - ctx: 上下文，用于控制超时和取消
// - identifier: 短信模板标识
// 返回值：
// - *sms_domain.SmsTemplate: 短信模板领域模型对象，未找到时为nil
// - error: 操作错误，成功时为nil
func (dao *SmsConfigDAO) FindTemplateByIdentifier(ctx context.Context, identifier string) (*sms_domain.SmsTemplate, error) {
	query := `
	SELECT SMB001, SMB002, SMB003, SMB004, SMB005, SMB006, SMB007, SMB008, SMB009
	FROM SM02 
	WHERE SMB003 = ?
	`

	// 初始化短信模板对象
	var template sms_domain.SmsTemplate
	// 执行SQL查询并扫描结果
	err := dao.db.QueryRowContext(ctx, query, identifier).Scan(
		&template.SMB001,
		&template.SMB002,
		&template.SMB003,
		&template.SMB004,
		&template.SMB005,
		&template.SMB006,
		&template.SMB007,
		&template.SMB008,
		&template.SMB009,
	)

	// 处理未找到记录的情况
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &template, err
}

// FindAllTemplates 查询所有短信模板
// 查询所有短信模板信息
// 参数：
// - ctx: 上下文，用于控制超时和取消
// 返回值：
// - []*sms_domain.SmsTemplate: 短信模板领域模型对象列表
// - error: 操作错误，成功时为nil
func (dao *SmsConfigDAO) FindAllTemplates(ctx context.Context) ([]*sms_domain.SmsTemplate, error) {
	query := `
	SELECT SMB001, SMB002, SMB003, SMB004, SMB005, SMB006, SMB007, SMB008, SMB009
	FROM SM02 
	ORDER BY SMB007 DESC
	`

	// 执行SQL查询
	rows, err := dao.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 初始化结果列表
	var templates []*sms_domain.SmsTemplate
	// 遍历查询结果
	for rows.Next() {
		var template sms_domain.SmsTemplate
		// 扫描结果到对象
		err := rows.Scan(
			&template.SMB001,
			&template.SMB002,
			&template.SMB003,
			&template.SMB004,
			&template.SMB005,
			&template.SMB006,
			&template.SMB007,
			&template.SMB008,
			&template.SMB009,
		)
		if err != nil {
			return nil, err
		}
		templates = append(templates, &template)
	}

	return templates, nil
}

// FindEnabledTemplates 查询所有启用的短信模板
// 查询所有启用的短信模板信息
// 参数：
// - ctx: 上下文，用于控制超时和取消
// 返回值：
// - []*sms_domain.SmsTemplate: 短信模板领域模型对象列表
// - error: 操作错误，成功时为nil
func (dao *SmsConfigDAO) FindEnabledTemplates(ctx context.Context) ([]*sms_domain.SmsTemplate, error) {
	query := `
	SELECT SMB001, SMB002, SMB003, SMB004, SMB005, SMB006, SMB007, SMB008, SMB009
	FROM SM02 
	WHERE SMB006 = 1
	ORDER BY SMB007 DESC
	`

	// 执行SQL查询
	rows, err := dao.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 初始化结果列表
	var templates []*sms_domain.SmsTemplate
	// 遍历查询结果
	for rows.Next() {
		var template sms_domain.SmsTemplate
		// 扫描结果到对象
		err := rows.Scan(
			&template.SMB001,
			&template.SMB002,
			&template.SMB003,
			&template.SMB004,
			&template.SMB005,
			&template.SMB006,
			&template.SMB007,
			&template.SMB008,
			&template.SMB009,
		)
		if err != nil {
			return nil, err
		}
		templates = append(templates, &template)
	}

	return templates, nil
}

// EnableTemplate 启用短信模板
// 根据ID启用短信模板
// 参数：
// - ctx: 上下文，用于控制超时和取消
// - id: 短信模板ID
// 返回值：
// - error: 操作错误，成功时为nil
func (dao *SmsConfigDAO) EnableTemplate(ctx context.Context, id string) error {
	query := `
	UPDATE SM02 SET SMB006 = 1, SMB008 = ? WHERE SMB001 = ?
	`

	// 执行SQL更新操作
	_, err := dao.db.ExecContext(ctx, query, time.Now(), id)
	return err
}

// DisableTemplate 禁用短信模板
// 根据ID禁用短信模板
// 参数：
// - ctx: 上下文，用于控制超时和取消
// - id: 短信模板ID
// 返回值：
// - error: 操作错误，成功时为nil
func (dao *SmsConfigDAO) DisableTemplate(ctx context.Context, id string) error {
	query := `
	UPDATE SM02 SET SMB006 = 2, SMB008 = ? WHERE SMB001 = ?
	`

	// 执行SQL更新操作
	_, err := dao.db.ExecContext(ctx, query, time.Now(), id)
	return err
}
