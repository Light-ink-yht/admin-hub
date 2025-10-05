package dao

import (
	"time"

	"github.com/Light-ink-yht/admin-hub/pkg/model"
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
