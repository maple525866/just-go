// Package stream 演示 io.Reader / io.Writer、io.Copy 与 bufio.Scanner。
//
// import 路径：just-go/stage-1-syntax/06-stdlib-essentials/stream
package stream

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

// CopyText 使用 io.Copy 从 reader 复制到 writer。
func CopyText(input string) (string, int64, error) {
	reader := strings.NewReader(input)
	var writer bytes.Buffer
	n, err := io.Copy(&writer, reader)
	return writer.String(), n, err
}

// ScanLines 使用 bufio.Scanner 按行读取文本。
func ScanLines(input string) ([]string, error) {
	scanner := bufio.NewScanner(strings.NewReader(input))
	lines := make([]string, 0)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
