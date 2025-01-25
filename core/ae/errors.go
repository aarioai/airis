package ae

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"github.com/aarioai/airis/pkg/utils"
	"github.com/redis/go-redis/v9"
	"regexp"
)

const (
	sqlBadConnMsg   = "sql bad conn: "
	sqlSkipMsg      = "sql skip: "
	sqlRemoveArgMsg = "sql remove argument: "
	sqlConnDoneMsg  = "sql conn done: "
	sqlTxDoneMsg    = "sql tx done: "
	sqlErrorMsg     = "sql error: "
)

var (
	duplicateKeyPattern = regexp.MustCompile(`Duplicate\s+entry\s+'([^']*)'\s+for\s+key\s+'([^']*)'`)
)

func FirstError(errs ...error) error {
	for _, e := range errs {
		if e != nil {
			return e
		}
	}
	return nil
}
func IsNotFound(err error) bool {
	return errors.Is(err, sql.ErrNoRows) || errors.Is(err, redis.Nil) || errors.Is(err, ErrNotFound)
}

// NewSQLError 处理 SQL 错误
func NewSQLError(err error, details ...any) *Error {
	if err == nil {
		return nil
	}
	msg := err.Error()
	caller := utils.Caller(1)

	errorMapping := map[error]func() *Error{
		driver.ErrBadConn:        func() *Error { return NewE(caller + sqlBadConnMsg + msg).WithDetail(details...) },
		driver.ErrSkip:           func() *Error { return NewE(caller + sqlSkipMsg + msg).WithDetail(details...) },
		driver.ErrRemoveArgument: func() *Error { return NewE(caller + sqlRemoveArgMsg + msg).WithDetail(details...) },
		sql.ErrNoRows:            func() *Error { return ErrorNotFound }, // can't WithDetail, locked
		sql.ErrConnDone:          func() *Error { return NewE(caller + sqlConnDoneMsg + msg).WithDetail(details...) },
		sql.ErrTxDone:            func() *Error { return NewE(caller + sqlTxDoneMsg + msg).WithDetail(details...) },
	}

	for errType, handler := range errorMapping {
		if errors.Is(err, errType) {
			return handler()
		}
	}

	// 处理重复键错误
	if matches := duplicateKeyPattern.FindStringSubmatch(msg); len(matches) == 3 {
		return NewConflict("sql key").WithDetail(details...)
	}

	return NewE(caller + sqlErrorMsg + msg).WithDetail(details...)
}

// NewRedisError 处理 Redis 错误
// @TODO
func NewRedisError(err error, details ...any) *Error {
	if err == nil {
		return nil
	}

	if errors.Is(err, redis.Nil) {
		return New(NotFound, "redis key not found").WithDetail(details...)
	}
	msg := err.Error()
	caller := utils.Caller(1)
	return New(InternalServerError, caller+" redis: "+msg).WithDetail(details...)
}
