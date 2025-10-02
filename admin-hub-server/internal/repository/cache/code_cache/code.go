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

type CodeCache struct {
	client redis.Cmdable
}

func NewCodeCache(client redis.Cmdable) *CodeCache {
	return &CodeCache{client: client}
}

// Set 设置到 redis 中，biz 区分业务场景，input 用户的phone/email，code 验证码
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

func (c *CodeCache) Key(biz, input string) string {
	return fmt.Sprintf("code:%s:%s", biz, input)
}
