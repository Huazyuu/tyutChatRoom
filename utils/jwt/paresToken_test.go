package jwt

import (
	"testing"
)

func TestParseToken(t *testing.T) {

	type args struct {
		tokenStr string
	}
	test := struct {
		name string
		args args
	}{
		name: "testUser",
		args: args{
			tokenStr: TOKEN,
		},
	}

	t.Run(test.name, func(t *testing.T) {

		res, err := ParseToken(test.args.tokenStr)

		if err != nil {
			t.Error("invalid token")
			return
		}
		t.Log("res: ", res)
		t.Log("ParseToken_CASE SUCCESS")
	})

}
