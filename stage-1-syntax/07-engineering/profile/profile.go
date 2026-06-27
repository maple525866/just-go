// Package profile 说明 pprof CPU / 内存 / 阻塞分析的用途与命令入口。
//
// import 路径：just-go/stage-1-syntax/07-engineering/profile
package profile

// Profile 描述一种 pprof profile。
type Profile struct {
	Name    string
	Purpose string
}

// Types 返回常见 profile 类型及用途。
func Types() []Profile {
	return []Profile{
		{Name: "CPU", Purpose: "定位消耗 CPU 时间最多的函数"},
		{Name: "memory", Purpose: "观察分配热点和对象保留情况"},
		{Name: "blocking", Purpose: "分析 goroutine 在锁、channel 等位置的阻塞"},
	}
}

// Commands 返回常用 pprof 命令提示。
func Commands() []string {
	return []string{
		"go test -bench=. -cpuprofile=cpu.out ./stage-1-syntax/07-engineering/calc",
		"go tool pprof cpu.out",
		"go test -bench=. -memprofile=mem.out ./stage-1-syntax/07-engineering/calc",
	}
}
