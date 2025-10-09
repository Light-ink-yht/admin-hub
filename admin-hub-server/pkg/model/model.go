package model

import "time"

// LN01 自定义 model
type LN01 struct {
	LAN001 int64     // id
	LAN002 time.Time // 创建时间
	LAN003 time.Time // 更新时间
	LAN004 time.Time // 删除时间
	LAN005 string    // 是否删除 1 否，2 是
	LAN006 int64     // 创建人
	LAN007 int64     // 修改人
}
