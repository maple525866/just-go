package main

import (
	"fmt"
	"strings"
	"time"

	"just-go/stage-1-syntax/06-stdlib-essentials/clock"
	"just-go/stage-1-syntax/06-stdlib-essentials/codec"
	"just-go/stage-1-syntax/06-stdlib-essentials/format"
	"just-go/stage-1-syntax/06-stdlib-essentials/inspect"
	"just-go/stage-1-syntax/06-stdlib-essentials/stream"
	"just-go/stage-1-syntax/06-stdlib-essentials/system"
	"just-go/stage-1-syntax/06-stdlib-essentials/web"
)

func main() {
	fmt.Println(buildReport())
}

// buildReport 串联各标准库主题子包，组装标准库精要学习报告。
func buildReport() string {
	copied, bytesCopied, _ := stream.CopyText("hello stdlib")
	lines, _ := stream.ScanLines("fmt\nio\ntime")
	fileText, _ := system.WriteReadTempFile("temp file ok")
	goVersion, _ := system.GoVersion()
	jsonText, jsonLesson, _ := codec.JSONRoundTrip(codec.Lesson{ID: 6, Title: "stdlib", Done: true})
	xmlText, _, _ := codec.XMLRoundTrip(codec.Lesson{ID: 6, Title: "stdlib", Done: true})
	timeSummary := clock.FormatSummary(time.Date(2026, 6, 27, 0, 0, 0, 0, time.UTC), 2*time.Hour)
	typeName, fields := inspect.DescribeStruct(jsonLesson)

	return strings.Join([]string{
		reportTitle(),
		"",
		"## fmt（格式化）",
		format.Progress("chapter-06", 6, 15),
		format.WriteSummary("fmt", 2),
		"",
		"## io / bufio（流式读写）",
		fmt.Sprintf("io.Copy: %q bytes=%d", copied, bytesCopied),
		fmt.Sprintf("Scanner lines: %v", lines),
		"",
		"## os / os/exec（文件与进程）",
		fmt.Sprintf("temp file: %s", fileText),
		fmt.Sprintf("go version: %s", goVersion),
		"",
		"## net/http（handler + client）",
		fmt.Sprintf("handler type: %T", web.HelloHandler()),
		"",
		"## encoding/json / xml（序列化）",
		fmt.Sprintf("json: %s", jsonText),
		fmt.Sprintf("xml: %s", xmlText),
		"",
		"## time（格式化 + ticker）",
		fmt.Sprintf("time: %s minutes=%d ticks=%d", timeSummary.Formatted, timeSummary.Minutes, clock.CountTicks(1, time.Millisecond)),
		"",
		"## reflect（只读元数据）",
		fmt.Sprintf("type=%s fields=%v", typeName, fields),
	}, "\n")
}

// reportTitle 生成报告标题行（纯函数，便于 main_test 断言）。
func reportTitle() string {
	return "=== 标准库精要学习报告 ==="
}
