package pwd

import (
	"testing"
)

func TestHashCheckPwd(t *testing.T) {
	type args struct {
		pwd string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "pwdHashCase1",
			args: args{pwd: "123456"},
		}, {
			name: "pwdHashCase2",
			args: args{pwd: "abcdef"},
		}, {
			name: "pwdHashCase3",
			args: args{"123456asdalskndlasdasx_dwqdas"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HashPwd(tt.args.pwd)
			if err != nil {
				t.Error(err)
			}
			if _, err = CheckPwd(got, tt.args.pwd); err != nil {
				t.Errorf("HashPwd() = %v, err: %v", got, err)
			}
		})
	}
}
