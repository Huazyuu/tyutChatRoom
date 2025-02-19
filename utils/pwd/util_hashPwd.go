package pwd

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPwd hash加密
func HashPwd(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// CheckPwd 解hash 验证密码
func CheckPwd(hashPwd string, pwd string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashPwd), []byte(pwd))
	if err != nil {
		return false, err
	}
	return true, nil
}
