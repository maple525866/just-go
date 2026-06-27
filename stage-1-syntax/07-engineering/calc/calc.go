// Package calc 提供适合 testing 与 benchmark 的确定性纯函数。
//
// import 路径：just-go/stage-1-syntax/07-engineering/calc
package calc

import "strings"

// Add 返回两个整数之和，用于最基础的表驱动测试。
func Add(a, b int) int {
	return a + b
}

// Fibonacci 返回第 n 个 Fibonacci 数；n<=0 时返回 0。
func Fibonacci(n int) int {
	if n <= 0 {
		return 0
	}
	if n == 1 {
		return 1
	}
	prev, curr := 0, 1
	for i := 2; i <= n; i++ {
		prev, curr = curr, prev+curr
	}
	return curr
}

// NormalizeWords 将文本拆词、去空白、转小写，适合作为 benchmark 对象。
func NormalizeWords(input string) []string {
	parts := strings.Fields(input)
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		out = append(out, strings.ToLower(part))
	}
	return out
}
