package main

import (
	"fmt"

	"just-go/stage-1-syntax/01-hello-go/greeting"
)

func main() {
	name := resolveName("")
	fmt.Println(greeting.Greet(name))
}

// resolveName 把「决定问候对象」的纯逻辑从带副作用的 main 中抽离出来，
// 便于在 main_test.go 中直接断言，无需捕获 stdout。
// arg 为空时返回 "world"，否则原样返回。
func resolveName(arg string) string {
	if arg == "" {
		return "world"
	}
	return arg
}
