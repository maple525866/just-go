// Package channel 演示 channel 通信：无缓冲、有缓冲、关闭与 select 超时。
//
// import 路径：just-go/stage-1-syntax/05-concurrency/channel
package channel

import "time"

// PingPong 使用无缓冲 channel 完成一次发送接收配对。
func PingPong(message string) string {
	ch := make(chan string)
	go func() {
		ch <- message
	}()
	return <-ch
}

// BufferedQueue 使用有缓冲 channel 暂存多个值，再关闭并读取。
func BufferedQueue(values []int) []int {
	ch := make(chan int, len(values))
	for _, v := range values {
		ch <- v
	}
	close(ch)

	out := make([]int, 0, len(values))
	for v := range ch {
		out = append(out, v)
	}
	return out
}

// CloseAndRange 启动发送方，关闭 channel 后由接收方 range 读取全部值。
func CloseAndRange(values []string) []string {
	ch := make(chan string)
	go func() {
		defer close(ch)
		for _, v := range values {
			ch <- v
		}
	}()

	out := make([]string, 0, len(values))
	for v := range ch {
		out = append(out, v)
	}
	return out
}

// ReceiveWithTimeout 使用 select 等待消息或超时。
func ReceiveWithTimeout(delay, timeout time.Duration) string {
	ch := make(chan string, 1)
	go func() {
		time.Sleep(delay)
		ch <- "ready"
	}()

	select {
	case msg := <-ch:
		return msg
	case <-time.After(timeout):
		return "timeout"
	}
}
