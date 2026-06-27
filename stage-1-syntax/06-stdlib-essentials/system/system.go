// Package system 演示 os 与 os/exec：临时文件、环境变量、外部命令。
//
// import 路径：just-go/stage-1-syntax/06-stdlib-essentials/system
package system

import (
	"os"
	"os/exec"
	"strings"
)

// WriteReadTempFile 在系统临时目录写入并读取文件，不污染固定工作目录。
func WriteReadTempFile(content string) (string, error) {
	file, err := os.CreateTemp("", "just-go-stdlib-*.txt")
	if err != nil {
		return "", err
	}
	name := file.Name()
	defer os.Remove(name)

	if _, err := file.WriteString(content); err != nil {
		file.Close()
		return "", err
	}
	if err := file.Close(); err != nil {
		return "", err
	}

	data, err := os.ReadFile(name)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// EnvOrDefault 读取环境变量，不存在时返回默认值。
func EnvOrDefault(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// GoVersion 使用 os/exec 调用 Go 工具链，避免依赖平台专属 shell 命令。
func GoVersion() (string, error) {
	out, err := exec.Command("go", "env", "GOVERSION").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
