package greeting

import "testing"

func TestGreet(t *testing.T) {
	cases := []struct {
		name string // 子测试名称
		in   string // Greet 的入参
		want string // 期望返回值
	}{
		{name: "普通名字", in: "Ada", want: "Hello, Ada!"},
		{name: "空字符串回退默认值", in: "", want: "Hello, Gopher!"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := Greet(tc.in)
			if got != tc.want {
				t.Errorf("Greet(%q) = %q, 期望 %q", tc.in, got, tc.want)
			}
		})
	}
}
