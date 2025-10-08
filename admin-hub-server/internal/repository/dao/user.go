package dao

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Light-ink-yht/admin-hub/internal/domain/user_domain"
	"github.com/go-sql-driver/mysql"
)

var (
	ErrUserDuplicateEmailOrPhone = errors.New("邮箱或者手机号冲突")
	ErrUserNotFound              = errors.New("未找到用户")
)

type UserDao interface {
	Insert(ctx context.Context, user *user_domain.AA01) error
	FindByEmail(ctx context.Context, aaa002 string) (*user_domain.AA01, error)
	FindByPhone(ctx context.Context, aaa003 string) (*user_domain.AA01, error)
	FindByUserId(ctx context.Context, aaa001 int64) (*user_domain.AA01, error)
}

type userDao struct {
	db *sql.DB
}

func NewUserDao(db *sql.DB) UserDao {
	return &userDao{
		db: db,
	}
}

func (dao *userDao) Insert(ctx context.Context, user *user_domain.AA01) error {
	query := `
			INSERT INTO aa01 (
				lan001, lan002, lan003, lan004, lan005, lan006, lan007,
				aaa001, aaa002, aaa003, aaa004, aaa005, aaa006, aaa007,
				aaa008, aaa009, aaa010, aaa011, aaa012, aaa013, aaa014, aaa015
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`
	_, err := dao.db.ExecContext(ctx, query,
		user.LAN001, user.LAN002, user.LAN003, user.LAN004, user.LAN005, user.LAN006, user.LAN007,
		user.AAA001, user.AAA002, user.AAA003, user.AAA004, user.AAA005, user.AAA006, user.AAA007,
		user.AAA008, user.AAA009, user.AAA010, user.AAA011, user.AAA012, user.AAA013, user.AAA014, user.AAA015,
	)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			const uniqueConflictsErrNo uint16 = 1062
			if mysqlErr.Number == uniqueConflictsErrNo {
				// 邮箱或者手机号冲突
				return ErrUserDuplicateEmailOrPhone
			}
		}
	}
	return nil
}

func (dao *userDao) FindByEmail(ctx context.Context, aaa002 string) (*user_domain.AA01, error) {
	query := `
		SELECT 
			lan001, lan002, lan003, lan004, lan005, lan006, lan007,
			aaa001, aaa002, aaa003, aaa004, aaa005, aaa006, aaa007,
			aaa008, aaa009, aaa010, aaa011, aaa012, aaa013, aaa014, aaa015
		FROM aa01 
		WHERE aaa002 = ? AND lan005 = '1'
		LIMIT 1
	`

	row := dao.db.QueryRowContext(ctx, query, aaa002)

	var user user_domain.AA01
	err := row.Scan(
		&user.LAN001, &user.LAN002, &user.LAN003, &user.LAN004, &user.LAN005, &user.LAN006, &user.LAN007,
		&user.AAA001, &user.AAA002, &user.AAA003, &user.AAA004, &user.AAA005, &user.AAA006, &user.AAA007,
		&user.AAA008, &user.AAA009, &user.AAA010, &user.AAA011, &user.AAA012, &user.AAA013, &user.AAA014, &user.AAA015,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (dao *userDao) FindByPhone(ctx context.Context, aaa003 string) (*user_domain.AA01, error) {
	query := `
		SELECT 
			lan001, lan002, lan003, lan004, lan005, lan006, lan007,
			aaa001, aaa002, aaa003, aaa004, aaa005, aaa006, aaa007,
			aaa008, aaa009, aaa010, aaa011, aaa012, aaa013, aaa014, aaa015
		FROM aa01 
		WHERE aaa003 = ? AND lan005 = '1'
		LIMIT 1
	`

	row := dao.db.QueryRowContext(ctx, query, aaa003)

	var user user_domain.AA01
	err := row.Scan(
		&user.LAN001, &user.LAN002, &user.LAN003, &user.LAN004, &user.LAN005, &user.LAN006, &user.LAN007,
		&user.AAA001, &user.AAA002, &user.AAA003, &user.AAA004, &user.AAA005, &user.AAA006, &user.AAA007,
		&user.AAA008, &user.AAA009, &user.AAA010, &user.AAA011, &user.AAA012, &user.AAA013, &user.AAA014, &user.AAA015,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (dao *userDao) FindByUserId(ctx context.Context, aaa001 int64) (*user_domain.AA01, error) {
	query := `
		SELECT 
			lan001, lan002, lan003, lan004, lan005, lan006, lan007,
			aaa001, aaa002, aaa003, aaa004, aaa005, aaa006, aaa007,
			aaa008, aaa009, aaa010, aaa011, aaa012, aaa013, aaa014, aaa015
		FROM aa01 
		WHERE aaa001 = ? AND lan005 = '1'
		LIMIT 1
	`

	row := dao.db.QueryRowContext(ctx, query, aaa001)

	var user user_domain.AA01
	err := row.Scan(
		&user.LAN001, &user.LAN002, &user.LAN003, &user.LAN004, &user.LAN005, &user.LAN006, &user.LAN007,
		&user.AAA001, &user.AAA002, &user.AAA003, &user.AAA004, &user.AAA005, &user.AAA006, &user.AAA007,
		&user.AAA008, &user.AAA009, &user.AAA010, &user.AAA011, &user.AAA012, &user.AAA013, &user.AAA014, &user.AAA015,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}
