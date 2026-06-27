// Package format 演示 fmt 的格式化输出。
//
// import 路径：just-go/stage-1-syntax/06-stdlib-essentials/format
package format

import (
	"bytes"
	"fmt"
)

// Progress 使用 fmt.Sprintf 生成学习进度说明。
func Progress(chapter string, done, total int) string {
	return fmt.Sprintf("%s: %02d/%02d tasks", chapter, done, total)
}

// WriteSummary 使用 fmt.Fprintf 写入任意 writer，这里返回字符串便于测试断言。
func WriteSummary(topic string, count int) string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "topic=%s examples=%d", topic, count)
	return buf.String()
}
