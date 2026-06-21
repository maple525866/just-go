// Package funcs 演示多返回值、命名返回、可变参数与闭包。
package funcs

import "fmt"

// MinMax 返回可变参数中的最小值与最大值；空参数时 ok 为 false。
func MinMax(nums ...int) (min, max int, ok bool) {
	if len(nums) == 0 {
		return 0, 0, false
	}
	min, max = nums[0], nums[0]
	for _, n := range nums[1:] {
		if n < min {
			min = n
		}
		if n > max {
			max = n
		}
	}
	return min, max, true
}

// Average 计算可变参数的算术平均值；空参数返回 0（命名返回值 avg）。
func Average(nums ...int) (avg float64) {
	if len(nums) == 0 {
		return 0
	}
	sum := 0
	for _, n := range nums {
		sum += n
	}
	avg = float64(sum) / float64(len(nums))
	return
}

// MakeGrader 返回一个闭包：捕获 threshold，多次调用时共享同一阈值配置。
func MakeGrader(threshold int) func(int) string {
	return func(score int) string {
		if score >= threshold {
			return "pass"
		}
		return "fail"
	}
}

// MinMaxLine 把 MinMax 结果格式化为可读字符串。
func MinMaxLine(nums ...int) string {
	min, max, ok := MinMax(nums...)
	if !ok {
		return "minmax: (empty)"
	}
	return fmt.Sprintf("minmax: min=%d max=%d", min, max)
}
