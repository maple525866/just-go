package main

import "testing"

func TestResolveName(t *testing.T) {
	cases := []struct {
		name string // 子测试名称
		in   string // resolveName 的入参
		want string // 期望返回值
	}{
		{name: "给定名字原样返回", in: "Grace", want: "Grace"},
		{name: "空字符串回退 world", in: "", want: "world"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := resolveName(tc.in)
			if got != tc.want {
				t.Errorf("resolveName(%q) = %q, 期望 %q", tc.in, got, tc.want)
			}
		})
	}
}
