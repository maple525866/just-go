// Package apperr 演示 Go 错误处理：sentinel error、错误包装、errors.Is 与 errors.As。
//
// import 路径：just-go/stage-1-syntax/04-interface-error/apperr
package apperr

import (
	"errors"
	"fmt"
)

// ErrUserNotFound 是可用 errors.Is 判别的 sentinel error。
var ErrUserNotFound = errors.New("user not found")

// QueryError 是包含上下文信息的自定义错误类型，可用 errors.As 提取。
type QueryError struct {
	User string
	Op   string
	Err  error
}

// Error 实现 error 接口。
func (e *QueryError) Error() string {
	return fmt.Sprintf("%s %q: %v", e.Op, e.User, e.Err)
}

// Unwrap 暴露被包装的底层错误，供 errors.Is / errors.As 遍历错误链。
func (e *QueryError) Unwrap() error {
	return e.Err
}

// FindUser 模拟按用户名查询；缺失用户会返回带 %w 包装的错误链。
func FindUser(name string) (string, error) {
	if name == "Ada" {
		return "Ada Lovelace", nil
	}
	return "", fmt.Errorf("lookup profile: %w", &QueryError{
		User: name,
		Op:   "find user",
		Err:  ErrUserNotFound,
	})
}

// IsUserNotFound 封装 errors.Is 判别，便于调用方不依赖具体包装层数。
func IsUserNotFound(err error) bool {
	return errors.Is(err, ErrUserNotFound)
}

// ExtractQueryError 使用 errors.As 提取错误链中的 QueryError。
func ExtractQueryError(err error) (*QueryError, bool) {
	var queryErr *QueryError
	if errors.As(err, &queryErr) {
		return queryErr, true
	}
	return nil, false
}

// Summary 返回本章错误处理示例的关键点。
func Summary() string {
	return "使用 fmt.Errorf + %w 包装错误；用 errors.Is 判断 sentinel；用 errors.As 提取自定义错误"
}
