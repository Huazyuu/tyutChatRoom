package pwd

import (
	"testing"
)

func TestCheckPasswordLever(t *testing.T) {
	type args struct {
		pwd string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "密码长度不足",
			args: args{
				pwd: "Ab1",
			},
			wantErr: true,
		},
		{
			name: "密码缺少数字",
			args: args{
				pwd: "Abcdefgh",
			},
			wantErr: true,
		},
		{
			name: "密码缺少小写字母",
			args: args{
				pwd: "ABCDEFGH1",
			},
			wantErr: true,
		},
		{
			name: "密码缺少大写字母",
			args: args{
				pwd: "abcdefgh1",
			},
			wantErr: true,
		},
		{
			name: "密码符合要求",
			args: args{
				pwd: "Abcdefgh1",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CheckPasswordLever(tt.args.pwd); (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordLever() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
