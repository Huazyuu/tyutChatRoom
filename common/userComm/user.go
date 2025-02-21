package userComm

type UserRegisterRequest struct {
	UserName string `json:"username" binding:"required" msg:"请输入用户名"`
	Password string `json:"password" binding:"required" msg:"请输入密码"`
	Email    string `json:"email" binding:"required" msg:"请输入邮箱"`
}
type UserLoginRequest struct {
	Email    string  `json:"email" binding:"required,email" msg:"非法邮箱"`
	Code     *string `json:"code"`
	Password string  `json:"password"`
}
type UserListRequest struct {
	Page  int    `form:"page"`
	Limit int    `form:"limit"`
	Sort  string `form:"sort"`
}
