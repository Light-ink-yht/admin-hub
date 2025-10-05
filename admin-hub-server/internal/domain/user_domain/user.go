package user_domain

import (
	"errors"
	"time"

	"github.com/Light-ink-yht/admin-hub/pkg/model"
	"github.com/dlclark/regexp2"
)

// AA01 用户模型
type AA01 struct {
	model.LN01
	AAA001 int64     // 用户id
	AAA002 string    // 邮箱 全局唯一
	AAA003 string    // 手机号 全局唯一
	AAA004 string    // 密码
	AAA005 string    // 昵称
	AAA006 string    // 姓名
	AAA007 string    // 头像
	AAA008 string    // 性别 1 男，2 女，3 未知
	AAA009 time.Time // 生日
	AAA010 string    // 状态 1 启用，2 禁用
	AAA011 string    // 备注
	AAA012 int       // 登录次数
	AAA013 time.Time // 最后登录时间
	AAA014 string    // 最后登录ip
	AAA015 time.Time // 密码修改时间
}

var (
	emailRegex                          = regexp2.MustCompile(`^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+$`, regexp2.None)
	passwordRegex                       = regexp2.MustCompile(`^(?=.*[a-zA-Z])(?=.*[0-9])(?=.*[._~!@#$^&*])[A-Za-z0-9._~!@#$^&*]{8,20}$`, regexp2.None)
	ErrTheMailboxIsNotInTheRightFormat  = errors.New("电子邮件格式无效")
	ErrThePasswordIsNotInTheRightFormat = errors.New("密码长度必须为 8-20 个字符，并包含字母、数字和特殊字符")
	ErrThePasswordIsInconsistentTwice   = errors.New("两次密码不一致")
)
