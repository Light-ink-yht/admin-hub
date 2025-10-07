package snowflake

import (
	"fmt"
	"sync"
	"time"
)

// 15位ID生成器
type ShortIDGenerator struct {
	nodeID   int64 // 机器ID，最多支持1024个节点(0-1023)
	sequence int64 // 序列号
	lastTime int64 // 上次生成ID的时间戳
	mu       sync.Mutex
}

// 初始化生成器，nodeID范围0-1023
func NewShortIDGenerator(nodeID int64) (*ShortIDGenerator, error) {
	if nodeID < 0 || nodeID > 1023 {
		return nil, fmt.Errorf("nodeID must be between 0 and 1023")
	}
	return &ShortIDGenerator{
		nodeID:   nodeID,
		sequence: 0,
		lastTime: 0,
	}, nil
}

// 生成15位唯一ID
func (g *ShortIDGenerator) Generate() int64 {
	g.mu.Lock()
	defer g.mu.Unlock()

	// 取当前时间戳(秒级，从2020-01-01开始计算)
	now := time.Since(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)).Seconds()
	currentTime := int64(now)

	// 如果当前时间与上次时间相同，则递增序列号
	if currentTime == g.lastTime {
		g.sequence++
		// 序列号溢出检查(4095是12位最大数)
		if g.sequence > 4095 {
			// 等待到下一毫秒
			for currentTime <= g.lastTime {
				now := time.Since(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)).Seconds()
				currentTime = int64(now)
			}
			g.sequence = 0
		}
	} else {
		// 时间不同，重置序列号
		g.sequence = 0
	}

	g.lastTime = currentTime

	// 15位ID结构: 31位时间戳(7-8位十进制) + 10位机器ID(3位十进制) + 12位序列号(4位十进制)
	id := currentTime<<22 | (g.nodeID << 12) | g.sequence

	// 确保是15位数字，如果超过则取低15位
	if id > 999999999999999 {
		id = id % 1000000000000000
	}

	return id
}
