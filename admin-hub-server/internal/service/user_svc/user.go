package user_svc

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Light-ink-yht/admin-hub/internal/domain/user_domain"
	"github.com/Light-ink-yht/admin-hub/internal/repository/user_repo"
	"github.com/Light-ink-yht/admin-hub/pkg/snowflake"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserDuplicateEmailOrPhone = user_repo.ErrUserDuplicateEmailOrPhone
	ErrInvalidEmailOrPassword    = errors.New("邮箱或密码不对")
	ErrInvalidPhoneOrPassword    = errors.New("手机号或密码不对")
)

type UserService interface {
	Signup(ctx context.Context, req *user_domain.SignUp) error
	EmailLogin(ctx context.Context, req *user_domain.Login, userAgent string) (string, error)
	PhoneLogin(ctx context.Context, req *user_domain.Login, userAgent string) (string, error)
}

type userService struct {
	repo user_repo.UserRepo
}

func NewUserService(repo user_repo.UserRepo) UserService {
	return &userService{
		repo: repo,
	}
}

func (svc *userService) Signup(ctx context.Context, req *user_domain.SignUp) error {
	// 校验请求
	if err := req.Validate(); err != nil {
		return err
	}

	// 密码加密
	hash, err := bcrypt.GenerateFromPassword([]byte(req.AAA004), bcrypt.DefaultCost)
	if err != nil {
		// 密码加密失败
		return err
	}

	// 生成唯一ID
	userId, err := svc.GenerateUniqueID()
	if err != nil {
		return err
	}

	// 初始化用户
	user := user_domain.DefaultUser(req.AAA002, string(hash), req.AAA003, userId)
	localTime := time.Now()

	// 设置嵌入结构体的字段
	user.LAN002 = localTime
	user.LAN003 = localTime
	user.LAN004 = time.Time{}
	user.LAN005 = "1"
	user.LAN006 = 0
	user.LAN007 = 0

	// 插入用户数据
	err = svc.repo.Signup(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (svc *userService) EmailLogin(ctx context.Context, req *user_domain.Login, userAgent string) (string, error) {
	// 校验请求
	if err := req.Validate(); err != nil {
		return "", err
	}

	// 查询用户
	user, err := svc.repo.FindByEmail(ctx, req.AAA002)
	if errors.Is(err, user_repo.ErrUserNotFound) {
		return "", ErrInvalidEmailOrPassword
	}
	if err != nil {
		return "", err
	}

	// 校验密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.AAA004), []byte(req.AAA004)); err != nil {
		return "", ErrInvalidEmailOrPassword
	}

	// 生成JWT Token
	return svc.setJWTToken(userAgent, user)
}

func (svc *userService) PhoneLogin(ctx context.Context, req *user_domain.Login, userAgent string) (string, error) {
	// 校验请求
	if err := req.Validate(); err != nil {
		return "", err
	}

	// 查询用户
	user, err := svc.repo.FindByPhone(ctx, req.AAA003)
	if errors.Is(err, user_repo.ErrUserNotFound) {
		return "", ErrInvalidPhoneOrPassword
	}
	if err != nil {
		return "", err
	}

	// 校验密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.AAA004), []byte(req.AAA004)); err != nil {
		return "", ErrInvalidEmailOrPassword
	}

	// 生成JWT Token
	return svc.setJWTToken(userAgent, user)
}

// setJWTToken 生成JWT Token
func (svc *userService) setJWTToken(userAgent string, user *user_domain.AA01) (string, error) {
	// 创建用户声明
	claims := UserClaims{
		// 设置过期时间
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
		UserId:    user.AAA001,
		UserAgent: userAgent,
	}

	// 创建JWT Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	// 签名JWT Token
	tokenStr, err := token.SignedString(JWTKey)
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

// GenerateUniqueID 用于每次调用时生成一个新的唯一ID。
func (svc *userService) GenerateUniqueID() (int64, error) {
	generator, err := snowflake.NewShortIDGenerator(1) // 节点ID为1
	if err != nil {
		fmt.Println("Error creating generator:", err)
		return 0, err
	}

	id := generator.Generate()
	return id, nil
}
