package user_domain

import (
	"time"

	"github.com/Light-ink-yht/admin-hub/pkg/model"
)

type GetInfoResponse struct {
	model.LN01
	AAA001 int64     // 用户id
	AAA002 string    // 邮箱 全局唯一
	AAA003 string    // 手机号 全局唯一
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

func (res *GetInfoResponse) ResData(user *AA01) *GetInfoResponse {
	var resData GetInfoResponse
	resData.LN01 = model.LN01{
		LAN001: user.LAN001,
		LAN002: user.LAN002,
		LAN003: user.LAN003,
		LAN004: user.LAN004,
		LAN005: user.LAN005,
		LAN006: user.LAN006,
		LAN007: user.LAN007,
	}
	resData.AAA001 = user.AAA001
	resData.AAA002 = user.AAA002
	resData.AAA003 = user.AAA003
	resData.AAA005 = user.AAA005
	resData.AAA006 = user.AAA006
	resData.AAA007 = user.AAA007
	resData.AAA008 = user.AAA008
	resData.AAA009 = user.AAA009
	resData.AAA010 = user.AAA010
	resData.AAA011 = user.AAA011
	resData.AAA012 = user.AAA012
	resData.AAA013 = user.AAA013
	resData.AAA014 = user.AAA014
	resData.AAA015 = user.AAA015
	return &resData
}
