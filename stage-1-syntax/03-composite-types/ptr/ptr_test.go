package ptr

import "testing"

func TestReceiverSemantics(t *testing.T) {
	cases := []struct {
		name        string
		start       int
		bonus       int
		wantValue   int // 值接收者返回的新余额
		wantPointer int // 指针接收者修改后的原对象余额
	}{
		{name: "加 10", start: 100, bonus: 10, wantValue: 110, wantPointer: 110},
		{name: "加 0", start: 50, bonus: 0, wantValue: 50, wantPointer: 50},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			acc := Account{Balance: tc.start}

			// 值接收者：返回新值，但不改原对象
			got := acc.WithBonus(tc.bonus)
			if got.Balance != tc.wantValue {
				t.Errorf("WithBonus 返回余额 = %d, 期望 %d", got.Balance, tc.wantValue)
			}
			if acc.Balance != tc.start {
				t.Errorf("值接收者不应改动原对象：acc.Balance = %d, 期望保持 %d", acc.Balance, tc.start)
			}

			// 指针接收者：原地修改原对象
			acc.AddBonus(tc.bonus)
			if acc.Balance != tc.wantPointer {
				t.Errorf("指针接收者应原地修改：acc.Balance = %d, 期望 %d", acc.Balance, tc.wantPointer)
			}
		})
	}
}

func TestDeref(t *testing.T) {
	x := 42
	if got := Deref(&x); got != 42 {
		t.Errorf("Deref(&x) = %d, 期望 42", got)
	}
}

func TestDoubleInPlace(t *testing.T) {
	x := 21
	DoubleInPlace(&x)
	if x != 42 {
		t.Errorf("DoubleInPlace 后 x = %d, 期望 42", x)
	}
}
