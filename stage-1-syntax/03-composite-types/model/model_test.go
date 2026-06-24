package model

import "testing"

func TestStudentFieldPromotion(t *testing.T) {
	cases := []struct {
		name      string
		student   Student
		wantEmail string
		wantLabel string
	}{
		{
			name: "嵌入字段提升",
			student: Student{
				Name:    "Ada",
				Score:   92,
				Contact: Contact{Email: "ada@go.dev", Phone: "123"},
			},
			wantEmail: "ada@go.dev",
			wantLabel: "Ada(92) <ada@go.dev>",
		},
		{
			name: "空联系方式",
			student: Student{
				Name:  "Bob",
				Score: 78,
			},
			wantEmail: "",
			wantLabel: "Bob(78) <>",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// s.Email 是被提升的嵌入字段，等价于 s.Contact.Email
			if got := tc.student.Email; got != tc.wantEmail {
				t.Errorf("被提升字段 Email = %q, 期望 %q", got, tc.wantEmail)
			}
			if got := tc.student.Email; got != tc.student.Contact.Email {
				t.Errorf("s.Email(%q) 应等价于 s.Contact.Email(%q)", got, tc.student.Contact.Email)
			}
			if got := tc.student.Label(); got != tc.wantLabel {
				t.Errorf("Label() = %q, 期望 %q", got, tc.wantLabel)
			}
		})
	}
}
