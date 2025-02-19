package pwd

import (
	"errors"
	"regexp"
)

func CheckPasswordLever(pwd string) error {
	if len(pwd) < 8 {
		return errors.New("password len is < 9")
	}
	num := `[0-9]{1}`
	a_z := `[a-z]{1}`
	A_Z := `[A-Z]{1}`
	// symbol := `[!@#~$%^&*()+|_]{1}`
	if b, err := regexp.MatchString(num, pwd); !b || err != nil {
		return errors.New("password need num")
	}
	if b, err := regexp.MatchString(a_z, pwd); !b || err != nil {
		return errors.New("password need lowercase")
	}
	if b, err := regexp.MatchString(A_Z, pwd); !b || err != nil {
		return errors.New("password need capital letter")
	}
	return nil
}
