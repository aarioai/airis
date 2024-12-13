package ae

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"regexp"

	"github.com/redis/go-redis/v9"
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

func First(es ...*Error) *Error {
	for _, e := range es {
		if e != nil {
			return e
		}
	}
	return nil
}

// FirstError 支持所有继承 error 的子类
func FirstError(errs ...error) error {
	for _, e := range errs {
		if e != nil {
			return e
		}
	}
	return nil
}

// NewSQLError 处理 SQL 错误
func NewSQLError(err error, details ...any) *Error {
	if err == nil {
		return nil
	}
	msg, pos := CallerMsg(err.Error(), 1)

	errorMapping := map[error]func() *Error{
		driver.ErrBadConn:        func() *Error { return NewE(pos + sqlBadConnMsg + msg).WithDetail(details...) },
		driver.ErrSkip:           func() *Error { return NewE(pos + sqlSkipMsg + msg).WithDetail(details...) },
		driver.ErrRemoveArgument: func() *Error { return NewE(pos + sqlRemoveArgMsg + msg).WithDetail(details...) },
		sql.ErrNoRows:            func() *Error { return ErrorNotFound.WithDetail(details...) },
		sql.ErrConnDone:          func() *Error { return NewE(pos + sqlConnDoneMsg + msg).WithDetail(details...) },
		sql.ErrTxDone:            func() *Error { return NewE(pos + sqlTxDoneMsg + msg).WithDetail(details...) },
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

	return NewE(pos + sqlErrorMsg + msg).WithDetail(details...)
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
	msg, pos := CallerMsg(err.Error(), 1)
	return New(InternalServerError, pos+" redis: "+msg).WithDetail(details...)
}
