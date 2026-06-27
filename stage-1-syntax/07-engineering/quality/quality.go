// Package quality 汇总 Go 工程质量门禁命令。
//
// import 路径：just-go/stage-1-syntax/07-engineering/quality
package quality

// Check 表示一个本地或 CI 质量检查命令。
type Check struct {
	Name    string
	Command string
	Purpose string
}

// Checks 返回与仓库 CI 对齐的本地验证命令。
func Checks() []Check {
	return []Check{
		{Name: "vet", Command: "go vet ./...", Purpose: "发现可疑代码和格式化调用错误"},
		{Name: "test", Command: "go test -race -count=1 ./...", Purpose: "运行测试并启用 race detector"},
		{Name: "build", Command: "go build ./...", Purpose: "确认所有 package 可编译"},
		{Name: "lint", Command: "golangci-lint run", Purpose: "运行聚合 linter，与 CI lint job 对齐"},
	}
}

// Commands 只返回命令文本，便于 README 和入口报告展示。
func Commands() []string {
	checks := Checks()
	commands := make([]string, 0, len(checks))
	for _, check := range checks {
		commands = append(commands, check.Command)
	}
	return commands
}
