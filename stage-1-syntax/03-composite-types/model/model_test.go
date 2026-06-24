package model

import "testing"

func TestStudentFieldPromotion(t *testing.T) {
	cases := []struct {
		name      string
		sName     string
		score     int
		contact   Contact
		wantLabel string
	}{
		{
			name:      "嵌入字段提升",
			sName:     "Ada",
			score:     92,
			contact:   Contact{Email: "ada@go.dev", Phone: "123"},
			wantLabel: "Ada(92) <ada@go.dev>",
		},
		{
			name:      "空联系方式",
			sName:     "Bob",
			score:     78,
			contact:   Contact{},
			wantLabel: "Bob(78) <>",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := Student{Name: tc.sName, Score: tc.score, Contact: tc.contact}

			// s.Email 是被提升的嵌入字段：无需写 s.Contact.Email 即可访问，
			// 且其值等于嵌入进来的 Contact 的 Email。
			if got := s.Email; got != tc.contact.Email {
				t.Errorf("被提升字段 Email = %q, 期望 %q", got, tc.contact.Email)
			}
			if got := s.Label(); got != tc.wantLabel {
				t.Errorf("Label() = %q, 期望 %q", got, tc.wantLabel)
			}
		})
	}
}
