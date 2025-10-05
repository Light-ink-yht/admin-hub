package ioc

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Light-ink-yht/admin-hub/internal/domain/email_domain"
	"github.com/Light-ink-yht/admin-hub/internal/repository/dao"
	"github.com/google/uuid"
)

// InitEmailTables 初始化邮件相关的数据库表
// 确保邮件服务器配置表和邮件模板表被正确创建
// 如果表不存在，则创建它们；如果存在，则不做任何操作
// 参数：
// - db: 数据库连接对象
// 返回值：
// - error: 操作错误，成功时为nil
func InitEmailTables(db *sql.DB) error {
	// 创建邮件配置相关的DAO实例
	emailConfigDAO := dao.NewEmailConfigDAO(db)
	emailTemplateDAO := dao.NewEmailTemplateDAO(db)

	// 创建上下文，设置超时
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 创建邮件服务器配置表
	err := emailConfigDAO.CreateEmailConfigTable(ctx)
	if err != nil {
		return fmt.Errorf("创建邮件服务器配置表失败: %w", err)
	}
	fmt.Println("邮件服务器配置表已创建或已存在")

	// 创建邮件模板表
	err = emailTemplateDAO.CreateEmailTemplateTable(ctx)
	if err != nil {
		return fmt.Errorf("创建邮件模板表失败: %w", err)
	}
	fmt.Println("邮件模板表已创建或已存在")

	return nil
}

// InitDefaultEmailConfig 初始化默认邮件配置
// 如果没有默认的邮件配置，则创建一个带有基本设置的默认配置
// 参数：
// - db: 数据库连接对象
// 返回值：
// - error: 操作错误，成功时为nil
func InitDefaultEmailConfig(db *sql.DB) error {
	// 创建邮件配置DAO实例
	emailConfigDAO := dao.NewEmailConfigDAO(db)

	// 创建上下文
	ctx := context.Background()

	// 检查是否已存在默认配置
	defaultConfig, err := emailConfigDAO.FindDefault(ctx)
	if err != nil {
		return fmt.Errorf("查询默认邮件配置失败: %w", err)
	}

	// 如果已存在默认配置，则不进行初始化
	if defaultConfig != nil {
		fmt.Println("默认邮件配置已存在，跳过初始化")
		return nil
	}

	// 创建默认邮件配置
	defaultConfig = email_domain.NewEmailConfig(
		"1480224563@qq.com", // 发送方邮箱地址
		"kglwbdxsfosrhhhi",  // 发送方邮箱授权码
		"smtp.qq.com",       // SMTP服务器地址
		465,                 // SMTP服务器端口
		"默认邮件服务器",           // 配置名称
	)
	// 设置额外字段
	defaultConfig.EMA001 = uuid.New().String()

	// 插入默认配置
	err = emailConfigDAO.Insert(ctx, defaultConfig)
	if err != nil {
		return fmt.Errorf("初始化默认邮件配置失败: %w", err)
	}

	fmt.Println("默认邮件配置已初始化")
	return nil
}

// InitDefaultEmailTemplates 初始化默认邮件模板
// 创建一些常用的邮件模板，如验证码模板、通知模板等
// 参数：
// - db: 数据库连接对象
// 返回值：
// - error: 操作错误，成功时为nil
func InitDefaultEmailTemplates(db *sql.DB) error {
	// 创建邮件模板DAO实例
	emailTemplateDAO := dao.NewEmailTemplateDAO(db)

	// 创建上下文
	ctx := context.Background()

	// 定义要初始化的默认模板
	defaultTemplates := []*email_domain.EmailTemplate{
		// 验证码模板
		func() *email_domain.EmailTemplate {
			t := email_domain.NewEmailTemplate(
				"验证码邮件",
				"verification_code",
				"【验证码】您的验证码已生成",
				`<p>您的验证码是：<strong>{code}</strong></p>
<p>请在10分钟内使用该验证码完成验证。</p>`,
			)
			t.EMB001 = uuid.New().String()
			return t
		}(),
		// 注册成功通知模板
		func() *email_domain.EmailTemplate {
			t := email_domain.NewEmailTemplate(
				"注册成功通知",
				"register_success",
				"【注册成功】欢迎加入我们",
				`<p>尊敬的用户，</p>
<p>恭喜您成功注册我们的系统！</p>
<p>如有任何问题，请联系客服。</p>`,
			)
			t.EMB001 = uuid.New().String()
			return t
		}(),
		// 密码重置模板
		func() *email_domain.EmailTemplate {
			t := email_domain.NewEmailTemplate(
				"密码重置邮件",
				"password_reset",
				"【密码重置】您的验证码已生成",
				`<p>您正在进行密码重置操作，验证码是：<strong>{code}</strong></p>
<p>请在10分钟内使用该验证码完成密码重置。</p>
<p>如非本人操作，请忽略此邮件。</p>`,
			)
			t.EMB001 = uuid.New().String()
			return t
		}(),
	}

	// 插入每个默认模板
	for _, template := range defaultTemplates {
		// 检查是否已存在相同类型的模板
		existingTemplate, err := emailTemplateDAO.FindByType(ctx, template.EMB003)
		if err != nil {
			return fmt.Errorf("查询邮件模板失败: %w", err)
		}

		// 如果已存在相同类型的模板，则跳过
		if existingTemplate != nil {
			fmt.Printf("邮件模板类型 %s 已存在，跳过初始化\n", template.EMB003)
			continue
		}

		// 插入模板
		err = emailTemplateDAO.Insert(ctx, template)
		if err != nil {
			return fmt.Errorf("初始化邮件模板 %s 失败: %w", template.EMB002, err)
		}

		fmt.Printf("邮件模板 %s 已初始化\n", template.EMB002)
	}

	return nil
}
