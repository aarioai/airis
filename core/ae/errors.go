package ae

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"regexp"

	"github.com/redis/go-redis/v9"
)

var (
	duplicateKeyPattern = regexp.MustCompile(`Duplicate\s+entry\s+'([^']*)'\s+for\s+key\s+'([^']*)'`)
)

func Catch(es ...*Error) *Error {
	for _, e := range es {
		if e != nil {
			return e
		}
	}
	return nil
}
func CatchError(es ...error) error {
	for _, e := range es {
		if e != nil {
			return e
		}
	}
	return nil
}

// NewSQLError 处理 SQL 错误
// @TODO
func NewSQLError(err error, details ...string) *Error {
	if err == nil {
		return nil
	}
	msg, pos := CallerMsg(err.Error(), 1)

	switch {
	case errors.Is(err, driver.ErrBadConn):
		return New(InternalServerError, pos+" sql bad conn: "+msg).withDetail(details...)
	case errors.Is(err, driver.ErrSkip):
		// ErrSkip may be returned by some optional interfaces' methods to
		// indicate at runtime that the fast path is unavailable and the sql
		// package should continue as if the optional interface was not
		// implemented. ErrSkip is only supported where explicitly
		// documented.
		return New(InternalServerError, pos+" sql skip: "+msg).withDetail(details...)
	case errors.Is(err, driver.ErrRemoveArgument):
		return New(InternalServerError, pos+" sql remove argument: "+msg).withDetail(details...)
	case errors.Is(err, sql.ErrNoRows):
		return NotFoundE.withDetail(details...) // 通过在 asql层，对数组转换为 ae.NoRows
	case errors.Is(err, sql.ErrConnDone):
		// ErrConnDone is returned by any operation that is performed on a connection
		// that has already been returned to the connection pool.
		return New(InternalServerError, pos+" sql conn done: "+msg).withDetail(details...)
	case errors.Is(err, sql.ErrTxDone):
		return New(InternalServerError, pos+" sql tx done: "+msg).withDetail(details...)
	}

	dupMatches := duplicateKeyPattern.FindAllStringSubmatch(msg, -1)
	if dupMatches != nil && len(dupMatches) > 0 && len(dupMatches[0]) == 3 {
		// dupMatches[0][1]
		return New(Conflict, "sql key conflict").withDetail(details...)
	}

	return New(InternalServerError, pos+" sql error: "+msg).withDetail(details...)
}

// NewRedisError 处理 Redis 错误
// @TODO
func NewRedisError(err error, details ...string) *Error {
	if err == nil {
		return nil
	}

	if errors.Is(err, redis.Nil) {
		return New(NotFound, "Key not found").withDetail(details...)
	}
	msg, pos := CallerMsg(err.Error(), 1)
	return New(InternalServerError, pos+" redis: "+msg).withDetail(details...)
}
