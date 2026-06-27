// Package pitfall 以安全方式说明常见并发坑。
//
// import 路径：just-go/stage-1-syntax/05-concurrency/pitfall
package pitfall

// Pitfall 描述一种并发风险及其规避方式。
type Pitfall struct {
	Name       string
	Risk       string
	Prevention string
}

// Summaries 返回常见并发坑的摘要。
//
// 注意：本包不提交真正会触发 data race、goroutine leak 或 channel deadlock 的代码，
// 因为这些示例会破坏测试稳定性。教学时用摘要说明风险，用 go test -race 验证安全代码。
func Summaries() []Pitfall {
	return []Pitfall{
		{
			Name:       "data race",
			Risk:       "多个 goroutine 未同步读写同一变量会产生数据竞争",
			Prevention: "用 channel 传递所有权，或用 sync.Mutex / atomic 保护共享状态，并运行 go test -race",
		},
		{
			Name:       "goroutine leak",
			Risk:       "发送方或接收方永久阻塞会让 goroutine 无法退出",
			Prevention: "使用 context 取消、关闭 channel，或确保每个 goroutine 有明确退出路径",
		},
		{
			Name:       "channel deadlock",
			Risk:       "没有对应接收方的发送、没有发送方的接收，或重复等待会导致死锁",
			Prevention: "设计清晰的 channel 所有者，发送方负责 close，并用 select timeout 暴露异常路径",
		},
	}
}

// Names 返回风险名称，便于入口报告展示。
func Names() []string {
	items := Summaries()
	names := make([]string, 0, len(items))
	for _, item := range items {
		names = append(names, item.Name)
	}
	return names
}
