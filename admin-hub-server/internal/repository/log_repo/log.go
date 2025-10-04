package log_repo

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Light-ink-yht/admin-hub/internal/domain/log_domain"
	"github.com/Light-ink-yht/admin-hub/internal/repository/dao"
)

type LogRepository struct {
	db *sql.DB
}

func NewLogRepository(db *sql.DB) *LogRepository {
	return &LogRepository{
		db: db,
	}
}

// SaveCommonLog 保存常用日志到数据库
func (repo *LogRepository) SaveCommonLog(ctx context.Context, log *log_domain.LL01) error {
	_, err := repo.db.ExecContext(ctx,
		"INSERT INTO ll01 (lla001, lla002, lla003, lla004, lla005, lla006, lla007) VALUES (?, ?, ?, ?, ?, ?, ?)",
		log.LLA001, log.LLA002, log.LLA003, log.LLA004, log.LLA005, log.LLA006, log.LLA007,
	)
	return err
}

// SaveLoginLog 保存登录日志到数据库
func (repo *LogRepository) SaveLoginLog(ctx context.Context, log *log_domain.LL02) error {
	_, err := repo.db.ExecContext(ctx,
		"INSERT INTO ll02 (lla001, lla002, lla003, lla004, lla005, lla006, lla007, lla008, lla009, lla010, lla011) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		log.LLA001, log.LLA002, log.LLA003, log.LLA004, log.LLA005, log.LLA006, log.LLA007, log.LLA008, log.LLA009, log.LLA010, log.LLA011,
	)
	return err
}

// SaveBusinessLog 保存业务日志到数据库
func (repo *LogRepository) SaveBusinessLog(ctx context.Context, log *log_domain.LL03) error {
	_, err := repo.db.ExecContext(ctx,
		"INSERT INTO ll03 (lla001, lla002, lla003, lla004, lla005, lla006, lla007, lla008, lla009, lla010, lla011, lla012, lla013) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		log.LLA001, log.LLA002, log.LLA003, log.LLA004, log.LLA005, log.LLA006, log.LLA007, log.LLA008, log.LLA009, log.LLA010, log.LLA011, log.LLA012, log.LLA013,
	)
	return err
}

// GetCommonLogs 获取常用日志列表
func (repo *LogRepository) GetCommonLogs(ctx context.Context, page, pageSize int, startTime, endTime string) ([]*dao.LL01, int64, error) {
	offset := (page - 1) * pageSize
	query := "SELECT * FROM ll01 WHERE lla007 BETWEEN ? AND ? AND lan005 = '1' ORDER BY lla007 DESC LIMIT ? OFFSET ?"
	countQuery := "SELECT COUNT(*) FROM ll01 WHERE lla007 BETWEEN ? AND ? AND lan005 = '1'"

	return getLogs[dao.LL01](ctx, repo.db, query, countQuery, startTime, endTime, pageSize, offset)
}

// GetLoginLogs 获取登录日志列表
func (repo *LogRepository) GetLoginLogs(ctx context.Context, page, pageSize int, startTime, endTime string, userID string) ([]*dao.LL02, int64, error) {
	offset := (page - 1) * pageSize
	var query, countQuery string
	var args []interface{}

	if userID != "" {
		query = "SELECT * FROM ll02 WHERE lla007 BETWEEN ? AND ? AND lla004 = ? AND lan005 = '1' ORDER BY lla007 DESC LIMIT ? OFFSET ?"
		countQuery = "SELECT COUNT(*) FROM ll02 WHERE lla007 BETWEEN ? AND ? AND lla004 = ? AND lan005 = '1'"
		args = append(args, startTime, endTime, userID, pageSize, offset)
	} else {
		query = "SELECT * FROM ll02 WHERE lla007 BETWEEN ? AND ? AND lan005 = '1' ORDER BY lla007 DESC LIMIT ? OFFSET ?"
		countQuery = "SELECT COUNT(*) FROM ll02 WHERE lla007 BETWEEN ? AND ? AND lan005 = '1'"
		args = append(args, startTime, endTime, pageSize, offset)
	}

	return getLogs[dao.LL02](ctx, repo.db, query, countQuery, args...)
}

// GetBusinessLogs 获取业务日志列表
func (repo *LogRepository) GetBusinessLogs(ctx context.Context, page, pageSize int, startTime, endTime string, module, operatorID string) ([]*dao.LL03, int64, error) {
	offset := (page - 1) * pageSize
	var query, countQuery string
	var args []interface{}

	query = "SELECT * FROM ll03 WHERE lla007 BETWEEN ? AND ? AND lan005 = '1'"
	countQuery = "SELECT COUNT(*) FROM ll03 WHERE lla007 BETWEEN ? AND ? AND lan005 = '1'"
	args = append(args, startTime, endTime)

	if module != "" {
		query += " AND lla003 = ?"
		countQuery += " AND lla003 = ?"
		args = append(args, module)
	}

	if operatorID != "" {
		query += " AND lla006 = ?"
		countQuery += " AND lla006 = ?"
		args = append(args, operatorID)
	}

	query += " ORDER BY lla007 DESC LIMIT ? OFFSET ?"
	args = append(args, pageSize, offset)

	return getLogs[dao.LL03](ctx, repo.db, query, countQuery, args...)
}

// DeleteLog 删除日志
func (repo *LogRepository) DeleteLog(ctx context.Context, logType log_domain.LogType, logID string) error {
	var tableName string
	switch logType {
	case log_domain.LogTypeCommon:
		tableName = "ll01"
	case log_domain.LogTypeLogin:
		tableName = "ll02"
	case log_domain.LogTypeBusiness:
		tableName = "ll03"
	default:
		return fmt.Errorf("未知的日志类型: %s", logType)
	}

	_, err := repo.db.ExecContext(ctx,
		fmt.Sprintf("UPDATE %s SET lan005 = '2', lan004 = NOW() WHERE lla001 = ?", tableName),
		logID,
	)
	return err
}

// 通用的获取日志函数
type LogModel interface{}

func getLogs[T LogModel](ctx context.Context, db *sql.DB, query, countQuery string, args ...interface{}) ([]*T, int64, error) {
	// 获取总数
	countArgs := make([]interface{}, 0)
	for i := 0; i < len(args); i++ {
		// 跳过LIMIT和OFFSET参数
		if i == len(args)-1 || i == len(args)-2 {
			continue
		}
		countArgs = append(countArgs, args[i])
	}

	var total int64
	err := db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// 获取日志列表
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var logs []*T
	for rows.Next() {
		var log T
		// 这里需要根据具体类型实现Scan逻辑，这里简化处理
		// 在实际实现中，应该为每种类型单独实现Scan方法
		err := rows.Scan()
		if err != nil {
			return nil, 0, err
		}
		logs = append(logs, &log)
	}

	return logs, total, nil
}
