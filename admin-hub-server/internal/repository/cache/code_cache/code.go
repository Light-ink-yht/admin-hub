package code_cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var (
	ErrCodeSendTooMany   = errors.New("发送验证码太频繁")
	RrrCodeVerifyTooMany = errors.New("验证次数太多")
	ErrUnknownForCode    = errors.New("我也不知发生什么了，反正是跟 code_cache 有关")
)

// 编译器会在编译的时候，把 set_code.lua 的代码放进这个 luaSetCode 变量里面
//
//go:embed lua/set_code.lua
var luaSetCode string

//go:embed lua/verify_code.lua
var luaVerifyCode string

// 验证码计数器key的后缀
const cntKeySuffix = ":cnt"

// CodeCache 验证码缓存结构体
// 封装了与Redis相关的验证码操作
// 提供设置、验证和删除验证码的功能

type CodeCache struct {
	client redis.Cmdable
}

// NewCodeCache 创建新的验证码缓存实例
func NewCodeCache(client redis.Cmdable) *CodeCache {
	return &CodeCache{client: client}
}

// Set 设置到 redis 中，biz 区分业务场景，input 用户的phone/email，code 验证码
// 通过Lua脚本实现原子操作，保证验证码的设置和计数器的更新是原子的
func (c *CodeCache) Set(ctx context.Context, biz, input, code string) error {
	res, err := c.client.Eval(ctx, luaSetCode, []string{c.Key(biz, input)}, code).Int()
	if err != nil {
		return err
	}

	switch res {

	case 0:
		// 毫无问题
		return nil

	case 1:
		//发送太频繁
		return ErrCodeSendTooMany

	default:
		// 系统错误
		return errors.New("系统错误")
	}
}

// Verify 验证用户输入的验证码是否正确
// 通过Lua脚本实现原子操作，保证验证和计数器更新是原子的
func (c *CodeCache) Verify(ctx context.Context, biz, input, code string) (bool, error) {
	res, err := c.client.Eval(ctx, luaVerifyCode, []string{c.Key(biz, input)}, code).Int()
	if err != nil {
		return false, err
	}

	switch res {

	case 0:
		// 毫无问题
		return true, nil

	case 1:
		//发送太频繁
		return false, RrrCodeVerifyTooMany

	default:
		// 系统错误
		return false, ErrUnknownForCode
	}
}

// Delete 删除指定的验证码记录
// 同时删除验证码本身和对应的计数器
func (c *CodeCache) Delete(ctx context.Context, biz, input string) error {
	key := c.Key(biz, input)
	cntKey := key + cntKeySuffix

	// 使用事务删除两个key
	pipe := c.client.Pipeline()
	pipe.Del(ctx, key)
	pipe.Del(ctx, cntKey)
	_, err := pipe.Exec(ctx)
	return err
}

// Key 生成验证码在Redis中的key
// 格式为: code:{biz}:{input}
func (c *CodeCache) Key(biz, input string) string {
	return fmt.Sprintf("code:%s:%s", biz, input)
}
