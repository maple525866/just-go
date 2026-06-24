package seq

import "testing"

func TestGrowSteps(t *testing.T) {
	cases := []struct {
		name string
		n    int
	}{
		{name: "零次", n: 0},
		{name: "一次", n: 1},
		{name: "四次", n: 4},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			steps := GrowSteps(tc.n)
			if len(steps) != tc.n {
				t.Fatalf("GrowSteps(%d) 返回 %d 个快照, 期望 %d", tc.n, len(steps), tc.n)
			}
			for i, snap := range steps {
				if gotLen := snap[0]; gotLen != i+1 {
					t.Errorf("第 %d 次 append 后 len=%d, 期望 %d", i+1, gotLen, i+1)
				}
				if snap[1] < snap[0] {
					t.Errorf("第 %d 次 append 后 cap=%d 不应小于 len=%d", i+1, snap[1], snap[0])
				}
			}
		})
	}
}

func TestSubSliceMutationDemo(t *testing.T) {
	base, sub := SubSliceMutationDemo()

	if sub[0] != 99 {
		t.Errorf("sub[0] = %d, 期望 99", sub[0])
	}
	if base[1] != 99 {
		t.Errorf("共享底层数组：base[1] 应被 sub 的修改污染为 99, got %d", base[1])
	}
}

func TestArrayValueCopy(t *testing.T) {
	original, modified := ArrayValueCopy()

	if original[0] != 1 {
		t.Errorf("数组是值类型：原数组不应被副本修改影响, original[0] = %d, 期望 1", original[0])
	}
	if modified[0] != 99 {
		t.Errorf("modified[0] = %d, 期望 99", modified[0])
	}
}
